package lib

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"
	"time"
)

type Timer struct {
	pastSlices     []TimeSlice
	currentSlice   *TimeSlice
	isStopped      bool
	isSuspended    bool
	suspendedSince time.Time
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

func (t *Timer) SuspendedFor() (time.Duration, error) {
	if !t.isSuspended {
		return 0, &TimerIsNotSuspended{}
	}

	return time.Now().Sub(t.suspendedSince), nil
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
	t.suspendedSince = time.Now()

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

func (t *Timer) Checkpoint(config Config) error {
	timerWasStarted := t.currentSlice != nil

	t.EndCurrentSliceIfExists("Checkpoint")

	data, err := json.MarshalIndent(t.pastSlices, "", "  ")
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(filepath.Join(config.StorageDir, time.Now().Local().Format("2006-01-02T15:04:05-0700")+".json"), data, 0644)
	if err != nil {
		return err
	}

	if timerWasStarted {
		t.StartNewSliceIfNotExists("Resume after checkpoint")
	}

	return nil
}

type VoidArgs struct{}

type VoidReply struct{}

type StartStopArgs struct {
	Reason string
}

type NoteArgs struct {
	Note string
}

func (t *Timer) RPC_Start(args *StartStopArgs, reply *VoidReply) error {
	return t.Start(args.Reason)
}

func (t *Timer) RPC_Stop(args *StartStopArgs, reply *VoidReply) error {
	return t.Stop(args.Reason)
}

func (t *Timer) RPC_Suspend(args *StartStopArgs, reply *VoidReply) error {
	return t.Suspend(args.Reason)
}

func (t *Timer) RPC_Wake(args *StartStopArgs, reply *VoidReply) error {
	return t.Wake(args.Reason)
}

func (t *Timer) GetPastSlices(args *VoidArgs, reply *[]TimeSlice) error {
	*reply = t.pastSlices

	return nil
}

func (t *Timer) RPC_Note(args *NoteArgs, reply *VoidReply) error {
	t.Wake("Note added")

	if t.IsStopped() {
		return &TimerIsNotYetStarted{}
	}

	t.currentSlice.Notes = append(t.currentSlice.Notes, args.Note)

	return nil
}
