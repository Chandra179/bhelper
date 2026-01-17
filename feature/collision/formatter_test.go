package collision

import (
	"math/big"
	"strings"
	"testing"
	"time"
)

func TestFormatResult(t *testing.T) {
	mathResult := &MathResult{
		TotalSpace:         1000,
		TotalIDs:           100,
		Probability:        big.NewFloat(0.005),
		ExpectedCollisions: 1,
		TimeToCollision: &TimeResult{
			P50:  time.Hour,
			P01:  time.Minute,
			P001: time.Second,
		},
	}

	simResult := &SimResult{
		Collisions:  2,
		Iterations:  1000,
		Probability: 0.002,
	}

	output := FormatResult("base64", 8, 1000, mathResult, simResult)

	if output == "" {
		t.Error("Expected non-empty output")
	}

	if !strings.Contains(output, "Collision Analysis") {
		t.Error("Expected 'Collision Analysis' in output")
	}

	if !strings.Contains(output, "Mathematical") {
		t.Error("Expected 'Mathematical' in output")
	}

	if !strings.Contains(output, "Simulation") {
		t.Error("Expected 'Simulation' in output")
	}
}

func TestFormatNumber(t *testing.T) {
	tests := []struct {
		input    uint64
		expected string
	}{
		{100, "100"},
		{1000, "1,000"},
		{1000000, "1,000,000"},
		{1000000000, "1,000,000,000"},
	}

	for _, test := range tests {
		result := formatNumber(test.input)
		if result != test.expected {
			t.Errorf("Expected '%s' for %d, got '%s'", test.expected, test.input, result)
		}
	}
}

func TestFormatDuration(t *testing.T) {
	tests := []struct {
		input    time.Duration
		expected string
	}{
		{time.Second, "1.0 seconds"},
		{time.Minute, "1.0 minutes"},
		{time.Hour, "1.0 hours"},
		{24 * time.Hour, "1.0 days"},
	}

	for _, test := range tests {
		result := formatDuration(test.input)
		if result != test.expected {
			t.Errorf("Expected '%s' for %v, got '%s'", test.expected, test.input, result)
		}
	}
}
