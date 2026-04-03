package main

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func renderStyledBorder(content, label string, style lipgloss.Style) string {
	// Get border color from the style
	borderFg := style.GetBorderTopForeground()
	borderStyle := lipgloss.NewStyle().Foreground(borderFg)
	labelFg := style.GetForeground()
	labelRendered := lipgloss.NewStyle().Foreground(labelFg).Bold(true).Render(label)

	// Render without top border
	noTop := style.BorderTop(false)
	body := noTop.Render(content)

	// Build top border: ╭─Label─────...──╮
	width := lipgloss.Width(body)
	b := lipgloss.RoundedBorder()

	labelWidth := lipgloss.Width(labelRendered)
	// Top line: corner + dash + label + remaining dashes + corner
	fillWidth := width - 2 - 1 - labelWidth // 2 corners, 1 dash before label
	if fillWidth < 0 {
		fillWidth = 0
	}
	topLine := borderStyle.Render(b.TopLeft+"─") + labelRendered + borderStyle.Render(strings.Repeat("─", fillWidth)+b.TopRight)

	return topLine + "\n" + body
}

func inputTypeLabel(inputType string) string {
	switch inputType {
	case "decimal":
		return "Decimal"
	case "hex":
		return "Hexadecimal"
	case "binary":
		return "Binary"
	case "octal":
		return "Octal"
	default:
		return inputType
	}
}

func renderTabBar(options []string, active string) string {
	var parts []string
	for _, opt := range options {
		if opt == active {
			parts = append(parts, activeTabStyle.Render(opt))
		} else {
			parts = append(parts, inactiveTabStyle.Render(opt))
		}
	}
	return strings.Join(parts, tabSepStyle.Render(" │ "))
}

func applyCursor(text string, pos int, focused bool) string {
	if !focused {
		return text
	}
	inputColor := lipgloss.NewStyle().Foreground(lipgloss.Color("255"))
	cs := lipgloss.NewStyle().Reverse(true).Foreground(lipgloss.Color("255"))
	if pos < len(text) {
		return text[:pos] + cs.Render(string(text[pos])) + inputColor.Render(text[pos+1:])
	}
	return text + cs.Render(" ")
}

type digitGroup struct {
	text string
	full bool
}

func groupDigits(digits string, groupSize int) []digitGroup {
	n := len(digits)
	if n == 0 {
		return nil
	}

	firstGroupSize := n % groupSize
	if firstGroupSize == 0 {
		firstGroupSize = groupSize
	}

	var groups []digitGroup
	first := digits[:firstGroupSize]
	groups = append(groups, digitGroup{first, len(first) == groupSize})

	for i := firstGroupSize; i < n; i += groupSize {
		end := i + groupSize
		if end > n {
			end = n
		}
		g := digits[i:end]
		groups = append(groups, digitGroup{g, len(g) == groupSize})
	}
	return groups
}

func binaryGroupToOctal(bits string) int {
	val := 0
	for _, b := range bits {
		val = val*2 + int(b-'0')
	}
	return val
}

func octalRWX(val int) string {
	var s strings.Builder
	for _, ch := range []struct {
		bit  int
		char string
	}{{4, "r"}, {2, "w"}, {1, "x"}} {
		if val&ch.bit != 0 {
			s.WriteString(permAnnotStyle.Render(ch.char))
		} else {
			s.WriteString(separatorStyle.Render("-"))
		}
	}
	return s.String()
}
