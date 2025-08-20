package cmd

import (
	"errors"
	"fmt"
	"gomove/pkg/mouse"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the mouse mover service",
	Long: `Start the mouse mover service that will periodically move the mouse cursor
to prevent your session from locking. The movement is minimal and shouldn't
interfere with your work.`,
	Run: func(cmd *cobra.Command, args []string) {
		startMouseMover()
	},
}

func init() {
	rootCmd.AddCommand(startCmd)

	// Here you will define your flags and configuration settings.
	startCmd.Flags().IntP("interval", "i", 60, "Interval in seconds between mouse movements")
	startCmd.Flags().IntP("distance", "d", 1, "Distance in pixels to move the mouse")

	// Bind flags to viper
	viper.BindPFlag("interval", startCmd.Flags().Lookup("interval"))
	viper.BindPFlag("distance", startCmd.Flags().Lookup("distance"))
}

func startMouseMover() {
	// Get configuration values
	interval := viper.GetInt("interval")
	distance := viper.GetInt("distance")

	if interval <= 0 {
		interval = 60 // default to 60 seconds
	}
	if distance <= 0 {
		distance = 1 // default to 1 pixel
	}

	fmt.Printf("Starting mouse mover with interval: %d seconds, distance: %d pixels\n", interval, distance)
	fmt.Println("Press Ctrl+C to stop...")

	// Create a channel to listen for interrupt signals
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)

	// Create a ticker for periodic mouse movement
	ticker := time.NewTicker(time.Duration(interval) * time.Second)
	defer ticker.Stop()

	// Main loop
	for {
		select {
		case <-ticker.C:
			err := mouse.Move(distance, time.Duration(interval)*time.Second)
			if err != nil {
				if errors.Is(err, mouse.ErrDurationTooLong) {
					log.Printf("Warning: Movement duration exceeds interval (%d seconds). Consider increasing interval or reducing movement complexity.\n", interval)
				} else if errors.Is(err, mouse.ErrUserInterruption) {
					// User interrupted movement, this is normal - no need to log
				} else {
					log.Printf("Error moving mouse: %v\n", err)
				}
			} else {
				fmt.Printf("Mouse moved at %s\n", time.Now().Format("15:04:05"))
			}
		case <-signalChan:
			fmt.Println("\nReceived interrupt signal. Stopping mouse mover...")
			return
		}
	}
}
