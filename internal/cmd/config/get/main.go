package get

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func Run(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		fmt.Println("Error: Please specify a configuration parameter to get")
		fmt.Println("Usage: gomove config get <parameter>")
		fmt.Println("Examples:")
		fmt.Println("  gomove config get debug")
		fmt.Println("  gomove config get behavior.idle_timeout")
		fmt.Println("  gomove config get activities.0.kind")
		return
	}

	param := args[0]

	// Check if the key exists in the configuration
	if !viper.IsSet(param) {
		fmt.Printf("Configuration parameter '%s' not found\n", param)
		return
	}

	value := viper.Get(param)
	fmt.Printf("%s = %v\n", param, value)
}
