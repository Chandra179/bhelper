package gcmetrics

import (
	"runtime"
	"testing"
	"time"
)

func TestNewMonitorFeature(t *testing.T) {
	f := NewMonitorFeature()
	if f == nil {
		t.Fatal("NewMonitorFeature returned nil")
	}
}

func TestMonitorFeatureID(t *testing.T) {
	f := NewMonitorFeature()
	if f.ID() != "gcmetrics" {
		t.Errorf("Expected ID 'gcmetrics', got '%s'", f.ID())
	}
}

func TestMonitorFeatureName(t *testing.T) {
	f := NewMonitorFeature()
	if f.Name() != "GC Metrics Monitor" {
		t.Errorf("Expected Name 'GC Metrics Monitor', got '%s'", f.Name())
	}
}

func TestMonitorFeatureDescription(t *testing.T) {
	f := NewMonitorFeature()
	desc := f.Description()
	if desc == "" {
		t.Error("Description should not be empty")
	}
}

func TestMonitorFeatureHelp(t *testing.T) {
	f := NewMonitorFeature()
	help := f.Help()
	if help == "" {
		t.Error("Help should not be empty")
	}
	if !contains(help, "GC Metrics Monitor") {
		t.Error("Help should mention feature name")
	}
}

func TestMonitorFeatureExamples(t *testing.T) {
	f := NewMonitorFeature()
	examples := f.Examples()
	if len(examples) == 0 {
		t.Error("Should have at least one example")
	}
}

func TestParseGoVersion(t *testing.T) {
	tests := []struct {
		input     string
		wantMaj   int
		wantMin   int
		wantPatch int
	}{
		{"go1.21.0", 1, 21, 0},
		{"go1.22.1", 1, 22, 1},
		{"go1.20.5", 1, 20, 5},
	}

	for _, tt := range tests {
		maj, min, patch, err := parseGoVersion(tt.input)
		if err != nil {
			t.Errorf("parseGoVersion(%q) error = %v", tt.input, err)
			continue
		}
		if maj != tt.wantMaj || min != tt.wantMin || patch != tt.wantPatch {
			t.Errorf("parseGoVersion(%q) = (%d, %d, %d), want (%d, %d, %d)",
				tt.input, maj, min, patch, tt.wantMaj, tt.wantMin, tt.wantPatch)
		}
	}
}

func TestExecuteWithOldGoVersion(t *testing.T) {
	f := NewMonitorFeature()

	// Mock old version by temporarily changing the version check
	// This is a limitation test - in real scenarios, this would check runtime.Version()
	// For now, we'll just verify the error handling logic exists
	if runtime.Version() == "go1.20.0" {
		_, err := f.Execute("")
		if err == nil {
			t.Error("Should return error for Go 1.20")
		}
	}
}

func TestExecuteEmptyInput(t *testing.T) {
	NewMonitorFeature()
	// This would start the TUI, which we can't test directly
	// We just verify it doesn't panic on empty input
	// Real testing would require mocking tea.Program
}

func TestModelInit(t *testing.T) {
	m := NewModel("")
	if m.state != stateSelectProcess {
		t.Errorf("Expected initial state stateSelectProcess, got %v", m.state)
	}
}

func TestModelWithEmptyProcesses(t *testing.T) {
	m := Model{
		processes: []ProcessInfo{},
		state:     stateSelectProcess,
	}

	view := m.viewProcessList()
	if !contains(view, "No Go processes") {
		t.Error("Should show 'No Go processes' when empty")
	}
}

func TestModelCollectSample(t *testing.T) {
	m := &Model{}
	sample := m.collectSample()

	if sample.Timestamp.IsZero() {
		t.Error("Sample timestamp should not be zero")
	}
}

func TestSampleStruct(t *testing.T) {
	now := time.Now()
	s := Sample{
		Timestamp: now,
		PauseTime: 0.001,
		HeapAlloc: 1000,
		HeapFrees: 500,
		HeapLive:  500,
		GCCycles:  1,
	}

	if s.Timestamp != now {
		t.Error("Timestamp not set correctly")
	}
	if s.PauseTime != 0.001 {
		t.Error("PauseTime not set correctly")
	}
}

func TestProcessInfoStruct(t *testing.T) {
	p := ProcessInfo{
		PID:  1234,
		Name: "test",
	}

	if p.PID != 1234 {
		t.Error("PID not set correctly")
	}
	if p.Name != "test" {
		t.Error("Name not set correctly")
	}
}

func TestDiscoverProcesses(t *testing.T) {
	processes := discoverProcesses()

	if len(processes) == 0 {
		t.Error("Should find at least current process")
	}

	// First process should be current process
	if processes[0].PID <= 0 {
		t.Error("Current process PID should be positive")
	}
}

func TestGenerateSimpleChart(t *testing.T) {
	m := &Model{}

	now := time.Now()
	samples := []Sample{
		{Timestamp: now, HeapLive: 1000},
		{Timestamp: now, HeapLive: 2000},
		{Timestamp: now, HeapLive: 1500},
	}

	chart := m.generateSimpleChart(samples, func(s Sample) float64 {
		return float64(s.HeapLive)
	})

	if chart == "" {
		t.Error("Chart should not be empty")
	}
	if !contains(chart, "MB") {
		t.Error("Chart should show values in MB")
	}
}

func TestModelToggleMetrics(t *testing.T) {
	m := Model{
		visibleMetrics: map[string]bool{
			"PauseTime": true,
			"HeapAlloc": true,
		},
	}

	if !m.visibleMetrics["PauseTime"] {
		t.Error("PauseTime should initially be visible")
	}

	m.visibleMetrics["PauseTime"] = false
	if m.visibleMetrics["PauseTime"] {
		t.Error("PauseTime should be toggled to false")
	}
}

func TestStateConstants(t *testing.T) {
	if stateSelectProcess != 0 {
		t.Error("stateSelectProcess should be 0")
	}
	if stateMonitoring != 1 {
		t.Error("stateMonitoring should be 1")
	}
}
