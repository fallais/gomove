package show

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func Run(cmd *cobra.Command, args []string) {
	fmt.Println("GoMove Configuration Status")
	fmt.Println("===========================")

	// Show config file location
	configFile := viper.ConfigFileUsed()
	if configFile == "" {
		home, err := os.UserHomeDir()
		if err == nil {
			defaultConfigFile := filepath.Join(home, ".gomove", "config.yaml")
			if _, err := os.Stat(defaultConfigFile); err == nil {
				configFile = defaultConfigFile
			} else {
				configFile = "No config file found"
			}
		}
	}

	fmt.Printf("Config file: %s\n", configFile)
	fmt.Println()

	// Show current settings
	fmt.Println("Current Settings:")
	fmt.Printf("  Interval: %d seconds\n", viper.GetInt("interval"))
	fmt.Printf("  Distance: %d pixels\n", viper.GetInt("distance"))
	fmt.Printf("  Debug: %t\n", viper.GetBool("debug"))

	logfile := viper.GetString("logfile")
	if logfile == "" {
		logfile = "stdout"
	}
	fmt.Printf("  Log output: %s\n", logfile)

	fmt.Println()
	fmt.Println("To start the mouse mover, run: gomove start")
	fmt.Println("To initialize a config file, run: gomove config init")
}
