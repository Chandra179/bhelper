package collision

import (
	"math"
	"testing"
)

func TestIDGeneratorInterface(t *testing.T) {
	gen, _ := NewBase64Generator(8)
	var _ IDGenerator = gen
}

func TestGenerate(t *testing.T) {
	gen, err := NewBase64Generator(8)
	if err != nil {
		t.Fatalf("NewBase64Generator failed: %v", err)
	}
	id := gen.Generate()
	if len(id) != 8 {
		t.Errorf("Expected length 8, got %d", len(id))
	}
}

func TestTotalSpace(t *testing.T) {
	gen, err := NewBase64Generator(8)
	if err != nil {
		t.Fatalf("NewBase64Generator failed: %v", err)
	}
	space := gen.TotalSpace()
	expected := uint64(1)
	for i := 0; i < 8; i++ {
		expected *= 64
	}
	if space != expected {
		t.Errorf("Expected %d, got %d", expected, space)
	}
}

func TestNewBase64GeneratorZeroLength(t *testing.T) {
	_, err := NewBase64Generator(0)
	if err == nil {
		t.Error("Expected error for zero length, got nil")
	}
}

func TestNewBase64GeneratorNegativeLength(t *testing.T) {
	_, err := NewBase64Generator(-1)
	if err == nil {
		t.Error("Expected error for negative length, got nil")
	}
}

func TestTotalSpaceOverflow(t *testing.T) {
	gen, err := NewBase64Generator(11)
	if err != nil {
		t.Fatalf("NewBase64Generator failed: %v", err)
	}
	space := gen.TotalSpace()
	if space != math.MaxUint64 {
		t.Errorf("Expected MaxUint64 for overflow case, got %d", space)
	}
}

func TestBase62Generator(t *testing.T) {
	gen, err := NewBase62Generator(10)
	if err != nil {
		t.Fatalf("NewBase62Generator failed: %v", err)
	}
	id := gen.Generate()
	if len(id) != 10 {
		t.Errorf("Expected length 10, got %d", len(id))
	}

	space := gen.TotalSpace()
	expected := uint64(1)
	for i := 0; i < 10; i++ {
		expected *= 62
	}
	if space != expected {
		t.Errorf("Expected %d, got %d", expected, space)
	}

	if gen.Name() != "base62" {
		t.Errorf("Expected 'base62', got '%s'", gen.Name())
	}
}

func TestSnowflakeGenerator(t *testing.T) {
	gen, err := NewSnowflakeGenerator()
	if err != nil {
		t.Fatalf("NewSnowflakeGenerator failed: %v", err)
	}
	id1 := gen.Generate()
	id2 := gen.Generate()

	if id1 == id2 {
		t.Error("Expected different IDs")
	}

	space := gen.TotalSpace()
	if space == 0 {
		t.Error("Expected non-zero space")
	}

	if gen.Name() != "snowflake" {
		t.Errorf("Expected 'snowflake', got '%s'", gen.Name())
	}
}
