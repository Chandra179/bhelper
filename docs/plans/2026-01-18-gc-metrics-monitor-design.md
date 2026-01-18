# GC Metrics Live Monitor Design

**Date**: 2026-01-18
**Feature Type**: Real-time monitoring and visualization

## Overview

A new feature for bhelper that provides real-time monitoring and visualization of Go garbage collector metrics. Users can attach to running Go processes and view live interactive charts showing GC behavior, heap allocation, pause times, and memory pressure.

## Architecture

The feature follows bhelper's plugin architecture, creating a new `feature/gcmetrics` package. The core component is a `GCMetricsMonitor` that implements the `Feature` interface using Go's `runtime/metrics` package (Go 1.21+) to capture real-time GC statistics.

Three main components:
- **Metrics Collector**: Captures GC data from runtime at regular intervals using `runtime/metrics.ReadAll()`
- **Process Manager**: Discovers and attaches to running Go processes (including bhelper itself)
- **Visualization Engine**: Renders interactive charts using bubble tea TUI components with live data updates

Integration with bhelper's TUI framework uses existing styles and patterns from `feature/collision`. Users select it from the feature list menu and switch between monitoring different processes.

## Components and Data Flow

**Metrics Collector** polls `runtime/metrics` every interval and captures:
- `/gc/pause:seconds` - GC pause durations
- `/gc/heap/allocs:bytes` - Heap allocations
- `/gc/heap/frees:bytes` - Heap deallocations
- `/gc/heap/go:bytes` - Live heap size
- `/gc/cycles:gc:seconds` - GC cycle times

Data flows: Collector → Channel → Buffer → Visualization. Each sample includes timestamp and all metric values. Buffered channels store recent samples for chart rendering.

**Process Manager** uses `github.com/karrick/gopsutil/process` to discover Go processes by reading `/proc/[pid]/cmdline` on Linux or platform APIs. Users select from a list or monitor the current process.

**Visualization Engine** uses `github.com/charmbracelet/bubbletea` with custom models:
- `ProcessList`: Selectable list of Go processes
- `MetricsChart`: Live-updating line charts using `bubbles/table` and lipgloss styling
- `StatsPanel`: Key metrics summary with color-coded thresholds (pause time warnings, heap growth alerts)

Charts auto-scale y-axis, show last N samples (scrollable), and can toggle different metrics. All updates happen through tea.Msg for smooth animations.

## Error Handling and Edge Cases

**Process Attachment Errors**: If a process terminates while monitoring, the feature shows an alert "Process [PID] ended" and returns to process selection. Invalid processes are filtered from the list automatically.

**Permission Issues**: Processes requiring elevated permissions are detected. Users can skip those processes or restart bhelper with elevated permissions. Failed access attempts are logged without crashing.

**Metrics Unavailable**: Go versions before 1.21 don't support `runtime/metrics`. The feature checks Go version on startup and displays "Requires Go 1.21+" if unavailable. Handles cases where specific metrics aren't exposed in older versions.

**Memory Limits**: Metrics buffer is capped at 1000 samples per metric. Old samples are discarded FIFO. Channel sizes are bounded (capacity: 100) to prevent unbounded growth.

**Concurrent Access**: All state updates happen through tea.Msg, ensuring thread-safe access to metrics data. The collector runs as a single goroutine per monitored process, properly cancelled when switching processes or exiting.

**Signal Handling**: Gracefully handles SIGINT/SIGTERM, ensuring goroutines terminate cleanly and temporary files are cleaned up.

## Testing Strategy

**Unit Tests** for the metrics collector:
- Mock `runtime/metrics` using test doubles to verify correct sampling intervals
- Test buffer overflow scenarios (1000+ samples) ensure FIFO behavior
- Validate metric parsing and timestamp generation
- Edge cases: empty samples, malformed data, concurrent access

**Integration Tests** for process manager:
- Spawn test Go processes and verify discovery works
- Test process termination during monitoring scenarios
- Validate permission error handling
- Cross-platform testing (Linux, macOS) for process discovery

**End-to-End Tests** for visualization:
- Use tea's `TestProgram` to simulate TUI interactions
- Verify process selection, metric toggling, and mode switching
- Test chart rendering with sample data streams
- Validate cleanup and termination behavior

**Performance Tests**:
- Benchmark collector with different sampling rates (100ms, 500ms, 1s)
- Measure memory footprint with 1K, 10K samples
- Verify no goroutine leaks after extended monitoring sessions

**Test Data**: Use synthetic metrics generators to simulate various GC patterns (low churn, high allocation, memory leaks) for reproducible testing.

## Key Dependencies

- `runtime/metrics` (Go 1.21+) - GC metrics collection
- `github.com/karrick/gopsutil/process` - Process discovery
- `github.com/charmbracelet/bubbletea` - TUI framework
- `github.com/charmbracelet/lipgloss` - Styling
- `github.com/charmbracelet/bubbles` - UI components

## Success Criteria

- Can monitor bhelper itself and other running Go processes
- Live charts update smoothly without lag (target: <50ms refresh)
- Charts show all key GC metrics with accurate timestamps
- Gracefully handles process termination and permission errors
- Memory usage stays bounded (<50MB) during extended monitoring
- All tests pass including cross-platform tests
