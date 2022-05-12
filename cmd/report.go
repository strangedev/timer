package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"time"

	"github.com/spf13/cobra"
	"github.com/strangedev/worktimer/lib"
)

func init() {
	rootCmd.AddCommand(reportCommand)
}

var dateRegex = regexp.MustCompile(`(?P<Year>\d{4})-(?P<Month>\d{2})-(?P<Day>\d{2})`)

var reportCommand = &cobra.Command{
	Use:   "report",
	Short: "Generates a report for a given day.",
	Long:  "Generates a report for a given day. The date is given as YYYY-MM-DD.",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		params := make(map[string]string)
		matches := dateRegex.FindStringSubmatch(args[0])
		for i, name := range dateRegex.SubexpNames() {
			if i > 0 && i <= len(matches) {
				params[name] = matches[i]
			}
		}
		yearString, ok := params["Year"]
		if !ok {
			log.Fatal("Could not parse the given date: Year is missing.")
		}
		year, err := strconv.Atoi(yearString)
		if err != nil {
			log.Fatal("Could not parse the given date: Year is malformed.")
		}
		monthString, ok := params["Month"]
		if !ok {
			log.Fatal("Could not parse the given date: Month is missing.")
		}
		month, err := strconv.Atoi(monthString)
		if err != nil {
			log.Fatal("Could not parse the given date: Month is malformed.")
		}
		if month < 1 || month > 12 {
			log.Fatal("Could not parse the given date: Month is not between 1 and 12.")
		}
		dayString, ok := params["Day"]
		if !ok {
			log.Fatal("Could not parse the given date: Day is missing.")
		}
		day, err := strconv.Atoi(dayString)
		if err != nil {
			log.Fatal("Could not parse the given date: Day is malformed.")
		}
		if day < 1 || day > 31 {
			log.Fatal("Could not parse the given date: Day is not between 1 and 31.")
		}
		wantedDate := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.Local)

		timeSlices, err := lib.GetTimeSlicesForDay(config, wantedDate)
		if err != nil {
			log.Fatalf("Could not load time slices: %v", err)
		}

		report := lib.GenerateReportFor(timeSlices)
		reportJson, err := json.MarshalIndent(report, "", "  ")
		if err != nil {
			log.Fatalf("Could not marshal report: %v", err)
		}

		fmt.Println(string(reportJson))
	},
}
