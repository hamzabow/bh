package main

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
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

func (m model) buildFloatBinaryDisplay(signBin, expBin, manBin string, expVal, bias int) string {
	indent := "  "

	// Header line with field labels
	signLabel := signStyle.Render("Sign")
	expLabel := exponentStyle.Render("Exponent")
	manLabel := mantissaStyle.Render("Mantissa")

	// Top brackets
	expTopWidth := len(expBin)
	manTopWidth := len(manBin)

	signTop := signStyle.Render("╭╮")
	expTop := exponentStyle.Render("╭" + strings.Repeat("─", expTopWidth-2) + "╮")
	manTop := mantissaStyle.Render("╭" + strings.Repeat("─", manTopWidth-2) + "╮")

	// Bit values
	signBits := signStyle.Render(signBin)
	expBits := exponentStyle.Render(expBin)
	manBits := mantissaStyle.Render(manBin)

	// Bottom brackets
	signBot := signStyle.Render("╰╯")
	expBot := exponentStyle.Render("╰" + strings.Repeat("─", expTopWidth-2) + "╯")
	manBot := mantissaStyle.Render("╰" + strings.Repeat("─", manTopWidth-2) + "╯")

	// Decoded info
	var decoded string
	if expVal == 0 {
		decoded = exponentStyle.Render(fmt.Sprintf("denorm (2^%d)", 1-bias))
	} else if (bias == 127 && expVal == 255) || (bias == 1023 && expVal == 2047) {
		decoded = exponentStyle.Render("special (NaN/Inf)")
	} else {
		decoded = exponentStyle.Render(fmt.Sprintf("%d - %d = %d", expVal, bias, expVal-bias))
	}

	// Build with spacing
	gap := "   "

	var s strings.Builder
	s.WriteString(indent + signLabel + gap + expLabel + strings.Repeat(" ", expTopWidth-len("Exponent")+3) + manLabel + "\n")
	s.WriteString(indent + signTop + gap + expTop + gap + manTop + "\n")
	s.WriteString(indent + signBits + gap + expBits + gap + manBits + "\n")
	s.WriteString(indent + signBot + gap + expBot + gap + manBot + "\n")
	s.WriteString(indent + "     " + decoded)

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
	s.WriteString(currentInputStyle.Render(inputDisplay))
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

	s.WriteString(helpStyle.Render("q: Quit · f: Integer mode · h: Help"))

	return s.String()
}
