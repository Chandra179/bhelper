package collision

import (
	"testing"
)

func TestParseInput(t *testing.T) {
	input := "base64:10:1000/sec"

	config, err := ParseInput(input)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if config.Format != "base64" {
		t.Errorf("Expected 'base64', got '%s'", config.Format)
	}
	if config.Length != 10 {
		t.Errorf("Expected 10, got %d", config.Length)
	}
	if config.Rate != 1000 {
		t.Errorf("Expected 1000, got %d", config.Rate)
	}
	if config.RateUnit != "sec" {
		t.Errorf("Expected 'sec', got '%s'", config.RateUnit)
	}
}

func TestParseInputWithMinute(t *testing.T) {
	input := "base62:8:500/min"

	config, err := ParseInput(input)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if config.Format != "base62" {
		t.Errorf("Expected 'base62', got '%s'", config.Format)
	}
	if config.Length != 8 {
		t.Errorf("Expected 8, got %d", config.Length)
	}
	if config.Rate != 500 {
		t.Errorf("Expected 500, got %d", config.Rate)
	}
	if config.RateUnit != "min" {
		t.Errorf("Expected 'min', got '%s'", config.RateUnit)
	}
}

func TestParseInputWithMillis(t *testing.T) {
	input := "snowflake:0:10000/ms"

	config, err := ParseInput(input)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if config.Format != "snowflake" {
		t.Errorf("Expected 'snowflake', got '%s'", config.Format)
	}
	if config.Length != 0 {
		t.Errorf("Expected 0, got %d", config.Length)
	}
	if config.Rate != 10000 {
		t.Errorf("Expected 10000, got %d", config.Rate)
	}
	if config.RateUnit != "ms" {
		t.Errorf("Expected 'ms', got '%s'", config.RateUnit)
	}
}

func TestParseInputWithNanos(t *testing.T) {
	input := "base64:8:1000/ns"

	config, err := ParseInput(input)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if config.Format != "base64" {
		t.Errorf("Expected 'base64', got '%s'", config.Format)
	}
	if config.Length != 8 {
		t.Errorf("Expected 8, got %d", config.Length)
	}
	if config.Rate != 1000 {
		t.Errorf("Expected 1000, got %d", config.Rate)
	}
	if config.RateUnit != "ns" {
		t.Errorf("Expected 'ns', got '%s'", config.RateUnit)
	}
}

func TestParseInputError(t *testing.T) {
	input := "invalid"

	_, err := ParseInput(input)
	if err == nil {
		t.Error("Expected error for invalid input")
	}
}

func TestParseInputInvalidLength(t *testing.T) {
	input := "base64:abc:1000/sec"

	_, err := ParseInput(input)
	if err == nil {
		t.Error("Expected error for invalid length")
	}
}
