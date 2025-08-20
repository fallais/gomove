package mouse

import (
	"errors"
	"fmt"
	"time"
	"unsafe"

	"golang.org/x/sys/windows"
)

var (
	// ErrDurationTooLong is returned when the movement duration exceeds the interval
	ErrDurationTooLong = errors.New("movement duration exceeds interval")
	// ErrUserInterruption is returned when user moves the cursor during execution
	ErrUserInterruption = errors.New("movement interrupted by user cursor activity")
)

// Move moves the mouse cursor in a square pattern by the specified distance on Windows
func Move(distance int, interval time.Duration) error {
	user32 := windows.NewLazySystemDLL("user32.dll")
	getCursorPos := user32.NewProc("GetCursorPos")
	setCursorPos := user32.NewProc("SetCursorPos")

	// Calculate estimated duration for the movement
	// Each edge has 10 steps with 50ms delay = 500ms per edge, 4 edges = 2000ms total
	estimatedDuration := time.Duration(4*10*50) * time.Millisecond
	if estimatedDuration >= interval {
		return ErrDurationTooLong
	}

	// Get current cursor position
	var point struct {
		X, Y int32
	}

	ret, _, err := getCursorPos.Call(uintptr(unsafe.Pointer(&point)))
	if ret == 0 {
		return fmt.Errorf("failed to get cursor position: %v", err)
	}

	// Store original position for comparison
	originalX, originalY := point.X, point.Y

	// Helper function to check if user moved the cursor
	checkUserMovement := func() error {
		var currentPoint struct {
			X, Y int32
		}
		ret, _, err := getCursorPos.Call(uintptr(unsafe.Pointer(&currentPoint)))
		if ret == 0 {
			return fmt.Errorf("failed to get cursor position: %v", err)
		}
		return nil
	}

	// Check if user is already moving the cursor (compare position after small delay)
	time.Sleep(10 * time.Millisecond)
	var checkPoint struct {
		X, Y int32
	}
	ret, _, err = getCursorPos.Call(uintptr(unsafe.Pointer(&checkPoint)))
	if ret == 0 {
		return fmt.Errorf("failed to get cursor position: %v", err)
	}

	// If cursor moved, user is actively using it
	if checkPoint.X != originalX || checkPoint.Y != originalY {
		return nil // Do nothing if user is already moving the cursor
	}

	// Store the expected position we're controlling
	expectedX, expectedY := originalX, originalY

	// Helper function for absolute value
	abs := func(x int32) int32 {
		if x < 0 {
			return -x
		}
		return x
	}

	// Helper function to smoothly move from one point to another with interruption detection
	smoothMove := func(startX, startY, endX, endY int32) error {
		steps := int32(10) // Number of steps to take between points
		deltaX := (endX - startX) / steps
		deltaY := (endY - startY) / steps

		for i := int32(0); i <= steps; i++ {
			// Before moving, check if user has moved the cursor unexpectedly
			if err := checkUserMovement(); err != nil {
				return err
			}

			var currentPoint struct {
				X, Y int32
			}
			ret, _, err := getCursorPos.Call(uintptr(unsafe.Pointer(&currentPoint)))
			if ret == 0 {
				return fmt.Errorf("failed to get cursor position: %v", err)
			}

			// If the cursor is not where we expect it to be, user has moved it
			tolerance := int32(5) // Allow small tolerance for rounding errors
			if abs(currentPoint.X-expectedX) > tolerance || abs(currentPoint.Y-expectedY) > tolerance {
				return ErrUserInterruption
			}

			currentX := startX + deltaX*i
			currentY := startY + deltaY*i

			ret, _, err = setCursorPos.Call(uintptr(currentX), uintptr(currentY))
			if ret == 0 {
				return fmt.Errorf("failed to set cursor position: %v", err)
			}

			// Update expected position
			expectedX, expectedY = currentX, currentY

			time.Sleep(50 * time.Millisecond) // Small delay for smooth movement
		}
		return nil
	}

	// Draw square by tracing each edge: right, down, left, up
	// Right edge
	if err := smoothMove(originalX, originalY, originalX+int32(distance), originalY); err != nil {
		return err
	}

	// Down edge
	if err := smoothMove(originalX+int32(distance), originalY, originalX+int32(distance), originalY+int32(distance)); err != nil {
		return err
	}

	// Left edge
	if err := smoothMove(originalX+int32(distance), originalY+int32(distance), originalX, originalY+int32(distance)); err != nil {
		return err
	}

	// Up edge (back to start)
	if err := smoothMove(originalX, originalY+int32(distance), originalX, originalY); err != nil {
		return err
	}

	return nil
}
