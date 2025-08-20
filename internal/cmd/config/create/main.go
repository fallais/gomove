package create

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

func Run(cmd *cobra.Command, args []string) {
	// Get home directory
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Printf("Error getting home directory: %v\n", err)
		return
	}

	// Create .gomove directory
	configDir := filepath.Join(home, ".gomove")
	err = os.MkdirAll(configDir, 0755)
	if err != nil {
		fmt.Printf("Error creating config directory: %v\n", err)
		return
	}

	// Create config file path
	configFile := filepath.Join(configDir, "config.yaml")

	// Check if config file already exists
	if _, err := os.Stat(configFile); err == nil {
		fmt.Printf("Config file already exists at: %s\n", configFile)
		return
	}

	// Create default config content
	defaultConfig := `# GoMove Configuration File
# This file configures the mouse mover behavior

# Interval between mouse movements (in seconds)
interval: 60

# Distance to move the mouse (in pixels)
distance: 1

# Enable/disable debug output
debug: false

# Log file path (optional, leave empty for stdout)
logfile: ""
`

	// Write config file
	err = os.WriteFile(configFile, []byte(defaultConfig), 0644)
	if err != nil {
		fmt.Printf("Error creating config file: %v\n", err)
		return
	}

	fmt.Printf("Configuration file created successfully at: %s\n", configFile)
	fmt.Println("You can now edit this file to customize your settings.")
}
