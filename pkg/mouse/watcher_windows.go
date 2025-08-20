package mouse

import (
	"time"

	"github.com/go-vgo/robotgo"
)

const (
	DefaultWatcherInterval = 500 * time.Millisecond
)

type Watcher struct {
	ticker *time.Ticker

	actualX int
	actualY int

	isMoving bool
}

func NewWatcher() Watcher {
	return Watcher{
		ticker: time.NewTicker(DefaultWatcherInterval),
	}
}

func (w Watcher) Start() {
	go func() {
		w.actualX, w.actualY = robotgo.Location()

		for range w.ticker.C {
			// Check mouse position and take action if needed
			newX, newY := robotgo.Location()

			// Check if the mouse has moved
			if w.actualX != newX || w.actualY != newY {
				w.isMoving = true
			} else {
				w.isMoving = false
			}

			// Update the actual position
			w.actualX, w.actualY = newX, newY
		}
	}()
}

func (w Watcher) Stop() {
	w.ticker.Stop()
}

func (w Watcher) IsUserMoving() bool {
	return w.isMoving
}
