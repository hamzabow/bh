package main

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit

		case "h":
			m.showHelp = !m.showHelp
			return m, nil

		case "ctrl+f":
			if !m.showHelp {
				if m.floatMode {
					// Float → Integer
					m.floatMode = false
					m = m.convertFloatToInteger()
				} else {
					// Integer → Float
					m.floatMode = true
					if m.bitSize != 32 && m.bitSize != 64 {
						m.bitSize = 32
					}
					m = m.convertIntegerToFloat()
				}
				return m, nil
			}
		}

		if m.showHelp {
			return m, nil
		}

		if m.floatMode {
			m = m.updateFloatKeys(msg)
		} else {
			m = m.updateIntegerKeys(msg)
		}
	}

	return m, nil
}

func (m model) convertIntegerToFloat() model {
	if m.decimal == "" || m.err != nil || m.overflow {
		m.input = ""
		m.cursor = 0
		m.err = nil
		m.overflow = false
		m.hex, m.binary, m.decimal, m.octal = "", "", "", ""
		m.floatVal, m.floatHex, m.floatBin = "", "", ""
		return m
	}

	m.input = m.decimal
	m.cursor = len(m.input)
	m.err = nil
	m.overflow = false
	m.hex, m.binary, m.decimal, m.octal = "", "", "", ""
	m = m.updateFloatConversions()
	return m
}

func (m model) convertFloatToInteger() model {
	if m.input == "" || m.err != nil || m.floatVal == "" {
		m.input = ""
		m.cursor = 0
		m.err = nil
		m.overflow = false
		m.hex, m.binary, m.decimal, m.octal = "", "", "", ""
		m.floatVal, m.floatHex, m.floatBin = "", "", ""
		return m
	}

	inputLower := strings.ToLower(m.input)
	var val float64
	switch inputLower {
	case "nan", "inf", "+inf", "infinity", "+infinity", "-inf", "-infinity":
		// No integer representation, clear
		m.input = ""
		m.cursor = 0
		m.err = nil
		m.overflow = false
		m.hex, m.binary, m.decimal, m.octal = "", "", "", ""
		m.floatVal, m.floatHex, m.floatBin = "", "", ""
		return m
	default:
		var err error
		val, err = strconv.ParseFloat(m.input, 64)
		if err != nil || math.IsNaN(val) || math.IsInf(val, 0) {
			m.input = ""
			m.cursor = 0
			m.err = nil
			m.overflow = false
			m.hex, m.binary, m.decimal, m.octal = "", "", "", ""
			m.floatVal, m.floatHex, m.floatBin = "", "", ""
			return m
		}
	}

	truncated := math.Trunc(val)
	m.input = fmt.Sprintf("%.0f", truncated)
	m.cursor = len(m.input)
	m.inputType = "decimal"
	m.err = nil
	m.overflow = false
	m.floatVal, m.floatHex, m.floatBin = "", "", ""
	m = m.updateConversions()
	return m
}
