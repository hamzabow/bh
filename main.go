package main

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"
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
	floatMode  bool
	floatVal   string
	floatHex   string
	floatBin   string
}

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

func main() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
