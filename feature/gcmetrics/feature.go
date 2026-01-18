package gcmetrics

import (
	"fmt"
	"runtime"
	"strings"

	"bhelper/feature"

	tea "github.com/charmbracelet/bubbletea"
)

// MonitorFeature provides GC metrics visualization
type MonitorFeature struct{}

func NewMonitorFeature() *MonitorFeature {
	return &MonitorFeature{}
}

func (m *MonitorFeature) ID() string {
	return "gcmetrics"
}

func (m *MonitorFeature) Name() string {
	return "GC Metrics Monitor"
}

func (m *MonitorFeature) Description() string {
	return "Live GC metrics visualization with interactive charts"
}

func (m *MonitorFeature) Help() string {
	var b strings.Builder
	b.WriteString("GC Metrics Monitor provides real-time visualization of Go garbage collector behavior:\n\n")
	b.WriteString("• Live GC pause times and cycle frequency\n")
	b.WriteString("• Heap allocation and deallocation rates\n")
	b.WriteString("• Live heap size tracking\n")
	b.WriteString("• Interactive charts with multiple metrics\n")
	b.WriteString("• Monitor any running Go process\n\n")
	b.WriteString("Requirements: Go 1.21+ (uses runtime/metrics)\n\n")
	b.WriteString("Use this feature to:\n")
	b.WriteString("  - Identify GC pressure in applications\n")
	b.WriteString("  - Compare GC behavior across processes\n")
	b.WriteString("  - Detect memory leaks through allocation patterns\n")
	b.WriteString("  - Optimize code for better GC performance")
	return b.String()
}

func (m *MonitorFeature) Examples() []feature.Example {
	return []feature.Example{
		{Input: "", Description: "Start interactive GC monitor for current process"},
		{Input: "1234", Description: "Monitor specific process ID (if accessible)"},
		{Input: "list", Description: "Show available Go processes"},
	}
}

func (m *MonitorFeature) Execute(input string) (string, error) {
	goVersion := runtime.Version()
	major, minor, _, _ := parseGoVersion(goVersion)
	if major < 1 || (major == 1 && minor < 21) {
		return "", fmt.Errorf("GC Metrics Monitor requires Go 1.21+, current version: %s", goVersion)
	}

	model := NewModel(input)
	p := tea.NewProgram(model, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		return "", fmt.Errorf("error running UI: %w", err)
	}

	return "GC monitoring session ended", nil
}

func parseGoVersion(version string) (major, minor, patch int, err error) {
	_, err = fmt.Sscanf(version, "go%d.%d.%d", &major, &minor, &patch)
	if err != nil {
		return 0, 0, 0, err
	}
	return major, minor, patch, nil
}
