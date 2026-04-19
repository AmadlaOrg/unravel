package main

import (
	"fmt"
	"os"

	"github.com/AmadlaOrg/unravel/cmd"
	"github.com/spf13/cobra"
)

const (
	appName = "unravel"
	version = "1.0.0"
)

var rootCmd = &cobra.Command{
	Use:     appName,
	Short:   "System discovery CLI — outputs existing system state as HERY entities",
	Version: version,
}

func init() {
	rootCmd.AddCommand(cmd.DiscoverCmd)
	rootCmd.AddCommand(cmd.PluginsCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
