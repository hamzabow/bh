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
	input      string
	inputType  string
	cursor     int
	err        error
	hex        string
	binary     string
	decimal    string
	focused    bool
	bitSize    int
	signedMode bool
	overflow   bool
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

	bitSizeStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("93")).
		Bold(true).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("93")).
		Padding(0, 1)

	labelStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("243"))

	separatorStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("240"))
)

func initialModel() model {
	return model{
		input:      "",
		inputType:  "decimal",
		focused:    true,
		hex:        "",
		binary:     "",
		decimal:    "",
		bitSize:    32,
		signedMode: false,
		overflow:   false,
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
			m.overflow = false

		case "shift+tab":
			switch m.bitSize {
			case 8:
				m.bitSize = 16
			case 16:
				m.bitSize = 32
			case 32:
				m.bitSize = 64
			case 64:
				m.bitSize = 8
			}
			m = m.updateConversions()

		case "s":
			m.signedMode = !m.signedMode
			m = m.updateConversions()

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
		m.overflow = false
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
		m.overflow = false
		return m
	}

	m.err = nil
	m.overflow = false

	maxUnsigned, maxSigned, minSigned := m.getBitLimits()

	if m.signedMode {
		if num > maxSigned || num < minSigned {
			m.overflow = true
		}
		if num < 0 {
			num = m.toTwosComplement(num)
		}
	} else {
		if num < 0 {
			m.overflow = true
		} else if m.bitSize == 64 {
			// For 64-bit unsigned, any positive int64 value is valid
			// since we're using int64 internally but representing uint64
			m.overflow = false
		} else {
			if num > maxUnsigned {
				m.overflow = true
			}
		}
	}

	if !m.overflow {
		num = num & ((1 << m.bitSize) - 1)
	}

	displayNum := num
	if m.signedMode && !m.overflow {
		displayNum = m.fromTwosComplement(num)
	}

	m.decimal = fmt.Sprintf("%d", displayNum)

	hexWidth := m.bitSize / 4
	hexRaw := strings.ToUpper(fmt.Sprintf("%0*x", hexWidth, num))
	m.hex = m.formatHexWithBytes(hexRaw)

	binaryRaw := fmt.Sprintf("%0*b", m.bitSize, num)
	m.binary = m.formatBinaryWithBytes(binaryRaw)

	return m
}

func (m model) getBitLimits() (maxUnsigned, maxSigned, minSigned int64) {
	if m.bitSize == 64 {
		// For 64-bit, we need special handling since (1<<64) overflows
		maxUnsigned = -1 // Will be displayed as 18446744073709551615 when cast to uint64
		maxSigned = 9223372036854775807
		minSigned = -9223372036854775808
	} else {
		maxUnsigned = (1 << m.bitSize) - 1
		maxSigned = (1 << (m.bitSize - 1)) - 1
		minSigned = -(1 << (m.bitSize - 1))
	}
	return
}

func (m model) toTwosComplement(num int64) int64 {
	if num >= 0 {
		return num
	}
	return (1 << m.bitSize) + num
}

func (m model) fromTwosComplement(num int64) int64 {
	signBit := int64(1 << (m.bitSize - 1))
	if num&signBit != 0 {
		return num - (1 << m.bitSize)
	}
	return num
}

func (m model) formatBinaryWithBytes(binary string) string {
	if len(binary) <= 8 {
		return binary
	}

	indent := strings.Repeat(" ", 2) // 2-space indentation for all lines
	var line1, line2, line3, line4 strings.Builder

	// Build each line by processing bytes (8 bits each) with spacing between bytes
	hexRaw := strings.ToUpper(fmt.Sprintf("%0*x", m.bitSize/4, m.getNumFromBinary(binary)))
	byteCount := 0

	for i := 0; i < len(binary); i += 8 {
		// Add space between bytes (except for the first byte)
		if i > 0 {
			line1.WriteString(" ")
			line2.WriteString(" ")
			line3.WriteString(" ")
			line4.WriteString(" ")
		}

		// Top brackets for this byte (2 nibbles = 2 sets of ╭──╮)
		line1.WriteString(separatorStyle.Render("╭──╮╭──╮"))

		// Binary byte (8 bits)
		byteBits := binary[i : i+8]
		line2.WriteString(byteBits)

		// Vertical connectors for this byte
		line3.WriteString(separatorStyle.Render("│      │"))

		// Bottom with hex value for this byte
		hexByte := ""
		if byteCount < len(hexRaw)/2 {
			hexByte = hexRaw[byteCount*2 : byteCount*2+2]
		}
		line4.WriteString(separatorStyle.Render(fmt.Sprintf("╰──%s──╯", hexByte)))

		byteCount++
	}

	return "\n" + indent + line1.String() + "\n" + indent + line2.String() + "\n" + indent + line3.String() + "\n" + indent + line4.String()
}

func (m model) getNumFromBinary(binary string) int64 {
	if num, err := strconv.ParseInt(binary, 2, 64); err == nil {
		return num
	}
	return 0
}

func (m model) formatHexWithBytes(hex string) string {
	if len(hex) <= 2 {
		return hex
	}

	var result strings.Builder
	for i, char := range hex {
		if i > 0 && i%2 == 0 {
			result.WriteString(separatorStyle.Render("_"))
		}
		result.WriteRune(char)
	}

	return result.String()
}

func (m model) View() string {
	var s strings.Builder

	s.WriteString(titleStyle.Render("Number Base Converter"))
	s.WriteString("\n\n")

	bitInfo := fmt.Sprintf("%d-bit %s", m.bitSize, map[bool]string{true: "Signed", false: "Unsigned"}[m.signedMode])
	s.WriteString(bitSizeStyle.Render(bitInfo))
	s.WriteString("  ")

	maxUnsigned, maxSigned, minSigned := m.getBitLimits()
	var rangeInfo string
	if m.signedMode {
		rangeInfo = fmt.Sprintf("Range: %d to %d", minSigned, maxSigned)
	} else {
		if m.bitSize == 64 {
			rangeInfo = "Range: 0 to 18446744073709551615"
		} else {
			rangeInfo = fmt.Sprintf("Range: 0 to %d", maxUnsigned)
		}
	}
	s.WriteString(helpStyle.Render(rangeInfo))
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

	if m.overflow {
		s.WriteString(warningStyle.Render("⚠ Overflow: Value exceeds bit range"))
		s.WriteString("\n\n")
	}

	if m.decimal != "" {
		s.WriteString("Conversions:\n")
		s.WriteString(outputStyle.Render(labelStyle.Render("Decimal: ") + m.decimal))
		s.WriteString("\n")
		s.WriteString(outputStyle.Render(labelStyle.Render("Hexadecimal: ") + m.hex))
		s.WriteString("\n")
		s.WriteString(outputStyle.Render(labelStyle.Render("Binary: ") + m.binary))
		s.WriteString("\n\n")
	}

	s.WriteString(helpStyle.Render("Tab: Input type • Shift+Tab: Bit size • S: Signed/Unsigned • q/Ctrl+C: Quit"))

	return s.String()
}

func main() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}