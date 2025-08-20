package cmd

import (
	"github.com/spf13/cobra"
)

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage configuration",
	Long:  `Manage the gomove configuration file.`,
}

func init() {
	configCmd.AddCommand(createConfigCmd)
	configCmd.AddCommand(showConfigCmd)
}
