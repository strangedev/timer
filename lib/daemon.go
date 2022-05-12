package lib

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"path/filepath"
	"time"
)

type Daemon struct {
	*Timer
	StartedAt        time.Time
	Config           Config
	LastCheckpointAt time.Time
	NextCheckpointAt time.Time
}

func CreateDaemon(config Config) Daemon {
	return Daemon{
		Timer:     CreateTimer(),
		StartedAt: time.Now(),
		Config:    config,
	}
}

type VoidArgs struct{}

type VoidReply struct{}

type StartStopArgs struct {
	Reason string
}

type NoteArgs struct {
	Note string
}

type StatusReply struct {
	IsTimerStarted      bool
	IsTimerSuspended    bool
	TimerSuspendedSince time.Duration
	DaemonStartedAt     time.Time
	LastCheckpointAt    time.Time
	NextCheckpointAt    time.Time
}

func (d *Daemon) Checkpoint(reason string) error {
	d.ContinueWithNewSlice("Checkpoint: " + reason)

	data, err := json.MarshalIndent(d.PastSlices(), "", "  ")
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(filepath.Join(d.Config.StorageDir, time.Now().Local().Format("2006-01-02T15:04:05-0700")+".json"), data, 0644)
	if err != nil {
		return err
	}

	d.LastCheckpointAt = time.Now()
	d.ClearPastSlices()

	return nil
}

func (d *Daemon) StartTimer(args *StartStopArgs, reply *VoidReply) error {
	return d.Start(args.Reason)
}

func (d *Daemon) StopTimer(args *StartStopArgs, reply *VoidReply) error {
	return d.Stop(args.Reason)
}

func (d *Daemon) SuspendTimer(args *StartStopArgs, reply *VoidReply) error {
	return d.Suspend(args.Reason)
}

func (d *Daemon) WakeTimer(args *StartStopArgs, reply *VoidReply) error {
	return d.Wake(args.Reason)
}

func (d *Daemon) GetStatus(args *VoidArgs, reply *StatusReply) error {

	*reply = StatusReply{
		IsTimerStarted:      d.IsStarted(),
		IsTimerSuspended:    d.IsSuspended(),
		TimerSuspendedSince: d.SuspendedSince(),
		DaemonStartedAt:     d.StartedAt,
		LastCheckpointAt:    d.LastCheckpointAt,
		NextCheckpointAt:    d.NextCheckpointAt,
	}

	return nil
}

func (d *Daemon) TriggerCheckpoint(args *VoidArgs, reply *VoidReply) error {
	return d.Checkpoint("Triggered manually")
}

func (d *Daemon) AddNote(args *NoteArgs, reply *VoidReply) error {
	return d.Note(args.Note)
}

func (d *Daemon) Run(stop chan bool) {
	x11 := NewX11()
	if err := x11.Init(); err != nil {
		log.Fatal(err)
	}

	rpc.Register(d)
	rpc.HandleHTTP()

	socket, err := net.Listen("tcp", "localhost:1234")
	if err != nil {
		log.Fatal("listen error:", err)
	}

	go func() {
		fmt.Println("Starting RPC server")
		go http.Serve(socket, nil)
	}()

	go func() {
		fmt.Println("Starting checkpointing routine")
		for {
			d.NextCheckpointAt = time.Now().Add(d.Config.CheckpointInterval)
			time.Sleep(d.Config.CheckpointInterval)
			fmt.Println("Checkpointing")
			if checkpointErr := d.Checkpoint("Triggered by schedule"); checkpointErr != nil {
				fmt.Printf("%+v\n", d)
				log.Fatal(checkpointErr)
			}
		}
	}()

	go func() {
		fmt.Println("Starting X11 idle watcher")
		for {
			time.Sleep(d.Config.IdleWatcherInterval)
			if d.IsStopped() {
				continue
			}

			idleTime, err := x11.GetIdleTime()
			if err != nil {
				log.Fatal(err)
			}

			if idleTime > d.Config.IdleDuration {
				if !d.IsSuspended() {
					fmt.Println("Suspending timer due to X11 idle")
					d.Suspend("X11 idle")
				}
			}
		}
	}()
	go func() {
		fmt.Println("Starting X11 wake watcher")
		for {
			time.Sleep(d.Config.WakeWatcherInterval)
			if d.IsStopped() {
				continue
			}

			currentIdleTime, err := x11.GetIdleTime()
			if err != nil {
				log.Fatal(err)
			}

			if currentIdleTime < d.Config.IdleDuration {
				if d.IsSuspended() {
					fmt.Println("Waking timer due to X11 activity")
					d.Wake("X11 activity")
				}
			}
		}
	}()

	<-stop
	if checkpointErr := d.Checkpoint("Triggered by exit"); checkpointErr != nil {
		fmt.Printf("%+v\n", d)
		log.Fatal(checkpointErr)
	}
}
