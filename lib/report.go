package lib

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"
	"time"
)

type DailyReport struct {
	TotalTimeWorked  time.Duration
	StartedWorkingAt time.Time
	StoppedWorkingAt time.Time
	ActivePeriods    []TimeSlice
	Breaks           []TimeSlice
}

func GenerateReportFor(slices []TimeSlice) DailyReport {
	report := DailyReport{
		ActivePeriods: slices,
	}
	var CurrentBreak *TimeSlice

	for _, slice := range slices {
		report.TotalTimeWorked = report.TotalTimeWorked + slice.Duration

		if report.StartedWorkingAt.IsZero() {
			report.StartedWorkingAt = slice.Started
		}
		if report.StoppedWorkingAt.IsZero() {
			report.StoppedWorkingAt = slice.Ended
		}

		if slice.Started.Before(report.StartedWorkingAt) {
			report.StartedWorkingAt = slice.Started
		}
		if slice.Ended.After(report.StoppedWorkingAt) {
			report.StoppedWorkingAt = slice.Ended
		}

		if CurrentBreak == nil {
			CurrentBreak = &TimeSlice{
				Started:   slice.Ended,
				StartedBy: slice.EndedBy,
				Notes:     []string{"Break"},
			}
		} else {
			CurrentBreak.Ended = slice.Started
			CurrentBreak.EndedBy = slice.StartedBy
			CurrentBreak.Duration = CurrentBreak.Ended.Sub(CurrentBreak.Started)

			report.Breaks = append(report.Breaks, *CurrentBreak)
			CurrentBreak = nil
		}
	}

	return report
}

func GetTimeSlicesForDay(config Config, wantedDay time.Time) ([]TimeSlice, error) {
	year, month, day := wantedDay.Local().Date()
	startOfDay := time.Date(year, month, day, 0, 0, 0, 0, time.Local)
	filePrefix := startOfDay.Local().Format("2006-01-02")

	matches, err := filepath.Glob(filepath.Join(config.StorageDir, filePrefix+"*.json"))
	if err != nil {
		return []TimeSlice{}, err
	}

	totalSlices := make([]TimeSlice, 0)
	for _, fileName := range matches {
		data, err := ioutil.ReadFile(fileName)
		if err != nil {
			return totalSlices, err
		}
		slices := []TimeSlice{}
		if err := json.Unmarshal(data, &slices); err != nil {
			return totalSlices, err
		}
		totalSlices = append(totalSlices, slices...)
	}

	return totalSlices, nil
}
