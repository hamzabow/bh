package main

import "github.com/charmbracelet/lipgloss"

var (
	titleStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("205")).
		Bold(true).
		Padding(1, 2)

	inputStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("255")).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("62")).
		Padding(0, 2).
		Width(50)

	focusedInputStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("255")).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("205")).
		Padding(0, 2).
		Width(50)

	outputStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("39")).
		Border(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("59")).
		Padding(1, 2).
		Width(80)

	errorStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("196")).
		Bold(true)

	helpStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("241")).
		Italic(true)

	warningStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("208")).
		Bold(true)

	activeTabStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("205")).
		Bold(true)

	inactiveTabStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("241"))

	tabSepStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("241"))

	keyHintStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("244"))

	labelStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("243"))

	separatorStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("240"))

	permAnnotStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("109"))

	signStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("196")).
		Bold(true)

	exponentStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("214"))

	mantissaStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("39"))
)
