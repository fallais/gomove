package cmd

import (
	"gomove/internal/cmd/start"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the mouse mover service",
	Long:  `Start the mouse mover service that will periodically move the mouse cursor to prevent your session from locking. The movement is minimal and shouldn't interfere with your work.`,
	Run:   start.Run,
}

func init() {
	// Here you will define your flags and configuration settings.
	startCmd.Flags().IntP("interval", "i", 60, "Interval in seconds between mouse movements")
	startCmd.Flags().IntP("distance", "d", 1, "Distance in pixels to move the mouse")

	// Bind flags to viper
	viper.BindPFlag("interval", startCmd.Flags().Lookup("interval"))
	viper.BindPFlag("distance", startCmd.Flags().Lookup("distance"))
}
