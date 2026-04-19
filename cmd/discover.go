package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/AmadlaOrg/unravel/plugin"
	"github.com/spf13/cobra"
)

var pluginService plugin.Service

func init() {
	pluginService = plugin.New()
}

// DiscoverCmd discovers system state via unravel-* plugins and outputs HERY entities.
var DiscoverCmd = &cobra.Command{
	Use:   "discover",
	Short: "Discover system state and output as HERY entities",
	Long: `Scans for unravel-* plugins on PATH and invokes their discover subcommand.
Each plugin outputs JSON entities describing the actual state of the system.

Use --from to target a specific plugin, or run all discovered plugins.
Use --type to filter output to a specific entity type.`,
	RunE: runDiscover,
}

func init() {
	DiscoverCmd.Flags().String("from", "", "Plugin to use (e.g., 'system' for unravel-system)")
	DiscoverCmd.Flags().String("type", "", "Filter output to a specific entity type URI")
	DiscoverCmd.Flags().StringP("output", "o", "json", "Output format: json, yaml")
	DiscoverCmd.Flags().StringP("file", "f", "", "Input file (entity data or query)")
}

func runDiscover(cmd *cobra.Command, args []string) error {
	from, _ := cmd.Flags().GetString("from")
	entityType, _ := cmd.Flags().GetString("type")
	file, _ := cmd.Flags().GetString("file")

	// Determine stdin
	var stdin io.Reader
	if file != "" && file != "-" {
		f, err := os.Open(file)
		if err != nil {
			return fmt.Errorf("cannot open file %s: %w", file, err)
		}
		defer f.Close()
		stdin = f
	} else {
		stdin = os.Stdin
	}

	if from != "" {
		// Run a specific plugin
		return runPlugin(from, entityType, stdin)
	}

	// Run all discovered plugins
	plugins, err := pluginService.Discover()
	if err != nil {
		return fmt.Errorf("plugin discovery failed: %w", err)
	}

	if len(plugins) == 0 {
		fmt.Fprintln(os.Stderr, "no unravel-* plugins found in PATH")
		return nil
	}

	for _, pluginName := range plugins {
		shortName := strings.TrimPrefix(pluginName, "unravel-")
		if err := runPlugin(shortName, entityType, stdin); err != nil {
			fmt.Fprintf(os.Stderr, "warning: plugin %s failed: %v\n", pluginName, err)
		}
	}

	return nil
}

func runPlugin(name, entityType string, stdin io.Reader) error {
	pluginName := "unravel-" + name

	pluginArgs := []string{"discover"}
	if entityType != "" {
		pluginArgs = append(pluginArgs, "--type", entityType)
	}

	exitCode, err := pluginService.Exec(pluginName, pluginArgs, stdin, os.Stdout, os.Stderr)
	if err != nil {
		return err
	}
	if exitCode != 0 {
		return fmt.Errorf("plugin %s exited with code %d", pluginName, exitCode)
	}
	return nil
}

// filterByType filters JSON entity output by _type field.
// Used when --type is set and the plugin doesn't support it natively.
func filterByType(data []byte, entityType string) ([]byte, error) {
	var entities []map[string]any
	if err := json.Unmarshal(data, &entities); err != nil {
		// Try as single entity
		var entity map[string]any
		if err2 := json.Unmarshal(data, &entity); err2 != nil {
			return data, nil // not JSON entities, pass through
		}
		entities = []map[string]any{entity}
	}

	var filtered []map[string]any
	for _, e := range entities {
		if t, ok := e["_type"].(string); ok && strings.Contains(t, entityType) {
			filtered = append(filtered, e)
		}
	}

	return json.Marshal(filtered)
}
