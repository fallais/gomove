package watcher

import (
	"gomove/pkg/log"
	"sync"
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

	isMoving      bool
	isTappingKeys bool

	mutex sync.RWMutex
}

func NewWatcher() *Watcher {
	return &Watcher{
		ticker: time.NewTicker(DefaultWatcherInterval),
	}
}

func (w *Watcher) Start() {
	log.Debug("starting watcher")
	go w.watchMouse()
	// TODO: go w.watchKeyboard()
}

func (w *Watcher) Stop() {
	w.ticker.Stop()
}

func (w *Watcher) IsUserMoving() bool {
	w.mutex.RLock()
	defer w.mutex.RUnlock()
	return w.isMoving || w.isTappingKeys
}

func (w *Watcher) watchMouse() {
	w.actualX, w.actualY = robotgo.Location()

	for range w.ticker.C {
		// Check mouse position and take action if needed
		newX, newY := robotgo.Location()

		// Check if the mouse has moved
		if w.actualX != newX || w.actualY != newY {
			w.mutex.Lock()
			w.isMoving = true
			w.mutex.Unlock()
		} else {
			w.mutex.Lock()
			w.isMoving = false
			w.mutex.Unlock()
		}

		// Update the actual position
		w.actualX, w.actualY = newX, newY
	}
}
