package activity

import (
	"errors"
	"gomove/internal/models"
	"gomove/pkg/log"
	"gomove/pkg/mouse"
	"gomove/pkg/watcher"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"go.uber.org/zap"
)

type ActivityManager struct {
	behavior   models.Behavior
	activities []models.Activity

	watcher *watcher.Watcher
	mover   *mouse.Mover

	isPaused bool

	tickers map[string]*time.Ticker

	signal chan os.Signal

	mutex sync.RWMutex
}

func NewActivityManager(behavior models.Behavior, activities []models.Activity, watcher *watcher.Watcher, mover *mouse.Mover) *ActivityManager {
	return &ActivityManager{
		behavior:   behavior,
		activities: activities,
		watcher:    watcher,
		mover:      mover,
		signal:     make(chan os.Signal, 1),
	}
}

func (am *ActivityManager) Start() {
	log.Debug("starting activity manager")

	// Create a channel to listen for interrupt signals
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)

	// Start the watcher
	am.watcher.Start()

	// Create tickers for each activity
	am.tickers = make(map[string]*time.Ticker)
	for _, activity := range am.activities {
		log.Debug("adding ticker", zap.String("activity", string(activity.Kind)))
		am.tickers[string(activity.Kind)] = time.NewTicker(activity.Interval)
	}

	events := make(chan string)

	// Fan-in: forward each ticker into events channel
	for kind, ticker := range am.tickers {
		go func(n string, t *time.Ticker) {
			for range t.C {
				events <- n
			}
		}(kind, ticker)
	}

	for {
		select {
		case e := <-events:
			switch e {
			case string(models.KindMouse):
				if am.watcher.IsUserMoving() {
					log.Info("user is moving the mouse, skipping activity", zap.String("activity", e))
					continue
				}

				am.handleMouse()
			case string(models.KindKeyboard):
				// Handle keyboard input event
			}

		case <-signalChan:
			log.Info("received interrupt signal")
			log.Info("stopping mouse mover")
			return
		}
	}
}

func (am *ActivityManager) handleMouse() {
	err := am.mover.Move()
	if err != nil {
		if errors.Is(err, mouse.ErrUserInterruption) {
			log.Info("user interrupted movement")
		} else if errors.Is(err, mouse.ErrUserAlreadyMoving) {
			log.Info("user is already moving the cursor")
		} else {
			log.Error("error moving mouse", zap.Error(err))
		}
	} else {
		log.Info("mouse moved", zap.Time("at", time.Now()))
	}
}
