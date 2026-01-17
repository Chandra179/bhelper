package collision

import "testing"

func TestNewCollisionAnalyzer(t *testing.T) {
	analyzer := NewCollisionAnalyzer()

	if analyzer == nil {
		t.Fatal("Expected non-nil analyzer")
	}

	if analyzer.ID() != "collision" {
		t.Errorf("Expected ID 'collision', got '%s'", analyzer.ID())
	}
}

func TestCollisionAnalyzerExecute(t *testing.T) {
	analyzer := NewCollisionAnalyzer()

	result, err := analyzer.Execute("base64:10:1000/sec")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if result == "" {
		t.Error("Expected non-empty result")
	}
}

func TestCollisionAnalyzerExecuteError(t *testing.T) {
	analyzer := NewCollisionAnalyzer()

	_, err := analyzer.Execute("invalid")
	if err == nil {
		t.Error("Expected error for invalid input")
	}
}

func TestCollisionAnalyzerBase62(t *testing.T) {
	analyzer := NewCollisionAnalyzer()

	result, err := analyzer.Execute("base62:8:500/min")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if result == "" {
		t.Error("Expected non-empty result")
	}
}

func TestCollisionAnalyzerSnowflake(t *testing.T) {
	analyzer := NewCollisionAnalyzer()

	result, err := analyzer.Execute("snowflake:0:10000/sec")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if result == "" {
		t.Error("Expected non-empty result")
	}
}

func TestCollisionAnalyzerNanosecond(t *testing.T) {
	analyzer := NewCollisionAnalyzer()

	result, err := analyzer.Execute("base64:10:1000/ns")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if result == "" {
		t.Error("Expected non-empty result")
	}
}
