package time

import (
	"bhelper/feature"
	"fmt"
	"strconv"
	"strings"
)

type TimeConverter struct{}

type ConversionResult struct {
	Nanoseconds  float64
	Microseconds float64
	Milliseconds float64
	Seconds      float64
	Minutes      float64
	Hours        float64
}

func parseInput(input string) (float64, string, error) {
	input = strings.TrimSpace(input)
	if input == "" {
		return 0, "", fmt.Errorf("empty input")
	}

	input = strings.ToLower(input)

	unit := ""
	i := len(input)
	for i > 0 {
		c := input[i-1]
		if (c >= 'a' && c <= 'z') || c == 'µ' {
			i--
		} else {
			break
		}
	}

	if i == len(input) {
		return 0, "", fmt.Errorf("no unit found")
	}

	unit = input[i:]
	valueStr := strings.TrimSpace(input[:i])

	value, err := strconv.ParseFloat(valueStr, 64)
	if err != nil {
		return 0, "", fmt.Errorf("invalid number: %w", err)
	}

	validUnits := map[string]bool{
		"ns": true, "us": true, "µs": true,
		"ms": true, "s": true,
		"min": true, "m": true,
		"h": true, "hr": true,
	}

	if !validUnits[unit] {
		return 0, "", fmt.Errorf("invalid unit: %s", unit)
	}

	return value, unit, nil
}

func convertToAllUnits(value float64, unit string) ConversionResult {
	var seconds float64

	switch unit {
	case "ns":
		seconds = value / 1e9
	case "us", "µs":
		seconds = value / 1e6
	case "ms":
		seconds = value / 1e3
	case "s":
		seconds = value
	case "min", "m":
		seconds = value * 60
	case "h", "hr":
		seconds = value * 3600
	default:
		return ConversionResult{}
	}

	return ConversionResult{
		Nanoseconds:  seconds * 1e9,
		Microseconds: seconds * 1e6,
		Milliseconds: seconds * 1e3,
		Seconds:      seconds,
		Minutes:      seconds / 60,
		Hours:        seconds / 3600,
	}
}

func NewTimeConverter() *TimeConverter {
	return &TimeConverter{}
}

func (tc *TimeConverter) ID() string {
	return "time"
}

func (tc *TimeConverter) Name() string {
	return "Time Converter"
}

func (tc *TimeConverter) Description() string {
	return "Convert time values between units (nanosecond to hour)"
}

func (tc *TimeConverter) Help() string {
	return `Time Converter converts time values between different units:

• Nanoseconds (ns)
• Microseconds (µs/us)
• Milliseconds (ms)
• Seconds (s)
• Minutes (min/m)
• Hours (h/hr)

This is useful for understanding time relationships, converting between
different time scales, or debugging timing issues.`
}

func (tc *TimeConverter) Examples() []feature.Example {
	return []feature.Example{
		{Input: "100ms", Description: "Convert 100 milliseconds to all units"},
		{Input: "1s", Description: "Convert 1 second to all units"},
		{Input: "5min", Description: "Convert 5 minutes to all units"},
		{Input: "1.5h", Description: "Convert 1.5 hours to all units"},
		{Input: "1000ns", Description: "Convert 1000 nanoseconds to all units"},
	}
}

func (tc *TimeConverter) Execute(input string) (string, error) {
	input = strings.TrimSpace(input)
	if input == "" {
		return "Please provide a time value to convert (e.g., '100ms', '1s', '5min')", nil
	}

	value, unit, err := parseInput(input)
	if err != nil {
		return fmt.Sprintf("Error: %v", err), nil
	}

	result := convertToAllUnits(value, unit)

	var output strings.Builder
	output.WriteString(fmt.Sprintf("Time Conversions (%s):\n", input))
	output.WriteString(fmt.Sprintf("  Nanoseconds:  %s\n", formatNumber(result.Nanoseconds)))
	output.WriteString(fmt.Sprintf("  Microseconds: %s\n", formatNumber(result.Microseconds)))
	output.WriteString(fmt.Sprintf("  Milliseconds: %s\n", formatNumber(result.Milliseconds)))
	output.WriteString(fmt.Sprintf("  Seconds:      %s\n", formatNumber(result.Seconds)))
	output.WriteString(fmt.Sprintf("  Minutes:      %s\n", formatNumber(result.Minutes)))
	output.WriteString(fmt.Sprintf("  Hours:        %s\n", formatNumber(result.Hours)))

	return output.String(), nil
}

func formatNumber(n float64) string {
	if n >= 1e6 {
		return fmt.Sprintf("%.0f", n)
	}
	if n < 0.001 && n != 0 {
		return fmt.Sprintf("%.6g", n)
	}
	if n >= 1 {
		return fmt.Sprintf("%.0f", n)
	}
	if n >= 0.001 {
		return fmt.Sprintf("%.6f", n)
	}
	return fmt.Sprintf("%.6g", n)
}
