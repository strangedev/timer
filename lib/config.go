package lib

import "time"

type Config struct {
	StorageDir          string
	CheckpointInterval  time.Duration
	IdleDuration        time.Duration
	IdleWatcherInterval time.Duration
	WakeWatcherInterval time.Duration
}
