package main

import (
	"bhelper/feature"
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

// CLIMode represents the current state of the CLI
type CLIMode int

const (
	ModeFeatureList CLIMode = iota
	ModeFeatureHelp
	ModeFeatureExecute
)

// CLI is the main TUI model
type CLI struct {
	registry        *feature.FeatureRegistry
	mode            CLIMode
	selectedIndex   int
	selectedFeature feature.Feature
	textInput       textinput.Model
	output          string
	history         *History
}

// NewCLI creates a new CLI instance
func NewCLI(registry *feature.FeatureRegistry) CLI {
	ti := textinput.New()
	ti.Placeholder = "Type your input..."
	ti.Width = 70

	return CLI{
		registry:      registry,
		mode:          ModeFeatureList,
		selectedIndex: 0,
		textInput:     ti,
		history:       NewHistory(50),
	}
}

func (c CLI) Init() tea.Cmd {
	return textinput.Blink
}

func (c CLI) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch c.mode {
		case ModeFeatureList:
			return c.updateFeatureList(msg)
		case ModeFeatureHelp:
			return c.updateFeatureHelp(msg)
		case ModeFeatureExecute:
			return c.updateFeatureExecute(msg)
		}
	}
	return c, nil
}

// updateFeatureList handles feature list navigation
func (c CLI) updateFeatureList(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	features := c.registry.List()

	switch msg.String() {
	case "ctrl+c", "q":
		return c, tea.Quit

	case "up", "k":
		if c.selectedIndex > 0 {
			c.selectedIndex--
		}

	case "down", "j":
		if c.selectedIndex < len(features)-1 {
			c.selectedIndex++
		}

	case "enter":
		c.selectedFeature = features[c.selectedIndex]
		c.mode = ModeFeatureExecute
		c.textInput.Focus()
		return c, textinput.Blink

	case "h", "?":
		c.selectedFeature = features[c.selectedIndex]
		c.mode = ModeFeatureHelp
	}

	return c, nil
}

// updateFeatureHelp handles help screen
func (c CLI) updateFeatureHelp(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "q":
		return c, tea.Quit

	case "esc", "backspace":
		c.mode = ModeFeatureList

	case "enter":
		c.mode = ModeFeatureExecute
		c.textInput.Focus()
		return c, textinput.Blink
	}

	return c, nil
}

// updateFeatureExecute handles feature execution
func (c CLI) updateFeatureExecute(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c":
		return c, tea.Quit

	case "esc":
		c.mode = ModeFeatureList
		c.textInput.Blur()
		c.textInput.SetValue("")
		c.output = ""
		return c, nil

	case "ctrl+h":
		c.mode = ModeFeatureHelp
		c.textInput.Blur()
		return c, nil

	case "enter":
		// Execute the feature
		input := c.textInput.Value()
		if input != "" {
			result, err := c.selectedFeature.Execute(input)
			if err != nil {
				c.output = fmt.Sprintf("Error: %v", err)
			} else {
				c.output = result
			}
		}
		return c, nil

	case "ctrl+z":
		if state := c.history.Undo(); state != nil {
			c.textInput.SetValue(*state)
			c.output = ""
		}
		return c, nil

	case "ctrl+y":
		if state := c.history.Redo(); state != nil {
			c.textInput.SetValue(*state)
			c.output = ""
		}
		return c, nil

	default:
		oldValue := c.textInput.Value()
		var cmd tea.Cmd
		c.textInput, cmd = c.textInput.Update(msg)

		if oldValue != c.textInput.Value() {
			c.history.Push(oldValue)
			c.output = ""
		}
		return c, cmd
	}
}

func (c CLI) View() string {
	switch c.mode {
	case ModeFeatureList:
		return c.renderFeatureList()
	case ModeFeatureHelp:
		return c.renderFeatureHelp()
	case ModeFeatureExecute:
		return c.renderFeatureExecute()
	}
	return ""
}

// renderFeatureList shows all available features
func (c CLI) renderFeatureList() string {
	var s strings.Builder

	title := titleStyle.Render("Choose: ")
	s.WriteString(title + "\n")

	features := c.registry.List()
	for i, f := range features {
		cursor := "  "
		style := featureStyle

		if i == c.selectedIndex {
			cursor = "→ "
			style = selectedFeatureStyle
		}

		line := fmt.Sprintf("%s%s - %s", cursor, f.Name(), f.Description())
		s.WriteString(style.Render(line) + "\n")
	}

	s.WriteString("\n" + helpStyle.Render("↑/↓: navigate • ENTER: select • H/?: help • Q: quit"))

	return s.String()
}

// renderFeatureHelp shows detailed help for selected feature
func (c CLI) renderFeatureHelp() string {
	var s strings.Builder

	title := titleStyle.Render(fmt.Sprintf("📖 Help: %s", c.selectedFeature.Name()))
	s.WriteString(title + "\n\n")

	s.WriteString(sectionStyle.Render("Description:") + "\n")
	s.WriteString(c.selectedFeature.Help() + "\n\n")

	examples := c.selectedFeature.Examples()
	if len(examples) > 0 {
		s.WriteString(sectionStyle.Render("Examples:") + "\n")
		for _, ex := range examples {
			s.WriteString(exampleStyle.Render(fmt.Sprintf("Input: %s", ex.Input)) + "\n")
			s.WriteString(fmt.Sprintf("  → %s\n\n", ex.Description))
		}
	}

	s.WriteString(helpStyle.Render("ENTER: use feature • ESC: back • Q: quit"))

	return s.String()
}

// renderFeatureExecute shows the feature execution interface
func (c CLI) renderFeatureExecute() string {
	var s strings.Builder

	title := titleStyle.Render(fmt.Sprintf("⚡ %s", c.selectedFeature.Name()))
	s.WriteString(title + "\n\n")

	// Input
	s.WriteString(labelStyle.Render("Input: ") + c.textInput.View() + "\n\n")

	// Output
	if c.output != "" {
		s.WriteString(sectionStyle.Render("Result:") + "\n")
		s.WriteString(outputBoxStyle.Render(c.output) + "\n\n")
	}

	s.WriteString(helpStyle.Render("ENTER: execute • CTRL+H: help • CTRL+Z: undo • CTRL+Y: redo • ESC: back"))

	return s.String()
}
