package start

import (
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"gomove/pkg/log"
	"gomove/pkg/mouse"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func Run(cmd *cobra.Command, args []string) {
	// Get configuration values
	interval := viper.GetInt("interval")
	distance := viper.GetInt("distance")

	if interval <= 0 {
		interval = 60 // default to 60 seconds
	}
	if distance <= 0 {
		distance = 1 // default to 1 pixel
	}

	log.Info("Starting mouse mover", zap.Int("interval", interval), zap.Int("distance", distance))
	log.Info("Press Ctrl+C to stop...")

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
					log.Warn("movement duration exceeds interval", zap.Int("interval", interval))
				} else if errors.Is(err, mouse.ErrUserInterruption) {
					log.Debug("user interrupted movement", zap.Error(err))
				} else if errors.Is(err, mouse.ErrUserAlreadyMoving) {
					log.Debug("user is already moving the cursor", zap.Error(err))
				} else {
					log.Error("Error moving mouse", zap.Error(err))
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
