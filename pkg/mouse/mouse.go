package mouse

import (
	"errors"
	"gomove/internal/models"
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
	Distance  int           = 10
	Tolerance int           = 5
)

type Mover struct {
}

func NewMover() *Mover {
	return &Mover{}
}

func (m *Mover) EstimatedDuration() time.Duration {
	return time.Duration(4*Steps) * Delay
}

// Move moves the mouse cursor in a square pattern by the specified distance on Windows
func (m *Mover) Move(pattern models.Pattern) error {
	// Store original position for comparison
	originalX, originalY := robotgo.Location()

	robotgo.MouseSleep = 100
	switch pattern {
	case models.PatternSquare:
		m.moveSquare(originalX, originalY)
	case models.PatternTriangle:
		m.moveTriangle(originalX, originalY)
	case models.PatternUpAndDown:
		m.moveUpAndDown(originalX, originalY)
	case models.PatternLeftAndRight:
		m.moveLeftAndRight(originalX, originalY)
	}

	// Draw square by tracing each edge: right, down, left, up
	robotgo.MouseSleep = 100 // Make it smooth

	return nil
}

func (m *Mover) moveSquare(originalX, originalY int) {
	robotgo.Move(originalX+Distance, originalY)          // Right edge
	robotgo.Move(originalX+Distance, originalY+Distance) // Down edge
	robotgo.Move(originalX, originalY+Distance)          // Left edge
	robotgo.Move(originalX, originalY)                   // Up edge (back to start)
}

func (m *Mover) moveTriangle(originalX, originalY int) {
	robotgo.Move(originalX+Distance, originalY)            // Right edge
	robotgo.Move(originalX+Distance/2, originalY+Distance) // Down edge (to the middle)
	robotgo.Move(originalX, originalY)                     // Left edge (back to start)
}

func (m *Mover) moveUpAndDown(originalX, originalY int) {
	robotgo.Move(originalX, originalY+Distance) // Down edge
	robotgo.Move(originalX, originalY)          // Up edge (back to start)
}

func (m *Mover) moveLeftAndRight(originalX, originalY int) {
	robotgo.Move(originalX+Distance, originalY) // Right edge
	robotgo.Move(originalX, originalY)          // Left edge (back to start)
}
