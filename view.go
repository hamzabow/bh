package main

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func (m model) View() string {
	if m.showHelp {
		return m.viewHelp()
	}
	if m.floatMode {
		return m.viewFloat()
	}
	return m.viewInteger()
}

func (m model) viewHelp() string {
	var s strings.Builder

	s.WriteString(titleStyle.Render("Keyboard Shortcuts"))
	s.WriteString("\n\n")

	sections := []struct {
		heading string
		keys    [][2]string
	}{
		{
			"Modes",
			[][2]string{
				{"F", "Toggle float mode (IEEE 754)"},
				{"T", "Cycle input base (Dec/Hex/Oct/Bin)"},
				{"W", "Cycle bit size (8/16/32/64)"},
				{"S", "Toggle signed/unsigned"},
				{"G", "Cycle grouping (Off/Brackets/Spaces/Both)"},
				{"P", "Toggle permissions view (octal/rwx)"},
			},
		},
		{
			"Navigation",
			[][2]string{
				{"←/→", "Move cursor left/right"},
				{"Home, Ctrl+←", "Move cursor to beginning"},
				{"End, Ctrl+→", "Move cursor to end"},
			},
		},
		{
			"Editing",
			[][2]string{
				{"Backspace", "Delete character before cursor"},
				{"Delete", "Delete character at cursor"},
				{"Ctrl+W", "Delete word backward"},
			},
		},
		{
			"General",
			[][2]string{
				{"H", "Toggle this help page"},
				{"q, Ctrl+C", "Quit"},
			},
		},
	}

	keyStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("205")).Bold(true).Width(16)
	descStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("252"))
	headingStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("39")).Bold(true)

	for _, section := range sections {
		s.WriteString("  ")
		s.WriteString(headingStyle.Render(section.heading))
		s.WriteString("\n")
		for _, kv := range section.keys {
			s.WriteString("    ")
			s.WriteString(keyStyle.Render(kv[0]))
			s.WriteString(descStyle.Render(kv[1]))
			s.WriteString("\n")
		}
		s.WriteString("\n")
	}

	s.WriteString(helpStyle.Render("Press h to return"))

	return s.String()
}
