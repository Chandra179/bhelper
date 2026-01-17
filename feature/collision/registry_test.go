package collision

import "testing"

func TestNewRegistry(t *testing.T) {
	reg := NewGeneratorRegistry()
	if reg == nil {
		t.Fatal("Expected non-nil registry")
	}
}

func TestRegistryRegister(t *testing.T) {
	reg := NewGeneratorRegistry()
	gen, err := NewBase64Generator(8)
	if err != nil {
		t.Fatalf("NewBase64Generator failed: %v", err)
	}

	reg.Register(gen)

	got, ok := reg.Get("base64")
	if !ok {
		t.Error("Expected generator to be registered")
	}
	if got.Name() != "base64" {
		t.Errorf("Expected 'base64', got '%s'", got.Name())
	}
}

func TestRegistryGetNotFound(t *testing.T) {
	reg := NewGeneratorRegistry()

	_, ok := reg.Get("unknown")
	if ok {
		t.Error("Expected false for unknown generator")
	}
}

func TestRegistryList(t *testing.T) {
	reg := NewGeneratorRegistry()
	gen1, _ := NewBase64Generator(8)
	gen2, _ := NewBase62Generator(10)
	reg.Register(gen1)
	reg.Register(gen2)

	list := reg.List()
	if len(list) != 2 {
		t.Errorf("Expected 2 generators, got %d", len(list))
	}
}

func TestRegistryNames(t *testing.T) {
	reg := NewGeneratorRegistry()
	gen1, _ := NewBase64Generator(8)
	reg.Register(gen1)

	names := reg.Names()
	if len(names) != 1 {
		t.Errorf("Expected 1 name, got %d", len(names))
	}
	if names[0] != "base64" {
		t.Errorf("Expected 'base64', got '%s'", names[0])
	}
}
