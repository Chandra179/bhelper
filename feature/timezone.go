package feature

import (
	"fmt"
	"strings"
	"time"
)

// TimezoneAnalyzer provides comprehensive timezone and time information
type TimezoneAnalyzer struct{}

func NewTimezoneAnalyzer() *TimezoneAnalyzer {
	return &TimezoneAnalyzer{}
}

func (ta *TimezoneAnalyzer) ID() string {
	return "timezone"
}

func (ta *TimezoneAnalyzer) Name() string {
	return "Timezone Analyzer"
}

func (ta *TimezoneAnalyzer) Description() string {
	return "Get Unix timestamp for a date (dd-mm-yyyy)"
}

func (ta *TimezoneAnalyzer) Help() string {
	return `Timezone Analyzer displays Unix timestamp and timezone information for a specific date.

Input Format:
  • dd-mm-yyyy (e.g., 16-01-2026 for January 16, 2026)

Outputs:
  • Unix timestamp (seconds, milliseconds, microseconds, nanoseconds)
  • Date and time information
  • Timezone offset and UTC time
  • Daylight saving status
  • Julian day number and season

Use this feature to convert dates to Unix timestamps for programming, debugging,
or understanding time zone behavior.`
}

func (ta *TimezoneAnalyzer) Examples() []Example {
	return []Example{
		{Input: "16-01-2026", Description: "Get Unix timestamp for January 16, 2026"},
		{Input: "01-01-2024", Description: "Get Unix timestamp for January 1, 2024"},
		{Input: "25-12-2025", Description: "Get Unix timestamp for December 25, 2025"},
	}
}

func (ta *TimezoneAnalyzer) Execute(input string) (string, error) {
	targetTime, err := ta.parseDate(input)
	if err != nil {
		return "", err
	}

	localZone, _ := targetTime.Zone()

	var result strings.Builder

	result.WriteString(ta.renderTimeFormats(targetTime, localZone) + "\n")

	return result.String(), nil
}

func (ta *TimezoneAnalyzer) parseDate(input string) (time.Time, error) {
	if input == "" {
		return time.Now(), nil
	}

	parsed, err := time.Parse("02-01-2006", input)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid date format. Use dd-mm-yyyy (e.g., 16-01-2026)")
	}

	return parsed, nil
}

func (ta *TimezoneAnalyzer) renderTimeFormats(now time.Time, zoneName string) string {
	var b strings.Builder

	b.WriteString(fmt.Sprintf("Date:          %s\n", now.Format("2006-01-02")))
	b.WriteString(fmt.Sprintf("Time:          %s %s\n", now.Format("15:04:05"), ta.formatOffset(now)))
	b.WriteString(fmt.Sprintf("UTC Time:      %s UTC\n", now.UTC().Format("15:04:05")))
	b.WriteString(fmt.Sprintf("Unix Timestamp: %d\n", now.Unix()))

	dayOfWeek := now.Weekday()
	dayOfYear := ta.dayOfYear(now)
	weekNumber := ta.isoWeekNumber(now)
	isLeap := ta.isLeapYear(now.Year())

	b.WriteString(fmt.Sprintf("Day of Week:      %s\n", dayOfWeek.String()))
	b.WriteString(fmt.Sprintf("Day of Year:      %d/%d\n", dayOfYear, ta.daysInYear(now.Year())))
	b.WriteString(fmt.Sprintf("ISO Week Number:  %d\n", weekNumber))
	b.WriteString(fmt.Sprintf("Julian Day:       %d\n", ta.julianDay(now)))
	b.WriteString(fmt.Sprintf("Season:           %s\n", ta.season(now)))
	b.WriteString(fmt.Sprintf("Leap Year:        %t\n", isLeap))

	return b.String()
}

func (ta *TimezoneAnalyzer) formatOffset(t time.Time) string {
	_, offset := t.Zone()
	hours := offset / 3600
	minutes := (offset % 3600) / 60

	if hours >= 0 {
		return fmt.Sprintf("+%02d:%02d", hours, minutes)
	}
	return fmt.Sprintf("%03d:%02d", hours, minutes)
}

func (ta *TimezoneAnalyzer) dayOfYear(t time.Time) int {
	start := time.Date(t.Year(), 1, 1, 0, 0, 0, 0, t.Location())
	duration := t.Sub(start)
	return int(duration.Hours()/24) + 1
}

func (ta *TimezoneAnalyzer) daysInYear(year int) int {
	if ta.isLeapYear(year) {
		return 366
	}
	return 365
}

func (ta *TimezoneAnalyzer) isoWeekNumber(t time.Time) int {
	year, week := t.ISOWeek()
	if year != t.Year() {
		return 0
	}
	return week
}

func (ta *TimezoneAnalyzer) isLeapYear(year int) bool {
	if year%4 != 0 {
		return false
	} else if year%100 != 0 {
		return true
	} else {
		return year%400 == 0
	}
}

func (ta *TimezoneAnalyzer) julianDay(t time.Time) int {
	year := t.Year()
	month := int(t.Month())
	day := t.Day()

	a := (14 - month) / 12
	y := year + 4800 - a
	m := month + 12*a - 3

	julianDayNumber := day + (153*m+2)/5 + 365*y + y/4 - y/100 + y/400 - 32045

	return julianDayNumber
}

func (ta *TimezoneAnalyzer) season(t time.Time) string {
	month := int(t.Month())

	switch {
	case month >= 3 && month <= 5:
		return "Spring"
	case month >= 6 && month <= 8:
		return "Summer"
	case month >= 9 && month <= 11:
		return "Autumn"
	default:
		return "Winter"
	}
}
