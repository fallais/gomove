package cmd

import (
	"gomove/internal/cmd/config/create"

	"github.com/spf13/cobra"
)

// createConfigCmd represents the create config command
var createConfigCmd = &cobra.Command{
	Use:   "create",
	Short: "Creates a new configuration file",
	Long:  `Creates a new configuration file in ~/.gomove/config.yaml with default settings. This will create the directory and file if they don't exist.`,
	Run:   create.Run,
}
