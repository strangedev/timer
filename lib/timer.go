package lib

import (
	"time"
)

type Timer struct {
	pastSlices   []TimeSlice
	currentSlice *TimeSlice
	isStopped    bool
	isSuspended  bool
	suspendedAt  time.Time
}

func CreateTimer() *Timer {
	return &Timer{
		pastSlices: make([]TimeSlice, 0),
		isStopped:  true,
	}
}

func (t *Timer) IsSuspended() bool {
	return t.isSuspended
}

func (t *Timer) IsStarted() bool {
	return !t.IsStopped()
}
func (t *Timer) IsStopped() bool {
	return t.isStopped
}

func (t *Timer) SuspendedSince() time.Duration {
	if !t.isSuspended {
		return 0
	}

	return time.Since(t.suspendedAt)
}

func (t *Timer) StartNewSliceIfNotExists(reason string) {
	if t.currentSlice != nil {
		return
	}

	t.currentSlice = &TimeSlice{
		Started:   time.Now(),
		StartedBy: reason,
	}
}

func (t *Timer) EndCurrentSliceIfExists(reason string) {
	if t.currentSlice == nil {
		return
	}

	t.currentSlice.Ended = time.Now()
	t.currentSlice.EndedBy = reason
	t.currentSlice.Duration = t.currentSlice.Ended.Sub(t.currentSlice.Started)
	t.pastSlices = append(t.pastSlices, *t.currentSlice)
	t.currentSlice = nil
}

func (t *Timer) Start(reason string) error {
	if t.IsStarted() {
		return &TimerIsAlreadyStarted{}
	}

	t.StartNewSliceIfNotExists(reason)
	t.isStopped = false
	t.isSuspended = false

	return nil
}

func (t *Timer) Stop(reason string) error {
	if t.IsStopped() {
		return &TimerIsNotYetStarted{}
	}

	t.EndCurrentSliceIfExists(reason)
	t.isStopped = true

	return nil
}

func (t *Timer) Suspend(reason string) error {
	if !t.IsStarted() {
		return &TimerIsNotYetStarted{}
	}
	if t.isSuspended {
		return &TimerIsAlreadySuspended{}
	}

	t.EndCurrentSliceIfExists(reason)
	t.isSuspended = true
	t.suspendedAt = time.Now()

	return nil
}

func (t *Timer) Wake(reason string) error {
	if !t.IsStarted() {
		return &TimerIsNotYetStarted{}
	}
	if !t.isSuspended {
		return &TimerIsNotSuspended{}
	}

	t.StartNewSliceIfNotExists(reason)
	t.isSuspended = false

	return nil
}

func (t *Timer) ContinueWithNewSlice(reason string) {
	timerWasStarted := t.currentSlice != nil

	t.EndCurrentSliceIfExists(reason)

	if timerWasStarted {
		t.StartNewSliceIfNotExists(reason)
	}
}

func (t *Timer) Note(note string) error {
	t.Wake("Note added")

	if t.IsStopped() {
		return &TimerIsNotYetStarted{}
	}

	t.currentSlice.Notes = append(t.currentSlice.Notes, note)

	return nil
}

func (t *Timer) PastSlices() []TimeSlice {
	return t.pastSlices
}

func (t *Timer) ClearPastSlices() {
	t.pastSlices = make([]TimeSlice, 0)
}
