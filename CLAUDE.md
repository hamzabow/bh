# CLAUDE.md

## Project overview

bh is a Terminal User Interface (TUI) for converting between number bases (decimal, hexadecimal, octal, binary), built with Go and Bubbletea.

## Build and run

```bash
go build          # build the binary
./bh              # run the TUI
```

## Code structure

- `main.go` — model struct, types, initialization, entry point
- `styles.go` — all lipgloss style definitions
- `common.go` — shared utilities (tab bar, cursor, digit grouping)
- `update.go` — top-level Update() dispatcher
- `view.go` — top-level View() dispatcher and help page
- `integer.go` — integer mode: key handling, conversions, rendering
- `float.go` — IEEE 754 float mode: key handling, parsing, rendering
