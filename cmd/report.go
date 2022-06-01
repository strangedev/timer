package cmd

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"text/tabwriter"
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

		fmt.Printf("Summary of %v\n", wantedDate)
		fmt.Println("================================\n")
		summaryTable := tabwriter.NewWriter(os.Stdout, 0, 2, 1, ' ', 0)
		fmt.Fprintf(summaryTable, "%s\t%v\n", "Total time worked", report.TotalTimeWorked)
		fmt.Fprintf(summaryTable, "%s\t%v\n", "Started working at", report.StartedWorkingAt)
		fmt.Fprintf(summaryTable, "%s\t%v\n", "Stopped working at", report.StoppedWorkingAt)
		for _, activePeriod := range report.ActivePeriods {
			for _, note := range activePeriod.Notes {
				fmt.Fprintf(summaryTable, "%s\t%v\n", "Note", note)
			}
		}

		summaryTable.Flush()
		fmt.Println("\n\nActive periods")
		fmt.Println("--------------------------------")
		activePeriodsTable := tabwriter.NewWriter(os.Stdout, 0, 2, 1, ' ', 0)
		for i, activePeriod := range report.ActivePeriods {
			fmt.Fprintf(activePeriodsTable, "\n%s\t%v\n", "Period", i)
			fmt.Fprintf(activePeriodsTable, "%s\t%v\n", "Duration", activePeriod.Duration)
			fmt.Fprintf(activePeriodsTable, "%s\t%v\n", "Started at", activePeriod.Started)
			fmt.Fprintf(activePeriodsTable, "%s\t%v\n", "Started by", activePeriod.StartedBy)
			fmt.Fprintf(activePeriodsTable, "%s\t%v\n", "Ended at", activePeriod.Ended)
			fmt.Fprintf(activePeriodsTable, "%s\t%v\n", "Ended by", activePeriod.EndedBy)
			for _, note := range activePeriod.Notes {
				fmt.Fprintf(activePeriodsTable, "%s\t%v\n", "Note", note)
			}
		}
		activePeriodsTable.Flush()
		fmt.Println("\n\nBreaks")
		fmt.Println("--------------------------------")
		breaksTable := tabwriter.NewWriter(os.Stdout, 0, 2, 1, ' ', 0)
		for i, breakPeriod := range report.Breaks {
			fmt.Fprintf(breaksTable, "\n%s\t%v\n", "Break", i)
			fmt.Fprintf(breaksTable, "%s\t%v\n", "Duration", breakPeriod.Duration)
			fmt.Fprintf(breaksTable, "%s\t%v\n", "Started at", breakPeriod.Started)
			fmt.Fprintf(breaksTable, "%s\t%v\n", "Started by", breakPeriod.StartedBy)
			fmt.Fprintf(breaksTable, "%s\t%v\n", "Ended at", breakPeriod.Ended)
			fmt.Fprintf(breaksTable, "%s\t%v\n", "Ended by", breakPeriod.EndedBy)
		}
		breaksTable.Flush()
	},
}
