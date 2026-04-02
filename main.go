package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type groupingMode int

const (
	groupOff      groupingMode = iota
	groupBrackets
	groupSpaces
	groupBoth
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
	groupMode  groupingMode
	permView   bool
}

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
			if m.cursor > 0 {
				// Delete backward to prefix boundary or start of input
				stop := 0
				if m.hasPrefix() {
					stop = 2
				}
				if m.cursor > stop {
					m.input = m.input[:stop] + m.input[m.cursor:]
					m.cursor = stop
				}
				// If only prefix remains, clear everything
				if m.hasPrefix() && len(m.input) <= 2 {
					m.inputType = "decimal"
					m.input = ""
					m.cursor = 0
				}
				m = m.updateConversions()
			}

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

		case "f4":
			m.groupMode = (m.groupMode + 1) % 4

		case "p":
			m.permView = !m.permView
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
	if m.permView {
		m.binary = m.formatBinaryPerms(binaryRaw)
	} else {
		m.binary = m.formatBinaryWithBytes(binaryRaw)
	}

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

func (m model) formatBinaryPerms(binary string) string {
	indent := strings.Repeat(" ", 2)
	groups := groupDigits(binary, 3)

	var topLine, digitLine, botLine, rwxLine strings.Builder

	for gi, g := range groups {
		if gi > 0 {
			topLine.WriteString(" ")
			digitLine.WriteString(" ")
			botLine.WriteString(" ")
			rwxLine.WriteString(" ")
		}

		if g.full {
			topLine.WriteString(separatorStyle.Render("╭─╮"))
			octalVal := binaryGroupToOctal(g.text)
			botLine.WriteString(separatorStyle.Render("╰") + permAnnotStyle.Render(fmt.Sprintf("%d", octalVal)) + separatorStyle.Render("╯"))
			rwxLine.WriteString(octalRWX(octalVal))
		} else {
			topLine.WriteString(strings.Repeat(" ", len(g.text)))
			botLine.WriteString(strings.Repeat(" ", len(g.text)))
			rwxLine.WriteString(strings.Repeat(" ", len(g.text)))
		}
		digitLine.WriteString(g.text)
	}

	return "\n" + indent + topLine.String() + "\n" + indent + digitLine.String() + "\n" + indent + botLine.String() + "\n" + indent + rwxLine.String()
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
				{"F4", "Cycle grouping (Off/Brackets/Spaces/Both)"},
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
	s.WriteString("\n")

	groupedOpts := []string{"Off", "Brackets", "Spaces", "Both"}
	activeGrouped := map[groupingMode]string{
		groupOff: "Off", groupBrackets: "Brackets", groupSpaces: "Spaces", groupBoth: "Both",
	}[m.groupMode]
	s.WriteString("  ")
	s.WriteString(keyHintStyle.Render("[F4]") + " Grouped: ")
	s.WriteString(renderTabBar(groupedOpts, activeGrouped))
	s.WriteString("\n")

	permOpts := []string{"Off", "On"}
	activePerm := "Off"
	if m.permView {
		activePerm = "On"
	}
	s.WriteString("  ")
	s.WriteString(keyHintStyle.Render("[P]") + "  Unix Permissions: ")
	s.WriteString(renderTabBar(permOpts, activePerm))
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

	inputDisplay := m.renderGroupedInputDisplay()

	currentInputStyle := inputStyle.Width(80)
	if m.focused {
		currentInputStyle = focusedInputStyle.Width(80)
	}
	s.WriteString(currentInputStyle.Render(inputDisplay))
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

func (m model) renderGroupedInputDisplay() string {
	digits := m.input
	prefixLen := 0
	if m.hasPrefix() {
		prefixLen = 2
		digits = m.input[2:]
	} else if m.inputType == "decimal" && len(m.input) > 0 && m.input[0] == '-' {
		prefixLen = 1
		digits = m.input[1:]
	}

	prefix := m.input[:prefixLen]
	cursorInDigits := m.cursor - prefixLen

	if len(digits) == 0 {
		cursor := applyCursor(m.input, m.cursor, m.focused)
		return "\n" + cursor + "\n"
	}

	permActive := m.permView && m.groupMode != groupOff

	switch m.inputType {
	case "decimal":
		return m.renderBracketDecOct(prefix, digits, cursorInDigits)
	case "octal":
		if permActive {
			return m.renderBracketOctalPerms(prefix, digits, cursorInDigits)
		}
		return m.renderBracketDecOct(prefix, digits, cursorInDigits)
	case "binary":
		if permActive {
			return m.renderBracketBinaryPerms(prefix, digits, cursorInDigits)
		}
		return m.renderBracketBinary(prefix, digits, cursorInDigits)
	case "hex":
		return m.renderBracketHex(prefix, digits, cursorInDigits)
	}
	return applyCursor(m.input, m.cursor, m.focused)
}

func (m model) renderBracketDecOct(prefix, digits string, cursorInDigits int) string {
	groups := groupDigits(digits, 3)
	n := len(digits)
	showBrackets := m.groupMode == groupBrackets || m.groupMode == groupBoth
	showSpaces := m.groupMode == groupSpaces || m.groupMode == groupBoth

	var digitLine, botLine strings.Builder
	if prefix != "" {
		digitLine.WriteString(prefix)
		botLine.WriteString(strings.Repeat(" ", len(prefix)))
	}

	displayPos := len(prefix)
	posMap := make([]int, n+1)
	rawIdx := 0

	for gi, g := range groups {
		gLen := len(g.text)

		// Add space between groups
		if showSpaces && gi > 0 {
			digitLine.WriteByte(' ')
			botLine.WriteString(" ")
			displayPos++
		}

		if showBrackets && g.full {
			botLine.WriteString(separatorStyle.Render("╰─╯"))
		} else {
			botLine.WriteString(strings.Repeat(" ", gLen))
		}

		for i := 0; i < gLen; i++ {
			posMap[rawIdx] = displayPos
			digitLine.WriteByte(g.text[i])
			rawIdx++
			displayPos++
		}
	}
	posMap[n] = displayPos

	cursorDisplayPos := 0
	if cursorInDigits < 0 {
		cursorDisplayPos = m.cursor
	} else if cursorInDigits >= n {
		cursorDisplayPos = posMap[n]
	} else {
		cursorDisplayPos = posMap[cursorInDigits]
	}

	renderedDigitLine := applyCursor(digitLine.String(), cursorDisplayPos, m.focused)
	if showBrackets {
		return "\n" + renderedDigitLine + "\n" + botLine.String()
	}
	return "\n" + renderedDigitLine + "\n" + strings.Repeat(" ", len(prefix))
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

func (m model) renderBracketBinary(prefix, digits string, cursorInDigits int) string {
	n := len(digits)
	showBrackets := m.groupMode == groupBrackets || m.groupMode == groupBoth
	showSpaces := m.groupMode == groupSpaces || m.groupMode == groupBoth

	nibbleGroups := groupDigits(digits, 4)
	byteGroups := groupDigits(digits, 8)

	pad := strings.Repeat(" ", len(prefix))

	// Top line: nibble brackets (╭──╮ per full nibble)
	var topLine strings.Builder
	topLine.WriteString(pad)

	// Digit line
	var digitLine strings.Builder
	digitLine.WriteString(prefix)
	displayPos := len(prefix)
	posMap := make([]int, n+1)

	// Bottom line
	var botLine strings.Builder
	botLine.WriteString(pad)


	// Build digit line and top line together (nibble-based iteration)
	nibbleIdx := 0 // tracks position in digits for nibble spacing
	for gi, g := range nibbleGroups {
		gLen := len(g.text)

		// Add space between nibble groups
		if showSpaces && gi > 0 {
			digitLine.WriteByte(' ')
			topLine.WriteString(" ")
			displayPos++
		}

		if showBrackets && g.full {
			topLine.WriteString(separatorStyle.Render("╭──╮"))
		} else {
			topLine.WriteString(strings.Repeat(" ", gLen))
		}

		for i := 0; i < gLen; i++ {
			posMap[nibbleIdx] = displayPos
			digitLine.WriteByte(digits[nibbleIdx])
			nibbleIdx++
			displayPos++
		}
	}
	posMap[n] = displayPos

	// Build bottom line (byte-based)
	if showBrackets {
		// We need to figure out how wide each byte group is in display chars
		// In Both mode, a full byte = 8 digits + 1 internal nibble space = 9 display chars
		// In Brackets mode, a full byte = 8 digits
		for gi, g := range byteGroups {
			gLen := len(g.text)
			if showSpaces && gi > 0 {
				botLine.WriteString(" ")
			}
			if g.full {
				hexVal := fmt.Sprintf("%02X", m.getNumFromBinary(g.text))
				if showSpaces {
					botLine.WriteString(separatorStyle.Render(fmt.Sprintf("╰───%s──╯", hexVal)))
				} else {
					botLine.WriteString(separatorStyle.Render(fmt.Sprintf("╰──%s──╯", hexVal)))
				}
			} else {
				displayWidth := gLen
				if showSpaces && gLen > 4 {
					displayWidth = gLen + (gLen-1)/4
				}
				botLine.WriteString(strings.Repeat(" ", displayWidth))
			}
		}
	}

	cursorDisplayPos := 0
	if cursorInDigits < 0 {
		cursorDisplayPos = m.cursor
	} else if cursorInDigits >= n {
		cursorDisplayPos = posMap[n]
	} else {
		cursorDisplayPos = posMap[cursorInDigits]
	}

	renderedDigitLine := applyCursor(digitLine.String(), cursorDisplayPos, m.focused)
	if showBrackets {
		return topLine.String() + "\n" + renderedDigitLine + "\n" + botLine.String()
	}
	return "\n" + renderedDigitLine + "\n"
}

func (m model) renderBracketBinaryPerms(prefix, digits string, cursorInDigits int) string {
	n := len(digits)
	showBrackets := m.groupMode == groupBrackets || m.groupMode == groupBoth
	showSpaces := m.groupMode == groupSpaces || m.groupMode == groupBoth

	groups := groupDigits(digits, 3)
	pad := strings.Repeat(" ", len(prefix))

	var topLine, digitLine, botLine, rwxLine strings.Builder
	topLine.WriteString(pad)
	digitLine.WriteString(prefix)
	botLine.WriteString(pad)
	rwxLine.WriteString(pad)

	displayPos := len(prefix)
	posMap := make([]int, n+1)
	rawIdx := 0

	for gi, g := range groups {
		gLen := len(g.text)

		if showSpaces && gi > 0 {
			digitLine.WriteByte(' ')
			topLine.WriteString(" ")
			botLine.WriteString(" ")
			rwxLine.WriteString(" ")
			displayPos++
		}

		if showBrackets && g.full {
			topLine.WriteString(separatorStyle.Render("╭─╮"))
			octalVal := binaryGroupToOctal(g.text)
			botLine.WriteString(separatorStyle.Render("╰") + permAnnotStyle.Render(fmt.Sprintf("%d", octalVal)) + separatorStyle.Render("╯"))
			rwxLine.WriteString(octalRWX(octalVal))
		} else {
			topLine.WriteString(strings.Repeat(" ", gLen))
			botLine.WriteString(strings.Repeat(" ", gLen))
			rwxLine.WriteString(strings.Repeat(" ", gLen))
		}

		for i := 0; i < gLen; i++ {
			posMap[rawIdx] = displayPos
			digitLine.WriteByte(digits[rawIdx])
			rawIdx++
			displayPos++
		}
	}
	posMap[n] = displayPos

	cursorDisplayPos := 0
	if cursorInDigits < 0 {
		cursorDisplayPos = m.cursor
	} else if cursorInDigits >= n {
		cursorDisplayPos = posMap[n]
	} else {
		cursorDisplayPos = posMap[cursorInDigits]
	}

	renderedDigitLine := applyCursor(digitLine.String(), cursorDisplayPos, m.focused)
	if showBrackets {
		return topLine.String() + "\n" + renderedDigitLine + "\n" + botLine.String() + "\n" + rwxLine.String()
	}
	return "\n" + renderedDigitLine + "\n"
}

func (m model) renderBracketOctalPerms(prefix, digits string, cursorInDigits int) string {
	n := len(digits)
	showBrackets := m.groupMode == groupBrackets || m.groupMode == groupBoth
	showSpaces := m.groupMode == groupSpaces || m.groupMode == groupBoth

	var digitLine, binLine, rwxLine strings.Builder
	if prefix != "" {
		digitLine.WriteString(prefix)
		binLine.WriteString(strings.Repeat(" ", len(prefix)))
		rwxLine.WriteString(strings.Repeat(" ", len(prefix)))
	}

	displayPos := len(prefix)
	posMap := make([]int, n+1)

	for i := 0; i < n; i++ {
		if showSpaces && i > 0 {
			digitLine.WriteByte(' ')
			binLine.WriteString(" ")
			rwxLine.WriteString(" ")
			displayPos++
		}

		// Center the digit: space + digit + space
		digitLine.WriteByte(' ')
		displayPos++
		posMap[i] = displayPos
		digitLine.WriteByte(digits[i])
		displayPos++
		digitLine.WriteByte(' ')
		displayPos++

		if showBrackets {
			octalVal := int(digits[i] - '0')
			binLine.WriteString(permAnnotStyle.Render(fmt.Sprintf("%03b", octalVal)))
			rwxLine.WriteString(octalRWX(octalVal))
		} else {
			binLine.WriteString("   ")
			rwxLine.WriteString("   ")
		}
	}
	posMap[n] = displayPos

	cursorDisplayPos := 0
	if cursorInDigits < 0 {
		cursorDisplayPos = m.cursor
	} else if cursorInDigits >= n {
		cursorDisplayPos = posMap[n]
	} else {
		cursorDisplayPos = posMap[cursorInDigits]
	}

	renderedDigitLine := applyCursor(digitLine.String(), cursorDisplayPos, m.focused)
	if showBrackets {
		return "\n" + renderedDigitLine + "\n" + binLine.String() + "\n" + rwxLine.String()
	}
	return "\n" + renderedDigitLine + "\n"
}

func (m model) renderBracketHex(prefix, digits string, cursorInDigits int) string {
	groups := groupDigits(digits, 2)
	n := len(digits)
	showBrackets := m.groupMode == groupBrackets || m.groupMode == groupBoth
	showSpaces := m.groupMode == groupSpaces || m.groupMode == groupBoth

	var digitLine, botLine strings.Builder
	if prefix != "" {
		digitLine.WriteString(prefix)
		botLine.WriteString(strings.Repeat(" ", len(prefix)))
	}

	displayPos := len(prefix)
	posMap := make([]int, n+1)
	rawIdx := 0

	for gi, g := range groups {
		gLen := len(g.text)

		// Add space between groups
		if showSpaces && gi > 0 {
			digitLine.WriteByte(' ')
			botLine.WriteString(" ")
			displayPos++
		}

		if showBrackets && g.full {
			botLine.WriteString(separatorStyle.Render("╰╯"))
		} else {
			botLine.WriteString(strings.Repeat(" ", gLen))
		}

		for i := 0; i < gLen; i++ {
			posMap[rawIdx] = displayPos
			digitLine.WriteByte(g.text[i])
			rawIdx++
			displayPos++
		}
	}
	posMap[n] = displayPos

	cursorDisplayPos := 0
	if cursorInDigits < 0 {
		cursorDisplayPos = m.cursor
	} else if cursorInDigits >= n {
		cursorDisplayPos = posMap[n]
	} else {
		cursorDisplayPos = posMap[cursorInDigits]
	}

	renderedDigitLine := applyCursor(digitLine.String(), cursorDisplayPos, m.focused)
	if showBrackets {
		return "\n" + renderedDigitLine + "\n" + botLine.String()
	}
	return "\n" + renderedDigitLine + "\n" + strings.Repeat(" ", len(prefix))
}

func main() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}