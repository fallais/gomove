package start

import (
	"errors"
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

	mover := mouse.NewMover(distance, time.Duration(interval)*time.Second)

	// Check if the estimated duration exceeds the interval
	if mover.EstimatedDuration() > time.Duration(interval)*time.Second {
		log.Warn("estimated move duration exceeds interval, increase the interval", zap.Int("interval", interval))
		return
	}

	// Main loop
	for {
		select {
		case <-ticker.C:
			err := mover.Move()
			if err != nil {
				if errors.Is(err, mouse.ErrUserInterruption) {
					log.Info("user interrupted movement")
				} else if errors.Is(err, mouse.ErrUserAlreadyMoving) {
					log.Info("user is already moving the cursor")
				} else {
					log.Error("Error moving mouse", zap.Error(err))
				}
			} else {
				log.Info("mouse moved", zap.Time("at", time.Now()))
			}
		case <-signalChan:
			log.Info("Received interrupt signal. Stopping mouse mover...")
			return
		}
	}
}
