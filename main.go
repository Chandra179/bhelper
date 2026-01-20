package main

import (
	"bhelper/feature"
	"bhelper/feature/collision"
	"bhelper/feature/time"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	// Register all features
	registry := feature.NewFeatureRegistry()
	registry.Register(feature.NewCharacterAnalyzer())
	registry.Register(feature.NewTimezoneAnalyzer())
	registry.Register(time.NewTimeConverter())
	registry.Register(collision.NewCollisionAnalyzer())
	// registry.Register(NewWeatherForecast())
	// ... register 100 features here

	// Start CLI with all registered features
	p := tea.NewProgram(NewCLI(registry))
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
