package lib

import "time"

type TimeSlice struct {
	Started   time.Time
	Ended     time.Time
	Duration  time.Duration
	Notes     []string
	StartedBy string
	EndedBy   string
}
