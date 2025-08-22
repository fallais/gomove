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

	"github.com/go-vgo/robotgo"
	"go.uber.org/zap"
)

type ActivityManager struct {
	behavior   models.Behavior
	activities []models.Activity

	watcher *watcher.Watcher
	mover   *mouse.Mover

	isPaused bool
	pausedAt time.Time

	activityTickers map[string]*time.Ticker
	pauseTicker     *time.Ticker

	signal chan os.Signal

	mutex sync.RWMutex
}

func NewActivityManager(behavior models.Behavior, activities []models.Activity, watcher *watcher.Watcher) *ActivityManager {
	// Create the mouse mover
	mover := mouse.NewMover()

	return &ActivityManager{
		behavior:    behavior,
		activities:  activities,
		watcher:     watcher,
		mover:       mover,
		signal:      make(chan os.Signal, 1),
		pauseTicker: time.NewTicker(500 * time.Millisecond),
	}
}

func (am *ActivityManager) Start() {
	log.Debug("starting activity manager")

	// Create a channel to listen for interrupt signals
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)

	// Start the watcher
	am.watcher.Start()

	// Start the resume goroutine if ResumeAfterInactivity is enabled
	if am.behavior.ResumeAfterInactivity {
		go am.resume()
	}

	// Create tickers for each activity
	am.activityTickers = make(map[string]*time.Ticker)
	for _, activity := range am.activities {
		if activity.Enabled != nil && !*activity.Enabled {
			log.Debug("activity is not enabled", zap.String("activity", string(activity.Kind)))
			continue
		}

		log.Debug("adding ticker", zap.String("activity", string(activity.Kind)))
		am.activityTickers[string(activity.Kind)] = time.NewTicker(activity.Interval)
	}

	events := make(chan string)

	// Fan-in: forward each ticker into events channel
	for kind, ticker := range am.activityTickers {
		go func(n string, t *time.Ticker) {
			for range t.C {
				events <- n
			}
		}(kind, ticker)
	}

	for {
		select {
		case e := <-events:
			// Check if the activity manager is paused
			am.mutex.RLock()
			if am.isPaused {
				log.Debug("activity manager is paused, skipping activity", zap.String("activity", e))
				am.mutex.RUnlock()
				continue
			}
			am.mutex.RUnlock()

			// Check if user is active
			if am.watcher.IsUserMoving() {
				log.Debug("user is moving the mouse, skipping activity", zap.String("activity", e))

				if am.behavior.PauseWhenUserIsActive {
					am.pause()
				}

				continue
			}

			// Handle the event kind
			switch e {
			case string(models.KindMouse):
				am.handleMouse()
			case string(models.KindKeyboard):
				am.handleKeyboard()
			}

		case <-am.pauseTicker.C:
			am.checkAndPauseIfNeeded()

		case <-signalChan:
			log.Info("received interrupt signal")
			log.Info("stopping mouse mover")
			return
		}
	}
}

func (am *ActivityManager) handleMouse() {
	// Find the pattern
	var pattern models.Pattern
	for _, activity := range am.activities {
		if activity.Kind == models.KindMouse {
			pattern = activity.Pattern
			break
		}
	}

	err := am.mover.Move(pattern)
	if err != nil {
		if errors.Is(err, mouse.ErrUserInterruption) {
			log.Info("user interrupted movement")
		} else if errors.Is(err, mouse.ErrUserAlreadyMoving) {
			log.Info("user is already moving the cursor")
		} else {
			log.Error("error moving mouse", zap.Error(err))
		}
	} else {
		log.Info("mouse moved", zap.Time("at", time.Now()), zap.String("pattern", string(pattern)))
	}
}

func (am *ActivityManager) handleKeyboard() {
	robotgo.KeyTap("cmd")
	time.Sleep(500 * time.Millisecond)
	robotgo.KeyTap("v", "cmd")
}

func (am *ActivityManager) pause() {
	am.mutex.Lock()
	defer am.mutex.Unlock()

	if am.isPaused {
		log.Debug("activity manager is already paused")
		return
	}

	log.Debug("pausing activity manager")
	am.isPaused = true
	am.pausedAt = time.Now()
}

func (am *ActivityManager) resume() {
	ticker := time.NewTicker(1 * time.Second) // Check every second
	defer ticker.Stop()

	for range ticker.C {
		am.checkAndResumeIfNeeded()
	}
}

func (am *ActivityManager) checkAndPauseIfNeeded() {
	// Only pause if the behavior is enabled and we're not already paused
	if !am.behavior.PauseWhenUserIsActive {
		return
	}

	am.mutex.RLock()
	if am.isPaused {
		am.mutex.RUnlock()
		return
	}
	am.mutex.RUnlock()

	// Check if user is currently active
	if am.watcher.IsUserMoving() {
		log.Debug("user is active, pausing activity manager")
		am.pause()
	}
}

func (am *ActivityManager) checkAndResumeIfNeeded() {
	am.mutex.Lock()
	defer am.mutex.Unlock()

	if !am.isPaused {
		return
	}

	// Check if we've been paused longer than the idle timeout
	if time.Since(am.pausedAt) > am.behavior.IdleTimeout {
		// Check if user is still active
		if am.watcher.IsUserMoving() {
			log.Debug("user is still active, not resuming")
			// Reset the pause time to avoid constant checking
			am.pausedAt = time.Now()
			return
		}

		log.Debug("idle timeout reached, resuming activity manager")
		am.isPaused = false
		am.pausedAt = time.Time{} // Reset pause time
	}
}
