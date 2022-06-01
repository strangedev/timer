package lib

import "fmt"

type TimerIsAlreadyStarted struct{}

func (e *TimerIsAlreadyStarted) Error() string {
	return "The timer is already started."
}

type TimerIsNotYetStarted struct{}

func (e *TimerIsNotYetStarted) Error() string {
	return "The timer is not yet started."
}

type TimerIsAlreadySuspended struct{}

func (e *TimerIsAlreadySuspended) Error() string {
	return "The timer is already suspended."
}

type TimerIsNotSuspended struct{}

func (e *TimerIsNotSuspended) Error() string {
	return "The timer is not suspended."
}

type X11Error struct {
	Code int
}

func (e *X11Error) Error() string {
	return fmt.Sprintf("X11 call failed (%v)", e.Code)
}
