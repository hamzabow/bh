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
	octal      string
	focused    bool
	bitSize    int
	signedMode bool
	overflow   bool
	showHelp   bool
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

		case "?":
			m.showHelp = !m.showHelp
			return m, nil
		}

		// When help is shown, ignore all other keys
		if m.showHelp {
			return m, nil
		}

		switch msg.String() {
		case "home", "ctrl+left":
			m.cursor = 0

		case "end", "ctrl+right":
			m.cursor = len(m.input)

		case "delete":
			if m.cursor < len(m.input) {
				wasPrefix := m.hasPrefix()
				m.input = m.input[:m.cursor] + m.input[m.cursor+1:]
				if wasPrefix && !m.hasPrefix() {
					m.inputType = "decimal"
				}
				m = m.updateConversions()
			}

		case "ctrl+w":
			if m.hasPrefix() {
				m.inputType = "decimal"
			}
			m.input = ""
			m.cursor = 0
			m = m.updateConversions()

		case "f1":
			switch m.inputType {
			case "decimal":
				m.inputType = "hex"
			case "hex":
				m.inputType = "octal"
			case "octal":
				m.inputType = "binary"
			case "binary":
				m.inputType = "decimal"
			}
			m.input = ""
			m.cursor = 0
			m.err = nil
			m.hex, m.binary, m.decimal, m.octal = "", "", "", ""
			m.overflow = false

		case "f2":
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

		case "f3":
			m.signedMode = !m.signedMode
			m = m.updateConversions()

		case "backspace":
			if len(m.input) > 0 && m.cursor > 0 {
				wasPrefix := m.hasPrefix()
				m.input = m.input[:m.cursor-1] + m.input[m.cursor:]
				m.cursor--
				// If we broke a prefix, revert to decimal and clear
				if wasPrefix && !m.hasPrefix() {
					m.inputType = "decimal"
					m = m.updateConversions()
				} else {
					m = m.updateConversions()
				}
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
			if len(char) == 1 {
				// Prefix detection: input is "0", cursor at 1, typing x/b/o
				if m.input == "0" && m.cursor == 1 {
					lower := strings.ToLower(char)
					switch lower {
					case "x":
						m.inputType = "hex"
						m.input = "0" + char
						m.cursor = 2
						m.err = nil
						m.overflow = false
						m.hex, m.binary, m.decimal, m.octal = "", "", "", ""
					case "b":
						m.inputType = "binary"
						m.input = "0" + char
						m.cursor = 2
						m.err = nil
						m.overflow = false
						m.hex, m.binary, m.decimal, m.octal = "", "", "", ""
					case "o":
						m.inputType = "octal"
						m.input = "0" + char
						m.cursor = 2
						m.err = nil
						m.overflow = false
						m.hex, m.binary, m.decimal, m.octal = "", "", "", ""
					default:
						if m.isValidChar(char) {
							m.input = m.input[:m.cursor] + char + m.input[m.cursor:]
							m.cursor++
							m = m.updateConversions()
						}
					}
				} else if m.isValidChar(char) {
					m.input = m.input[:m.cursor] + char + m.input[m.cursor:]
					m.cursor++
					m = m.updateConversions()
				}
			}
		}
	}

	return m, nil
}

func (m model) hasPrefix() bool {
	return len(m.input) >= 2 && m.input[0] == '0' &&
		(m.input[1] == 'x' || m.input[1] == 'X' ||
			m.input[1] == 'b' || m.input[1] == 'B' ||
			m.input[1] == 'o' || m.input[1] == 'O')
}

func (m model) isValidChar(char string) bool {
	// Don't allow typing before or inside a prefix
	if m.hasPrefix() && m.cursor < 2 {
		return false
	}

	switch m.inputType {
	case "decimal":
		if char == "-" && m.signedMode && m.cursor == 0 && !strings.Contains(m.input, "-") {
			return true
		}
		return char >= "0" && char <= "9"
	case "hex":
		return (char >= "0" && char <= "9") ||
			(strings.ToLower(char) >= "a" && strings.ToLower(char) <= "f")
	case "octal":
		return char >= "0" && char <= "7"
	case "binary":
		return char == "0" || char == "1"
	}
	return false
}

func (m model) updateConversions() model {
	if m.input == "" {
		m.hex, m.binary, m.decimal, m.octal = "", "", "", ""
		m.err = nil
		m.overflow = false
		return m
	}

	var num int64
	var err error

	parseInput := m.input
	if m.hasPrefix() {
		parseInput = m.input[2:]
	}

	if parseInput == "" || parseInput == "-" {
		m.hex, m.binary, m.decimal, m.octal = "", "", "", ""
		m.err = nil
		m.overflow = false
		return m
	}

	switch m.inputType {
	case "decimal":
		num, err = strconv.ParseInt(parseInput, 10, 64)
	case "hex":
		num, err = strconv.ParseInt(parseInput, 16, 64)
	case "octal":
		num, err = strconv.ParseInt(parseInput, 8, 64)
	case "binary":
		num, err = strconv.ParseInt(parseInput, 2, 64)
	}

	if err != nil {
		m.err = err
		m.hex, m.binary, m.decimal, m.octal = "", "", "", ""
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

	m.octal = fmt.Sprintf("%o", num)

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

		// Binary byte (8 bits) - with bounds checking
		endIdx := i + 8
		if endIdx > len(binary) {
			endIdx = len(binary)
		}
		byteBits := binary[i:endIdx]
		// Pad with zeros if needed to maintain 8-bit alignment
		for len(byteBits) < 8 {
			byteBits = "0" + byteBits
		}
		line2.WriteString(byteBits)

		// Vertical connectors for this byte
		line3.WriteString(separatorStyle.Render("│      │"))

		// Bottom with hex value for this byte
		hexByte := ""
		startIdx := byteCount * 2
		hexEndIdx := startIdx + 2
		if startIdx < len(hexRaw) {
			if hexEndIdx > len(hexRaw) {
				hexEndIdx = len(hexRaw)
			}
			hexByte = hexRaw[startIdx:hexEndIdx]
			// Pad with spaces if needed to maintain alignment
			for len(hexByte) < 2 {
				hexByte += " "
			}
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
				{"F1", "Cycle input base (Dec/Hex/Oct/Bin)"},
				{"F2", "Cycle bit size (8/16/32/64)"},
				{"F3", "Toggle signed/unsigned"},
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
				{"Ctrl+W", "Clear input"},
			},
		},
		{
			"General",
			[][2]string{
				{"?", "Toggle this help page"},
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

	s.WriteString(helpStyle.Render("Press ? to return"))

	return s.String()
}

func (m model) View() string {
	if m.showHelp {
		return m.viewHelp()
	}

	var s strings.Builder

	s.WriteString(titleStyle.Render("Number Base Converter"))
	s.WriteString("\n\n")

	// Input type tab bar
	inputTypes := []string{"Decimal", "Hex", "Octal", "Binary"}
	activeInput := map[string]string{"decimal": "Decimal", "hex": "Hex", "octal": "Octal", "binary": "Binary"}[m.inputType]
	s.WriteString("  ")
	s.WriteString(keyHintStyle.Render("[F1]")+" ")
	s.WriteString(renderTabBar(inputTypes, activeInput))
	s.WriteString("\n")

	// Bit size + signed/unsigned tab bars
	bitSizes := []string{"8-bit", "16-bit", "32-bit", "64-bit"}
	activeBit := fmt.Sprintf("%d-bit", m.bitSize)
	signedOpts := []string{"Unsigned", "Signed"}
	activeSigned := map[bool]string{true: "Signed", false: "Unsigned"}[m.signedMode]
	s.WriteString("  ")
	s.WriteString(keyHintStyle.Render("[F2]")+" ")
	s.WriteString(renderTabBar(bitSizes, activeBit))
	s.WriteString("\n")
	s.WriteString("  ")
	s.WriteString(keyHintStyle.Render("[F3]")+" ")
	s.WriteString(renderTabBar(signedOpts, activeSigned))
	s.WriteString("\n\n")

	// Range info
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
	s.WriteString("  ")
	s.WriteString(helpStyle.Render(rangeInfo))
	s.WriteString("\n\n")

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
		s.WriteString(outputStyle.Render(labelStyle.Render("Octal: ") + m.octal))
		s.WriteString("\n")
		s.WriteString(outputStyle.Render(labelStyle.Render("Binary: ") + m.binary))
		s.WriteString("\n\n")
	}

	s.WriteString(helpStyle.Render("q: Quit · ?: Help"))

	return s.String()
}

func main() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}