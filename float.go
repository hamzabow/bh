package main

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func (m model) updateFloatKeys(msg tea.KeyMsg) model {
	switch msg.String() {
	case "home", "ctrl+left":
		m.cursor = 0

	case "end", "ctrl+right":
		m.cursor = len(m.input)

	case "delete":
		if m.cursor < len(m.input) {
			m.input = m.input[:m.cursor] + m.input[m.cursor+1:]
			m = m.updateFloatConversions()
		}

	case "ctrl+w":
		if m.cursor > 0 {
			m.input = m.input[m.cursor:]
			m.cursor = 0
			m = m.updateFloatConversions()
		}

	case "w":
		if m.bitSize == 32 {
			m.bitSize = 64
		} else {
			m.bitSize = 32
		}
		m = m.updateFloatConversions()

	case "backspace":
		if len(m.input) > 0 && m.cursor > 0 {
			m.input = m.input[:m.cursor-1] + m.input[m.cursor:]
			m.cursor--
			m = m.updateFloatConversions()
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
		m = m.updateFloatConversions()

	default:
		char := msg.String()
		if len(char) == 1 && m.isValidFloatChar(char) {
			m.input = m.input[:m.cursor] + char + m.input[m.cursor:]
			m.cursor++
			m = m.updateFloatConversions()
		}
	}

	return m
}

func (m model) isValidFloatChar(char string) bool {
	lower := strings.ToLower(char)

	// Digits always valid
	if char >= "0" && char <= "9" {
		return true
	}

	// Decimal point: only one allowed, not after 'e'
	if char == "." {
		if strings.Contains(m.input, ".") {
			return false
		}
		// Don't allow dot after e
		if strings.ContainsAny(m.input, "eE") {
			return false
		}
		return true
	}

	// Minus sign: only at start or right after 'e'/'E'
	if char == "-" {
		if m.cursor == 0 && (len(m.input) == 0 || m.input[0] != '-') {
			return true
		}
		if m.cursor > 0 && (m.input[m.cursor-1] == 'e' || m.input[m.cursor-1] == 'E') {
			return true
		}
		return false
	}

	// Plus sign: only right after 'e'/'E'
	if char == "+" {
		if m.cursor > 0 && (m.input[m.cursor-1] == 'e' || m.input[m.cursor-1] == 'E') {
			return true
		}
		return false
	}

	// 'e'/'E' for scientific notation: only one, must have digits before it
	if lower == "e" {
		if strings.ContainsAny(m.input, "eE") {
			return false
		}
		// Must have at least one digit before e
		hasDigit := false
		for _, c := range m.input[:m.cursor] {
			if c >= '0' && c <= '9' {
				hasDigit = true
				break
			}
		}
		return hasDigit
	}

	// Letters for special values: allow n, a, i, f (for nan, inf)
	if lower == "n" || lower == "a" || lower == "i" || lower == "f" {
		candidate := m.input[:m.cursor] + char + m.input[m.cursor:]
		candidateLower := strings.ToLower(candidate)
		// Allow if building toward nan, inf, -inf, +inf, infinity, -infinity
		for _, special := range []string{"nan", "inf", "-inf", "+inf", "infinity", "-infinity", "+infinity"} {
			if strings.HasPrefix(special, candidateLower) || strings.HasPrefix(candidateLower, special) {
				return true
			}
		}
		return false
	}

	return false
}

func (m model) updateFloatConversions() model {
	if m.input == "" || m.input == "-" || m.input == "+" {
		m.floatVal, m.floatHex, m.floatBin = "", "", ""
		m.err = nil
		return m
	}

	inputLower := strings.ToLower(m.input)

	var val float64
	var err error

	switch inputLower {
	case "nan":
		val = math.NaN()
	case "inf", "+inf", "infinity", "+infinity":
		val = math.Inf(1)
	case "-inf", "-infinity":
		val = math.Inf(-1)
	default:
		val, err = strconv.ParseFloat(m.input, 64)
		if err != nil {
			m.err = err
			m.floatVal, m.floatHex, m.floatBin = "", "", ""
			return m
		}
	}

	m.err = nil

	if m.bitSize == 32 {
		val32 := float32(val)
		bits := math.Float32bits(val32)
		m.floatVal = formatFloatValue(float64(val32), bits == 0x80000000)
		m.floatHex = formatFloatHex32(bits)
		m.floatBin = m.formatFloatBinary32(bits)
	} else {
		bits := math.Float64bits(val)
		m.floatVal = formatFloatValue(val, bits == 0x8000000000000000)
		m.floatHex = formatFloatHex64(bits)
		m.floatBin = m.formatFloatBinary64(bits)
	}

	return m
}

func formatFloatValue(val float64, isNegZero bool) string {
	if math.IsNaN(val) {
		return "NaN"
	}
	if math.IsInf(val, 1) {
		return "+Infinity"
	}
	if math.IsInf(val, -1) {
		return "-Infinity"
	}
	if isNegZero {
		return "-0"
	}
	return strconv.FormatFloat(val, 'g', -1, 64)
}

func formatFloatHex32(bits uint32) string {
	hex := fmt.Sprintf("%08X", bits)
	return hex[:2] + separatorStyle.Render("_") + hex[2:4] + separatorStyle.Render("_") + hex[4:6] + separatorStyle.Render("_") + hex[6:8]
}

func formatFloatHex64(bits uint64) string {
	hex := fmt.Sprintf("%016X", bits)
	var result strings.Builder
	for i := 0; i < 16; i += 2 {
		if i > 0 {
			result.WriteString(separatorStyle.Render("_"))
		}
		result.WriteString(hex[i : i+2])
	}
	return result.String()
}

func (m model) formatFloatBinary32(bits uint32) string {
	sign := (bits >> 31) & 1
	exp := (bits >> 23) & 0xFF
	mantissa := bits & 0x7FFFFF

	signBin := fmt.Sprintf("%b", sign)
	expBin := fmt.Sprintf("%08b", exp)
	manBin := fmt.Sprintf("%023b", mantissa)

	return m.buildFloatBinaryDisplay(signBin, expBin, manBin, int(exp), 127)
}

func (m model) formatFloatBinary64(bits uint64) string {
	sign := (bits >> 63) & 1
	exp := (bits >> 52) & 0x7FF
	mantissa := bits & 0xFFFFFFFFFFFFF

	signBin := fmt.Sprintf("%b", sign)
	expBin := fmt.Sprintf("%011b", exp)
	manBin := fmt.Sprintf("%052b", mantissa)

	return m.buildFloatBinaryDisplay(signBin, expBin, manBin, int(exp), 1023)
}

// renderFieldBits renders a binary field with nibble brackets on top (4-bit),
// byte brackets on bottom (8-bit), and nibble-spaced bits in between.
// Partial leading groups get full-width brackets with leading positions empty.
// Returns 3 lines: top brackets, bits, bottom brackets.
func renderFieldBits(bits string, style lipgloss.Style) (top, bitLine, bot string) {
	nibbles := groupDigits(bits, 4)
	bytes := groupDigits(bits, 8)

	var topB, bitB, botB strings.Builder

	// Top + bit lines: nibble-based
	for gi, g := range nibbles {
		gLen := len(g.text)
		if gi > 0 {
			topB.WriteString(" ")
			bitB.WriteString(" ")
		}

		topB.WriteString(separatorStyle.Render("╭──╮"))

		if g.full {
			bitB.WriteString(style.Render(g.text))
		} else {
			pad := 4 - gLen
			bitB.WriteString(strings.Repeat(" ", pad))
			bitB.WriteString(style.Render(g.text))
		}
	}

	// Bottom line: byte-based
	for gi, g := range bytes {
		gLen := len(g.text)
		if gi > 0 {
			botB.WriteString(" ")
		}

		// Display width: digits + internal nibble spaces
		// A full byte = "xxxx xxxx" = 9 chars; partial depends on size
		nibbleSpaces := 0
		if gLen > 4 {
			nibbleSpaces = 1
		}
		// Partial leading group is padded to full nibble width on display
		displayDigits := gLen
		if gLen <= 4 && !g.full {
			displayDigits = 4
		} else if gLen > 4 && !g.full {
			// Partial byte: first nibble padded to 4, rest as-is
			firstNibble := gLen % 4
			if firstNibble == 0 {
				firstNibble = 4
			}
			displayDigits = 4 + (gLen - firstNibble) // pad first nibble to 4
			nibbleSpaces = 1
		}
		displayWidth := displayDigits + nibbleSpaces

		botB.WriteString(separatorStyle.Render("╰" + strings.Repeat("─", displayWidth-2) + "╯"))
	}

	return topB.String(), bitB.String(), botB.String()
}

// fieldDisplayWidth returns the visual width of a nibble-grouped binary field.
// Partial leading groups are expanded to full 4-char nibble width.
func fieldDisplayWidth(bits string) int {
	nibbles := groupDigits(bits, 4)
	w := 0
	for i, g := range nibbles {
		if i > 0 {
			w++ // nibble separator space
		}
		if g.full {
			w += len(g.text)
		} else {
			w += 4 // partial group padded to full nibble width
		}
	}
	return w
}

func (m model) buildFloatBinaryDisplay(signBin, expBin, manBin string, expVal, bias int) string {
	indent := "  "
	gap := "   "

	// Render each field
	_, signBits, _ := renderFieldBits(signBin, signStyle)
	expTop, expBits, expBot := renderFieldBits(expBin, exponentStyle)
	manTop, manBits, manBot := renderFieldBits(manBin, mantissaStyle)

	// Field display widths for label alignment
	signWidth := fieldDisplayWidth(signBin)
	expWidth := fieldDisplayWidth(expBin)

	// Labels centered above each field's bits
	manWidth := fieldDisplayWidth(manBin)

	// Center a label within a field width, returning (leftPad, rightPad)
	centerPad := func(labelLen, fieldWidth int) (int, int) {
		if fieldWidth >= labelLen {
			left := (fieldWidth - labelLen) / 2
			right := fieldWidth - labelLen - left
			return left, right
		}
		return 0, labelLen - fieldWidth
	}

	expLabelLeft, _ := centerPad(8, expWidth)
	manLabelLeft, _ := centerPad(8, manWidth)

	// Sign column is at least as wide as "Sign" (4 chars)
	signColWidth := signWidth
	if signColWidth < 4 {
		signColWidth = 4
	}

	// Build label line with proper centering
	var labelLine strings.Builder
	labelLine.WriteString(strings.Repeat(" ", (signColWidth-4)/2))
	labelLine.WriteString(signStyle.Render("Sign"))
	labelLine.WriteString(strings.Repeat(" ", signColWidth-4-(signColWidth-4)/2))
	labelLine.WriteString(gap)
	labelLine.WriteString(strings.Repeat(" ", expLabelLeft))
	labelLine.WriteString(exponentStyle.Render("Exponent"))
	// Fill remaining exp width + gap before mantissa label
	expRemainder := expWidth - 8 - expLabelLeft
	if expRemainder < 0 {
		expRemainder = 0
	}
	labelLine.WriteString(strings.Repeat(" ", expRemainder))
	labelLine.WriteString(gap)
	labelLine.WriteString(strings.Repeat(" ", manLabelLeft))
	labelLine.WriteString(mantissaStyle.Render("Mantissa"))

	// Sign field: no brackets, just the digit. Pad other lines to match width.
	signDisplayW := signWidth
	if signDisplayW < 4 {
		signDisplayW = 4 // at least as wide as "Sign" label
	}
	signTopPad := strings.Repeat(" ", signDisplayW)
	signBotPad := strings.Repeat(" ", signDisplayW)
	// Center the sign bit in the field width
	signBitPad := (signDisplayW - signWidth) / 2
	signBitsLine := strings.Repeat(" ", signBitPad) + signBits + strings.Repeat(" ", signDisplayW-signWidth-signBitPad)

	// Decoded exponent info
	var decoded string
	if expVal == 0 {
		decoded = exponentStyle.Render(fmt.Sprintf("denorm (2^%d)", 1-bias))
	} else if (bias == 127 && expVal == 255) || (bias == 1023 && expVal == 2047) {
		decoded = exponentStyle.Render("special (NaN/Inf)")
	} else {
		decoded = exponentStyle.Render(fmt.Sprintf("%d - %d = %d", expVal, bias, expVal-bias))
	}

	var s strings.Builder
	// Label line
	s.WriteString(indent + labelLine.String() + "\n")
	// Top brackets
	s.WriteString(indent + signTopPad + gap + expTop + gap + manTop + "\n")
	// Bit values
	s.WriteString(indent + signBitsLine + gap + expBits + gap + manBits + "\n")
	// Bottom brackets
	s.WriteString(indent + signBotPad + gap + expBot + gap + manBot + "\n")
	// Decoded info centered under exponent field
	decodedPad := signDisplayW + len(gap) - 2
	if m.bitSize == 64 {
		decodedPad += 2
	}
	if decodedPad < 0 {
		decodedPad = 0
	}
	s.WriteString(indent + strings.Repeat(" ", decodedPad) + decoded)

	return s.String()
}

func (m model) viewFloat() string {
	var s strings.Builder

	s.WriteString(titleStyle.Render("IEEE 754 Float Converter"))
	s.WriteString("\n\n")

	// Bit size tab bar (32/64 only)
	bitSizes := []string{"32-bit", "64-bit"}
	activeBit := fmt.Sprintf("%d-bit", m.bitSize)
	s.WriteString("  ")
	s.WriteString(keyHintStyle.Render("[W]") + " ")
	s.WriteString(renderTabBar(bitSizes, activeBit))
	s.WriteString("\n\n")

	// Input
	inputDisplay := "\n" + applyCursor(m.input, m.cursor, m.focused) + "\n"
	currentInputStyle := inputStyle.Width(80)
	if m.focused {
		currentInputStyle = focusedInputStyle.Width(80)
	}
	s.WriteString(renderStyledBorder(inputDisplay, "Decimal", currentInputStyle))
	s.WriteString("\n\n")

	if m.err != nil {
		s.WriteString(errorStyle.Render(fmt.Sprintf("Error: %s", m.err.Error())))
		s.WriteString("\n\n")
	}

	if m.floatVal != "" {
		s.WriteString(outputStyle.Render(labelStyle.Render("Value: ") + m.floatVal))
		s.WriteString("\n")
		s.WriteString(outputStyle.Render(labelStyle.Render("Hex:   ") + m.floatHex))
		s.WriteString("\n\n")
		s.WriteString(m.floatBin)
		s.WriteString("\n\n")
	}

	s.WriteString(helpStyle.Render("q: Quit · Ctrl+F: Integer mode · h: Help"))

	return s.String()
}
