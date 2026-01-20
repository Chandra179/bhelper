package time

import (
	"strings"
	"testing"
)

func TestTimeConverter_Execute(t *testing.T) {
	converter := NewTimeConverter()

	tests := []struct {
		name    string
		input   string
		wantErr bool
		check   func(t *testing.T, output string)
	}{
		{
			name:    "100 milliseconds",
			input:   "100 ms",
			wantErr: false,
			check: func(t *testing.T, output string) {
				if !strings.Contains(output, "100 ms") {
					t.Errorf("Output should contain input value")
				}
				if !strings.Contains(output, "Nanoseconds") {
					t.Errorf("Output should contain Nanoseconds")
				}
				if !strings.Contains(output, "100000000") {
					t.Errorf("Output should contain 100000000 nanoseconds")
				}
			},
		},
		{
			name:    "1 second",
			input:   "1s",
			wantErr: false,
			check: func(t *testing.T, output string) {
				if !strings.Contains(output, "1s") {
					t.Errorf("Output should contain input value")
				}
				if !strings.Contains(output, "1000000000") {
					t.Errorf("Output should contain 1000000000 nanoseconds")
				}
			},
		},
		{
			name:    "empty input",
			input:   "",
			wantErr: false,
			check: func(t *testing.T, output string) {
				if !strings.Contains(output, "Please provide") {
					t.Errorf("Output should contain error message")
				}
			},
		},
		{
			name:    "invalid input",
			input:   "invalid",
			wantErr: false,
			check: func(t *testing.T, output string) {
				if !strings.Contains(output, "Error") {
					t.Errorf("Output should contain error message")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output, err := converter.Execute(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Execute() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			tt.check(t, output)
		})
	}
}

func TestParseInput(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantValue float64
		wantUnit  string
		wantErr   bool
	}{
		{
			name:      "millisecond with space",
			input:     "100 ms",
			wantValue: 100,
			wantUnit:  "ms",
			wantErr:   false,
		},
		{
			name:      "millisecond without space",
			input:     "100ms",
			wantValue: 100,
			wantUnit:  "ms",
			wantErr:   false,
		},
		{
			name:      "second with space",
			input:     "1 s",
			wantValue: 1,
			wantUnit:  "s",
			wantErr:   false,
		},
		{
			name:      "second without space",
			input:     "1s",
			wantValue: 1,
			wantUnit:  "s",
			wantErr:   false,
		},
		{
			name:      "nanosecond",
			input:     "1000ns",
			wantValue: 1000,
			wantUnit:  "ns",
			wantErr:   false,
		},
		{
			name:      "microsecond",
			input:     "500us",
			wantValue: 500,
			wantUnit:  "us",
			wantErr:   false,
		},
		{
			name:      "minute",
			input:     "5min",
			wantValue: 5,
			wantUnit:  "min",
			wantErr:   false,
		},
		{
			name:      "hour",
			input:     "2h",
			wantValue: 2,
			wantUnit:  "h",
			wantErr:   false,
		},
		{
			name:      "decimal value",
			input:     "1.5s",
			wantValue: 1.5,
			wantUnit:  "s",
			wantErr:   false,
		},
		{
			name:    "empty input",
			input:   "",
			wantErr: true,
		},
		{
			name:    "invalid unit",
			input:   "100xyz",
			wantErr: true,
		},
		{
			name:    "invalid number",
			input:   "abcs",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			value, unit, err := parseInput(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseInput() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if value != tt.wantValue {
					t.Errorf("parseInput() value = %v, want %v", value, tt.wantValue)
				}
				if unit != tt.wantUnit {
					t.Errorf("parseInput() unit = %v, want %v", unit, tt.wantUnit)
				}
			}
		})
	}
}

func TestConvertToAllUnits(t *testing.T) {
	tests := []struct {
		name    string
		value   float64
		unit    string
		wantNS  float64
		wantUS  float64
		wantMS  float64
		wantS   float64
		wantMin float64
		wantH   float64
	}{
		{
			name:    "1 second",
			value:   1,
			unit:    "s",
			wantNS:  1e9,
			wantUS:  1e6,
			wantMS:  1e3,
			wantS:   1,
			wantMin: 1.0 / 60.0,
			wantH:   1.0 / 3600.0,
		},
		{
			name:    "100 milliseconds",
			value:   100,
			unit:    "ms",
			wantNS:  100 * 1e6,
			wantUS:  100 * 1e3,
			wantMS:  100,
			wantS:   0.1,
			wantMin: 0.1 / 60.0,
			wantH:   0.1 / 3600.0,
		},
		{
			name:    "60 minutes",
			value:   60,
			unit:    "min",
			wantNS:  60 * 60 * 1e9,
			wantUS:  60 * 60 * 1e6,
			wantMS:  60 * 60 * 1e3,
			wantS:   3600,
			wantMin: 60,
			wantH:   1,
		},
		{
			name:    "1 hour",
			value:   1,
			unit:    "h",
			wantNS:  3600 * 1e9,
			wantUS:  3600 * 1e6,
			wantMS:  3600 * 1e3,
			wantS:   3600,
			wantMin: 60,
			wantH:   1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := convertToAllUnits(tt.value, tt.unit)
			if result.Nanoseconds != tt.wantNS {
				t.Errorf("convertToAllUnits() Nanoseconds = %v, want %v", result.Nanoseconds, tt.wantNS)
			}
			if result.Microseconds != tt.wantUS {
				t.Errorf("convertToAllUnits() Microseconds = %v, want %v", result.Microseconds, tt.wantUS)
			}
			if result.Milliseconds != tt.wantMS {
				t.Errorf("convertToAllUnits() Milliseconds = %v, want %v", result.Milliseconds, tt.wantMS)
			}
			if result.Seconds != tt.wantS {
				t.Errorf("convertToAllUnits() Seconds = %v, want %v", result.Seconds, tt.wantS)
			}
			if result.Minutes != tt.wantMin {
				t.Errorf("convertToAllUnits() Minutes = %v, want %v", result.Minutes, tt.wantMin)
			}
			if result.Hours != tt.wantH {
				t.Errorf("convertToAllUnits() Hours = %v, want %v", result.Hours, tt.wantH)
			}
		})
	}
}
