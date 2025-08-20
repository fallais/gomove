package cmd

import (
	"gomove/internal/cmd/config/show"

	"github.com/spf13/cobra"
)

// showConfigCmd represents the config show command
var showConfigCmd = &cobra.Command{
	Use:   "show",
	Short: "Show current configuration status",
	Long:  `Display the current configuration values and the location of the config file. This helps you verify your settings before starting the mouse mover.`,
	Run:   show.Run,
}
