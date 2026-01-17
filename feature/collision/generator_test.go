package collision

import "testing"

func TestIDGeneratorInterface(t *testing.T) {
	var _ IDGenerator = NewBase64Generator(8)
}

func TestGenerate(t *testing.T) {
	gen := NewBase64Generator(8)
	id := gen.Generate()
	if len(id) != 8 {
		t.Errorf("Expected length 8, got %d", len(id))
	}
}

func TestTotalSpace(t *testing.T) {
	gen := NewBase64Generator(8)
	space := gen.TotalSpace()
	expected := uint64(1)
	for i := 0; i < 8; i++ {
		expected *= 64
	}
	if space != expected {
		t.Errorf("Expected %d, got %d", expected, space)
	}
}
