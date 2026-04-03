package main

import tea "github.com/charmbracelet/bubbletea"

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit

		case "h":
			m.showHelp = !m.showHelp
			return m, nil

		case "f":
			if !m.showHelp {
				m.floatMode = !m.floatMode
				m.input = ""
				m.cursor = 0
				m.err = nil
				m.overflow = false
				m.hex, m.binary, m.decimal, m.octal = "", "", "", ""
				m.floatVal, m.floatHex, m.floatBin = "", "", ""
				if m.floatMode && m.bitSize != 32 && m.bitSize != 64 {
					m.bitSize = 32
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
