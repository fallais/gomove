package create

import (
	"gomove/pkg/log"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

func Run(cmd *cobra.Command, args []string) {
	// Get home directory
	home, err := os.UserHomeDir()
	if err != nil {
		log.Error("Error getting home directory", zap.Error(err))
		return
	}

	// Create .gomove directory
	configDir := filepath.Join(home, ".gomove")
	err = os.MkdirAll(configDir, 0755)
	if err != nil {
		log.Error("Error creating config directory", zap.Error(err))
		return
	}

	// Create config file path
	configFile := filepath.Join(configDir, "config.yaml")

	// Check if config file already exists
	if _, err := os.Stat(configFile); err == nil {
		log.Info("Config file already exists", zap.String("path", configFile))
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
		log.Error("Error creating config file", zap.Error(err))
		return
	}

	log.Info("Configuration file created successfully", zap.String("path", configFile))
	log.Info("You can now edit this file to customize your settings.")
}
