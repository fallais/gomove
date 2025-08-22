package cmd

import (
	"gomove/internal/cmd/config/get"

	"github.com/spf13/cobra"
)

// getConfigCmd represents the config get command
var getConfigCmd = &cobra.Command{
	Use:   "get <parameter>",
	Short: "Get a specific configuration parameter value",
	Long: `Get the value of a specific configuration parameter. Supports nested parameters using dot notation.

Examples:
  gomove config get debug
  gomove config get behavior.idle_timeout
  gomove config get activities.0.kind
  gomove config get activities.0.interval`,
	Run: get.Run,
}
