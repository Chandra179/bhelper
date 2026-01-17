package collision

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	Format   string
	Length   int
	Rate     int64
	RateUnit string
}

func ParseInput(input string) (*Config, error) {
	parts := strings.Split(input, ":")
	if len(parts) != 3 {
		return nil, errors.New("invalid format: expected 'format:length:rate/unit'")
	}

	format := parts[0]
	length, err := strconv.Atoi(parts[1])
	if err != nil {
		return nil, fmt.Errorf("invalid length: %v", err)
	}

	rateParts := strings.Split(parts[2], "/")
	if len(rateParts) != 2 {
		return nil, errors.New("invalid rate format: expected 'rate/unit'")
	}

	rate, err := strconv.ParseInt(rateParts[0], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid rate: %v", err)
	}
	rateUnit := rateParts[1]

	return &Config{
		Format:   format,
		Length:   length,
		Rate:     rate,
		RateUnit: rateUnit,
	}, nil
}

func parseDuration(s string) (time.Duration, error) {
	if len(s) < 2 {
		return 0, errors.New("duration too short")
	}

	var numStr, unit string
	if strings.HasSuffix(s, "day") || strings.HasSuffix(s, "days") {
		numStr = strings.TrimSuffix(s, "day")
		numStr = strings.TrimSuffix(numStr, "s")
		unit = "day"
	} else if strings.HasSuffix(s, "hour") || strings.HasSuffix(s, "hours") {
		numStr = strings.TrimSuffix(s, "hour")
		numStr = strings.TrimSuffix(numStr, "s")
		unit = "hour"
	} else if strings.HasSuffix(s, "year") || strings.HasSuffix(s, "years") {
		numStr = strings.TrimSuffix(s, "year")
		numStr = strings.TrimSuffix(numStr, "s")
		unit = "year"
	} else {
		numStr = s[:len(s)-1]
		unit = s[len(s)-1:]
	}

	num, err := strconv.Atoi(numStr)
	if err != nil {
		return 0, fmt.Errorf("invalid duration number: %v", err)
	}

	switch unit {
	case "s", "sec":
		return time.Duration(num) * time.Second, nil
	case "m", "min":
		return time.Duration(num) * time.Minute, nil
	case "h", "hour":
		return time.Duration(num) * time.Hour, nil
	case "d", "day":
		return time.Duration(num) * 24 * time.Hour, nil
	case "y", "year":
		return time.Duration(num) * 365 * 24 * time.Hour, nil
	default:
		return 0, fmt.Errorf("unknown duration unit: %s", unit)
	}
}
