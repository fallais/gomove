package show

import (
	"bytes"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

func Run(cmd *cobra.Command, args []string) {
	// Get all settings from Viper
	allSettings := viper.AllSettings()

	// Create a custom YAML encoder with 2-space indentation
	var buf bytes.Buffer
	encoder := yaml.NewEncoder(&buf)
	encoder.SetIndent(2) // Set indentation to 2 spaces

	err := encoder.Encode(allSettings)
	if err != nil {
		fmt.Printf("Error formatting configuration as YAML: %v\n", err)
		return
	}
	encoder.Close()

	// Output the YAML directly without any headers
	fmt.Print(buf.String())
}
