package calendarManager

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/mitchs-dev/library-go/generator"
	"github.com/mitchs-dev/worklog/internal/configuration"
	log "github.com/sirupsen/logrus"
)

// This file is used to fetch the periods from the calendar

var (
	// validPeriods is a map of valid periods
	validPeriods = []string{
		"today",
		"yesterday",
		"3day",
		"week",
		"cweek",
		"month",
		"quarter",
		"year",
	}

	// periodMap is a map of periods to their respective start and end dates
	periodMap = map[string]int{
		"today":     0,
		"yesterday": -1,
		"3day":      -2,
		"cweek":     0,
		"week":      -7,
		"month":     -30,
		"quarter":   -90,
		"year":      -365,
	}
)

// validPeriod checks if the period is valid
func validPeriod(period string) bool {
	for _, validPeriod := range validPeriods {
		if validPeriod == strings.ToLower(period) {
			return true
		}
	}
	return false
}

// PeriodFetch fetches the period from the calendar and returns all of the weeks in the period in the format "YYYY/WW", the month/days in the period in the format MMDD, and the first and last day of the period in the format of MMDD
func PeriodFetch(period string) ([]string, YearTree, string, string, error) {
	if !validPeriod(period) {
		return nil, YearTree{}, "", "", errors.New("invalid period")
	}

	// Get the current date
	epochTimestamp := int64(generator.EpochTimestamp(configuration.ScheduleWorkdayTimezone))

	// Convert epoch to time.Time
	endDate := time.Unix(epochTimestamp, 0)

	// Get the current year and week (end date)
	endYear, endWeek := endDate.ISOWeek()

	// Subtract the period from the current date
	startDate := endDate.AddDate(0, 0, periodMap[period])

	// Get the year and week of the start date
	startYear, startWeek := startDate.ISOWeek()

	// Get the day and month of the start and end dates
	startDay := startDate.Format("0102")

	// Get the day and month of the end date
	endDay := endDate.Format("0102")

	// Create a slice of weeks
	weeks := []string{}
	// Build the year tree
	// Initialize YearTree properly
	yearTree := YearTree{
		Years: make(map[int]WeekTree),
	}
	for year := startYear; year <= endYear; year++ {
		// Get the start and end week
		start := 1
		end := 52
		if year == startYear {
			start = startWeek
		}
		if year == endYear {
			end = endWeek
		}

		// Add the weeks to the slice
		// Add the year to the year tree
		log.Debug("Adding year to year tree: ", year)
		yearTree.Years[year] = WeekTree{
			Weeks: make(map[int]MonthDayTree),
		}
		for week := start; week <= end; week++ {
			log.Debug("Adding week to year tree: ", week)
			yearTree.Years[year].Weeks[week] = MonthDayTree{
				MonthDays: make([]string, 0),
			}

			// Add the week to the slice
			weeks = append(weeks, fmt.Sprintf("%d/%02d", year, week))

			// Calculate week start/end dates based on period
			weekStartDate := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)
			now := time.Now().Local()

			var weekEndDate time.Time

			switch period {
			case "today":
				weekStartDate = now
				weekEndDate = now
			case "yesterday":
				weekStartDate = now.AddDate(0, 0, -1)
				weekEndDate = now.AddDate(0, 0, -1)
			case "3day":
				weekStartDate = now.AddDate(0, 0, -2)
				weekEndDate = now
			default:
				// Existing weekly logic
				startDay := parseWeekday(configuration.ScheduleDaysStart)
				endDay := parseWeekday(configuration.ScheduleDaysEnd)

				// Find first occurrence of start day
				for weekStartDate.Weekday() != startDay {
					weekStartDate = weekStartDate.AddDate(0, 0, 1)
				}

				// Calculate correct week offset
				if time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC).Weekday() > startDay {
					weekStartDate = weekStartDate.AddDate(0, 0, (week-2)*7)
				} else {
					weekStartDate = weekStartDate.AddDate(0, 0, (week-1)*7)
				}

				daysInPeriod := int(endDay-startDay) + 1
				if daysInPeriod < 0 {
					daysInPeriod += 7
				}
				weekEndDate = weekStartDate.AddDate(0, 0, daysInPeriod-1)
			}

			// Populate dates
			currentDate := weekStartDate
			for currentDate.Before(weekEndDate) || currentDate.Equal(weekEndDate) {
				monthDay := currentDate.Format("0102")
				weekPtr := yearTree.Years[year].Weeks[week]
				weekPtr.MonthDays = append(weekPtr.MonthDays, monthDay)
				yearTree.Years[year].Weeks[week] = weekPtr
				currentDate = currentDate.AddDate(0, 0, 1)
			}
			// ...existing code...

		}
	}

	return weeks, yearTree, startDay, endDay, nil
}

// parseWeekday parses the weekday string to time.Weekday
func parseWeekday(day string) time.Weekday {
	switch strings.ToLower(day) {
	case "sunday":
		return time.Sunday
	case "monday":
		return time.Monday
	case "tuesday":
		return time.Tuesday
	case "wednesday":
		return time.Wednesday
	case "thursday":
		return time.Thursday
	case "friday":
		return time.Friday
	case "saturday":
		return time.Saturday
	default:
		return time.Monday // Default to Monday if invalid
	}
}
