package gcmetrics

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
)

func discoverProcesses() []ProcessInfo {
	processes := []ProcessInfo{}

	// Add current process
	processes = append(processes, ProcessInfo{
		PID:  os.Getpid(),
		Name: fmt.Sprintf("bhelper (current)"),
	})

	// Try to discover other Go processes on Linux
	if runtime.GOOS == "linux" {
		goProcesses := discoverLinuxProcesses()
		for _, p := range goProcesses {
			if p.PID != os.Getpid() {
				processes = append(processes, p)
			}
		}
	}

	return processes
}

func discoverLinuxProcesses() []ProcessInfo {
	processes := []ProcessInfo{}

	procDir := "/proc"
	entries, err := os.ReadDir(procDir)
	if err != nil {
		return processes
	}

	goPattern := regexp.MustCompile(`^/tmp/go-build\d+|/usr/local/go/`)

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		pid, err := strconv.Atoi(entry.Name())
		if err != nil {
			continue
		}

		cmdlinePath := filepath.Join(procDir, entry.Name(), "cmdline")
		cmdlineBytes, err := os.ReadFile(cmdlinePath)
		if err != nil {
			continue
		}

		cmdline := string(cmdlineBytes)
		if len(cmdline) == 0 {
			continue
		}

		// Check if it's a Go process by looking for Go indicators
		if isGoProcess(cmdline, goPattern) {
			name := extractProcessName(cmdline)
			processes = append(processes, ProcessInfo{
				PID:  pid,
				Name: name,
			})
		}
	}

	return processes
}

func isGoProcess(cmdline string, goPattern *regexp.Regexp) bool {
	return goPattern.MatchString(cmdline) ||
		contains(cmdline, "go-build") ||
		contains(cmdline, "/go") ||
		contains(cmdline, "/usr/bin/go")
}

func extractProcessName(cmdline string) string {
	// Get the first part before null byte
	if idx := findByte(cmdline, 0); idx >= 0 {
		cmdline = cmdline[:idx]
	}

	// Extract just the binary name
	if len(cmdline) > 50 {
		cmdline = "..." + cmdline[len(cmdline)-47:]
	}

	return cmdline
}

func findByte(s string, b byte) int {
	for i := range s {
		if s[i] == b {
			return i
		}
	}
	return -1
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) &&
		(s == substr ||
			(len(s) > len(substr) &&
				(s[:len(substr)] == substr ||
					s[len(s)-len(substr):] == substr ||
					findSubstring(s, substr) >= 0)))
}

func findSubstring(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}
