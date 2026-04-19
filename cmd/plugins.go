package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var (
	pluginsOutputFlag string
	pluginsHeryFlag   bool
)

func init() {
	PluginsCmd.Flags().StringVarP(&pluginsOutputFlag, "output", "o", "table", "Output format: table, json, yaml")
	PluginsCmd.Flags().BoolVar(&pluginsHeryFlag, "hery", false, "Wrap output in HERY envelope (_type, _body)")
}

// PluginsCmd lists all discovered unravel-* plugins.
var PluginsCmd = &cobra.Command{
	Use:   "plugins",
	Short: "List all discovered unravel-* plugins",
	RunE:  runPlugins,
}

type heryEnvelope struct {
	Type string `json:"_type" yaml:"_type"`
	Body any    `json:"_body" yaml:"_body"`
}

type pluginRow struct {
	Name        string `json:"name" yaml:"name"`
	Backend     string `json:"backend" yaml:"backend"`
	Version     string `json:"version" yaml:"version"`
	Description string `json:"description" yaml:"description"`
}

func runPlugins(cmd *cobra.Command, args []string) error {
	plugins, err := pluginService.Discover()
	if err != nil {
		return fmt.Errorf("plugin discovery failed: %w", err)
	}

	if len(plugins) == 0 {
		fmt.Fprintln(os.Stderr, "no unravel-* plugins found in PATH")
		return nil
	}

	var rows []pluginRow
	for _, pluginName := range plugins {
		info, err := pluginService.GetInfo(pluginName)
		if err != nil {
			shortName := strings.TrimPrefix(pluginName, "unravel-")
			rows = append(rows, pluginRow{
				Name:        shortName,
				Backend:     "?",
				Version:     "?",
				Description: fmt.Sprintf("error: %v", err),
			})
			continue
		}
		rows = append(rows, pluginRow{
			Name:        info.Name,
			Backend:     info.Backend,
			Version:     info.Version,
			Description: info.Description,
		})
	}

	var data any = rows
	if pluginsHeryFlag {
		data = heryEnvelope{
			Type: "amadla.org/entity/tools/plugins@v1.0.0",
			Body: rows,
		}
	}

	switch pluginsOutputFlag {
	case "json":
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(data)
	case "yaml":
		return yaml.NewEncoder(os.Stdout).Encode(data)
	default:
		table := tablewriter.NewWriter(os.Stdout)
		table.Header("Name", "Backend", "Version", "Description")
		for _, r := range rows {
			table.Append(r.Name, r.Backend, r.Version, r.Description)
		}
		table.Render()
		return nil
	}
}
