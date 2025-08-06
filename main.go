package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	input       string
	inputType   string
	cursor      int
	err         error
	hex         string
	binary      string
	decimal     string
	focused     bool
}

var (
	titleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("205")).
			Bold(true).
			Padding(1, 2)

	inputStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("86")).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("62")).
			Padding(1, 2).
			Width(50)

	focusedInputStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("86")).
				Border(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color("205")).
				Padding(1, 2).
				Width(50)

	outputStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("39")).
			Border(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("59")).
			Padding(1, 2).
			Width(50)

	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("196")).
			Bold(true)

	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241")).
			Italic(true)
)

func initialModel() model {
	return model{
		input:     "",
		inputType: "decimal",
		focused:   true,
		hex:       "",
		binary:    "",
		decimal:   "",
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit

		case "tab":
			switch m.inputType {
			case "decimal":
				m.inputType = "hex"
			case "hex":
				m.inputType = "binary"
			case "binary":
				m.inputType = "decimal"
			}
			m.input = ""
			m.cursor = 0
			m.err = nil
			m.hex, m.binary, m.decimal = "", "", ""

		case "backspace":
			if len(m.input) > 0 && m.cursor > 0 {
				m.input = m.input[:m.cursor-1] + m.input[m.cursor:]
				m.cursor--
				m = m.updateConversions()
			}

		case "left":
			if m.cursor > 0 {
				m.cursor--
			}

		case "right":
			if m.cursor < len(m.input) {
				m.cursor++
			}

		case "enter":
			m = m.updateConversions()

		default:
			char := msg.String()
			if len(char) == 1 && m.isValidChar(char) {
				m.input = m.input[:m.cursor] + char + m.input[m.cursor:]
				m.cursor++
				m = m.updateConversions()
			}
		}
	}

	return m, nil
}

func (m model) isValidChar(char string) bool {
	switch m.inputType {
	case "decimal":
		return char >= "0" && char <= "9"
	case "hex":
		return (char >= "0" && char <= "9") || 
			   (strings.ToLower(char) >= "a" && strings.ToLower(char) <= "f")
	case "binary":
		return char == "0" || char == "1"
	}
	return false
}

func (m model) updateConversions() model {
	if m.input == "" {
		m.hex, m.binary, m.decimal = "", "", ""
		m.err = nil
		return m
	}

	var num int64
	var err error

	switch m.inputType {
	case "decimal":
		num, err = strconv.ParseInt(m.input, 10, 64)
	case "hex":
		num, err = strconv.ParseInt(m.input, 16, 64)
	case "binary":
		num, err = strconv.ParseInt(m.input, 2, 64)
	}

	if err != nil {
		m.err = err
		m.hex, m.binary, m.decimal = "", "", ""
		return m
	}

	m.err = nil
	m.decimal = fmt.Sprintf("%d", num)
	m.hex = strings.ToUpper(fmt.Sprintf("%x", num))
	m.binary = fmt.Sprintf("%b", num)

	return m
}

func (m model) View() string {
	var s strings.Builder

	s.WriteString(titleStyle.Render("Number Base Converter"))
	s.WriteString("\n\n")

	inputLabel := fmt.Sprintf("Input (%s):", strings.Title(m.inputType))
	s.WriteString(inputLabel)
	s.WriteString("\n")

	inputDisplay := m.input
	if m.focused {
		if m.cursor < len(inputDisplay) {
			inputDisplay = inputDisplay[:m.cursor] + "│" + inputDisplay[m.cursor:]
		} else {
			inputDisplay += "│"
		}
	}

	if m.focused {
		s.WriteString(focusedInputStyle.Render(inputDisplay))
	} else {
		s.WriteString(inputStyle.Render(inputDisplay))
	}
	s.WriteString("\n\n")

	if m.err != nil {
		s.WriteString(errorStyle.Render(fmt.Sprintf("Error: %s", m.err.Error())))
		s.WriteString("\n\n")
	}

	if m.decimal != "" {
		s.WriteString("Conversions:\n")
		s.WriteString(outputStyle.Render(fmt.Sprintf("Decimal: %s", m.decimal)))
		s.WriteString("\n")
		s.WriteString(outputStyle.Render(fmt.Sprintf("Hexadecimal: %s", m.hex)))
		s.WriteString("\n")
		s.WriteString(outputStyle.Render(fmt.Sprintf("Binary: %s", m.binary)))
		s.WriteString("\n\n")
	}

	s.WriteString(helpStyle.Render("Tab: Switch input type • Enter: Convert • q/Ctrl+C: Quit"))

	return s.String()
}

func main() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}