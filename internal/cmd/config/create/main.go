package create

import (
	"gomove/pkg/log"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

const DefaultConfigurationFile = `
behavior:
  idle_timeout: 10s
  resume_after_inactivity: true
  pause_when_user_is_active: true
activities:
  - kind: mouse
    interval: 3s
debug: true
logfile: ""
`

func Run(cmd *cobra.Command, args []string) {
	// Get home directory
	home, err := os.UserHomeDir()
	if err != nil {
		log.Error("error getting home directory", zap.Error(err))
		return
	}

	// Create .gomove directory
	configDir := filepath.Join(home, ".gomove")
	err = os.MkdirAll(configDir, 0755)
	if err != nil {
		log.Error("error creating config directory", zap.Error(err))
		return
	}

	// Create config file path
	configFile := filepath.Join(configDir, "config.yaml")

	// Check if config file already exists
	if _, err := os.Stat(configFile); err == nil {
		log.Info("configuration file already exists", zap.String("path", configFile))
		return
	}

	// Write config file
	err = os.WriteFile(configFile, []byte(DefaultConfigurationFile), 0644)
	if err != nil {
		log.Error("Error creating config file", zap.Error(err))
		return
	}

	log.Info("configuration file created successfully", zap.String("path", configFile))
	log.Info("you can now edit this file to customize your settings")
}
