package mouse

import (
	"errors"
	"time"

	"github.com/go-vgo/robotgo"
)

var (
	// ErrDurationTooLong is returned when the movement duration exceeds the interval
	ErrDurationTooLong = errors.New("movement duration exceeds interval")

	// ErrUserInterruption is returned when user moves the cursor during execution
	ErrUserInterruption = errors.New("movement interrupted by user cursor activity")

	// ErrUserAlreadyMoving is returned when the user is already moving the cursor.
	ErrUserAlreadyMoving = errors.New("user is already moving the cursor")
)

const (
	Steps     int           = 10
	Delay     time.Duration = 50 * time.Millisecond
	Tolerance int           = 5
)

type Mover struct {
	distance int
	interval time.Duration
}

func NewMover(distance int, interval time.Duration) *Mover {
	return &Mover{
		distance: distance,
		interval: interval,
	}
}

func (m *Mover) EstimatedDuration() time.Duration {
	return time.Duration(4*Steps) * Delay
}

// Move moves the mouse cursor in a square pattern by the specified distance on Windows
func (m *Mover) Move() error {
	// Store original position for comparison
	originalX, originalY := robotgo.Location()

	// Check if user is already moving the cursor (compare position after small delay)
	isMoving, err := userIsAlreadyMoving(originalX, originalY)
	if err != nil {
		return err
	}
	if isMoving {
		return ErrUserAlreadyMoving
	}

	// Draw square by tracing each edge: right, down, left, up
	robotgo.MouseSleep = 100                                           // Make it smooth
	robotgo.Move(originalX+int(m.distance), originalY)                 // Right edge
	robotgo.Move(originalX+int(m.distance), originalY+int(m.distance)) // Down edge
	robotgo.Move(originalX, originalY+int(m.distance))                 // Left edge
	robotgo.Move(originalX, originalY)                                 // Up edge (back to start)

	return nil
}

func userIsAlreadyMoving(originalX, originalY int) (bool, error) {
	time.Sleep(10 * time.Millisecond)

	nowX, nowY := robotgo.Location()

	// If cursor moved, user is actively using it
	if nowX != originalX || nowY != originalY {
		return true, nil
	}

	return false, nil
}
