package gcmetrics

import (
	"fmt"
	"runtime/metrics"
	"strings"
	"sync"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	maxSamples     = 100
	sampleInterval = 500 * time.Millisecond
)

type state int

const (
	stateSelectProcess state = iota
	stateMonitoring
)

type Model struct {
	state          state
	processes      []ProcessInfo
	selectedIdx    int
	metrics        []Sample
	ticker         *time.Ticker
	quit           chan struct{}
	mu             sync.Mutex
	width, height  int
	visibleMetrics map[string]bool
}

type Sample struct {
	Timestamp time.Time
	PauseTime float64
	HeapAlloc uint64
	HeapFrees uint64
	HeapLive  uint64
	GCCycles  uint64
}

type ProcessInfo struct {
	PID  int
	Name string
}

type tickMsg time.Time

func NewModel(input string) Model {
	processes := discoverProcesses()
	selectedIdx := 0

	if len(input) > 0 && input != "list" {
		for i, p := range processes {
			if fmt.Sprintf("%d", p.PID) == input {
				selectedIdx = i
				break
			}
		}
	}

	visibleMetrics := map[string]bool{
		"PauseTime": true,
		"HeapAlloc": true,
		"HeapLive":  true,
	}

	return Model{
		state:          stateSelectProcess,
		processes:      processes,
		selectedIdx:    selectedIdx,
		metrics:        make([]Sample, 0, maxSamples),
		ticker:         time.NewTicker(sampleInterval),
		quit:           make(chan struct{}),
		visibleMetrics: visibleMetrics,
	}
}

func (m Model) Init() tea.Cmd {
	if m.state == stateMonitoring && len(m.processes) > 0 {
		return m.startMonitoringCmd()
	}
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return m.handleKey(msg)
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil
	case tickMsg:
		return m.handleTick()
	}

	return m, nil
}

func (m Model) View() string {
	switch m.state {
	case stateSelectProcess:
		return m.viewProcessList()
	case stateMonitoring:
		return m.viewMonitoring()
	default:
		return "Unknown state"
	}
}

func (m *Model) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch m.state {
	case stateSelectProcess:
		return m.handleProcessListKeys(msg)
	case stateMonitoring:
		return m.handleMonitoringKeys(msg)
	}
	return m, nil
}

func (m Model) handleProcessListKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q", "esc":
		if m.ticker != nil {
			m.ticker.Stop()
		}
		return m, tea.Quit
	case "enter", " ":
		if len(m.processes) > 0 {
			m.state = stateMonitoring
			return m, m.startMonitoringCmd()
		}
	case "up", "k":
		if m.selectedIdx > 0 {
			m.selectedIdx--
		}
	case "down", "j":
		if m.selectedIdx < len(m.processes)-1 {
			m.selectedIdx++
		}
	}
	return m, nil
}

func (m Model) handleMonitoringKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q", "esc":
		m.ticker.Stop()
		close(m.quit)
		return m, tea.Quit
	case "1":
		m.visibleMetrics["PauseTime"] = !m.visibleMetrics["PauseTime"]
	case "2":
		m.visibleMetrics["HeapAlloc"] = !m.visibleMetrics["HeapAlloc"]
	case "3":
		m.visibleMetrics["HeapFrees"] = !m.visibleMetrics["HeapFrees"]
	case "4":
		m.visibleMetrics["HeapLive"] = !m.visibleMetrics["HeapLive"]
	case "5":
		m.visibleMetrics["GCCycles"] = !m.visibleMetrics["GCCycles"]
	case "b":
		m.state = stateSelectProcess
		m.metrics = make([]Sample, 0, maxSamples)
	}
	return m, nil
}

func (m Model) startMonitoringCmd() tea.Cmd {
	return func() tea.Msg {
		ticker := time.NewTicker(sampleInterval)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				return tickMsg(time.Now())
			case <-m.quit:
				return tea.QuitMsg{}
			}
		}
	}
}

func (m *Model) handleTick() (tea.Model, tea.Cmd) {
	sample := m.collectSample()

	m.mu.Lock()
	m.metrics = append(m.metrics, sample)
	if len(m.metrics) > maxSamples {
		m.metrics = m.metrics[1:]
	}
	m.mu.Unlock()

	return m, tea.Tick(sampleInterval, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func (m Model) viewProcessList() string {
	var b strings.Builder

	title := titleStyle.Render("Select Go Process to Monitor")
	b.WriteString(title)
	b.WriteString("\n\n")

	if len(m.processes) == 0 {
		b.WriteString("No Go processes found.")
		return b.String()
	}

	for i, p := range m.processes {
		prefix := "  "
		if i == m.selectedIdx {
			prefix = selectedStyle.Render("▶ ")
		}
		line := fmt.Sprintf("%sPID %d: %s", prefix, p.PID, p.Name)
		b.WriteString(line)
		b.WriteString("\n")
	}

	b.WriteString("\n")
	b.WriteString("↑/j: up, ↓/k: down, Enter/Space: select, q/Esc: quit")

	return b.String()
}

func (m Model) viewMonitoring() string {
	var b strings.Builder

	title := titleStyle.Render("GC Metrics Monitor")
	b.WriteString(title)
	b.WriteString("\n\n")

	if len(m.metrics) == 0 {
		b.WriteString("Collecting metrics...")
		return b.String()
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	// Show latest metrics
	latest := m.metrics[len(m.metrics)-1]

	b.WriteString("Current Metrics:\n")
	b.WriteString(fmt.Sprintf("  GC Pause:    %.6f ms\n", latest.PauseTime*1000))
	b.WriteString(fmt.Sprintf("  Heap Alloc:  %d bytes\n", latest.HeapAlloc))
	b.WriteString(fmt.Sprintf("  Heap Frees:  %d bytes\n", latest.HeapFrees))
	b.WriteString(fmt.Sprintf("  Heap Live:   %d bytes\n", latest.HeapLive))
	b.WriteString(fmt.Sprintf("  GC Cycles:   %d\n", latest.GCCycles))

	b.WriteString("\n")

	// Show simple chart for Heap Live
	b.WriteString("Heap Live History (last " + fmt.Sprintf("%d", len(m.metrics)) + " samples):\n")
	chart := m.generateSimpleChart(m.metrics, func(s Sample) float64 { return float64(s.HeapLive) })
	b.WriteString(chart)

	b.WriteString("\n\n")
	b.WriteString("1-5: toggle metric, b: back to process list, q/Esc: quit")

	return b.String()
}

func (m Model) generateSimpleChart(samples []Sample, valueFunc func(Sample) float64) string {
	if len(samples) == 0 {
		return "No data"
	}

	// Find max value for scaling
	maxVal := 0.0
	for _, s := range samples {
		if v := valueFunc(s); v > maxVal {
			maxVal = v
		}
	}

	if maxVal == 0 {
		maxVal = 1
	}

	var b strings.Builder
	barWidth := 50

	for i, s := range samples {
		v := valueFunc(s)
		barLen := int((v / maxVal) * float64(barWidth))
		if barLen > barWidth {
			barLen = barWidth
		}

		bar := strings.Repeat("█", barLen)
		spaces := strings.Repeat(" ", barWidth-barLen)

		// Show timestamp (HH:MM:SS)
		timeStr := s.Timestamp.Format("15:04:05")
		b.WriteString(fmt.Sprintf("%s [%s] %s%.2f MB\n",
			timeStr, bar, spaces, v/(1024*1024)))

		// Only show last 10 samples to avoid too much output
		if i >= 9 {
			if len(samples) > 10 {
				b.WriteString(fmt.Sprintf("... and %d more samples\n", len(samples)-10))
			}
			break
		}
	}

	return b.String()
}

func (m *Model) collectSample() Sample {
	sample := Sample{
		Timestamp: time.Now(),
	}

	descriptions := make([]metrics.Description, 5)
	descriptions[0].Name = "/gc/pause:seconds"
	descriptions[1].Name = "/gc/heap/allocs:bytes"
	descriptions[2].Name = "/gc/heap/frees:bytes"
	descriptions[3].Name = "/gc/heap/go:bytes"
	descriptions[4].Name = "/gc/cycles:gc:seconds"

	samples := make([]metrics.Sample, 5)
	for i, d := range descriptions {
		samples[i].Name = d.Name
	}

	metrics.Read(samples)

	for i, s := range samples {
		switch descriptions[i].Name {
		case "/gc/pause:seconds":
			if descriptions[i].Kind == metrics.KindFloat64 {
				sample.PauseTime = s.Value.Float64()
			} else if descriptions[i].Kind == metrics.KindFloat64Histogram {
				sample.PauseTime = extractFromHistogram(s.Value.Float64Histogram())
			}
		case "/gc/heap/allocs:bytes":
			if descriptions[i].Kind == metrics.KindUint64 {
				sample.HeapAlloc = s.Value.Uint64()
			}
		case "/gc/heap/frees:bytes":
			if descriptions[i].Kind == metrics.KindUint64 {
				sample.HeapFrees = s.Value.Uint64()
			}
		case "/gc/heap/go:bytes":
			if descriptions[i].Kind == metrics.KindUint64 {
				sample.HeapLive = s.Value.Uint64()
			}
		case "/gc/cycles:gc:seconds":
			if descriptions[i].Kind == metrics.KindUint64 {
				sample.GCCycles = s.Value.Uint64()
			}
		}
	}

	return sample
}

func extractFromHistogram(h *metrics.Float64Histogram) float64 {
	if h == nil || len(h.Counts) == 0 {
		return 0
	}
	var total float64
	for _, count := range h.Counts {
		total += float64(count)
	}
	return total
}

var (
	titleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#7D56F4")).
			Bold(true)

	itemStyle     = lipgloss.NewStyle()
	selectedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#FA7979")).Bold(true)
)
