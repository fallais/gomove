package cmd

import (
	"gomove/internal/cmd/start"

	"github.com/spf13/cobra"
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the mouse mover service",
	Long:  `Start the mouse mover service that will periodically move the mouse cursor to prevent your session from locking. The movement is minimal and shouldn't interfere with your work.`,
	Run:   start.Run,
}
