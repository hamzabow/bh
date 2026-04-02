---
title: CLI Arguments
status: not started
priority: medium
category: feature
---

## Description

Support one-shot conversions from the command line without entering the TUI, and allow filtering which bases are displayed in both CLI and TUI modes.

## Requirements

### One-shot conversion

- `bh 255` converts decimal 255 and prints all bases
- `bh 0xFF` detects hex prefix and converts
- `bh 0b1010` detects binary prefix and converts
- `bh 0o77` detects octal prefix and converts
- If no arguments, launch TUI as usual

### Output filtering (`--show`)

- `--show` flag accepts a comma-separated list of bases to display
- Valid values: `dec`, `hex`, `oct`, `bin`
- When omitted, all bases are shown (default behavior)
- Examples:
  - `bh --show hex,bin 255` — only print hex and binary
  - `bh --show dec 0xFF` — only print the decimal conversion
- In TUI mode, `--show` sets the initially visible output rectangles:
  - `bh --show hex,bin` — launch TUI with only hex and binary outputs visible

### Input base override (`--input`)

- `--input` flag explicitly sets the input base, overriding prefix auto-detection
- Valid values: `dec`, `hex`, `oct`, `bin`
- Examples:
  - `bh --input bin --show hex 1010` — interpret input as binary, show only hex output
  - `bh --input hex FF` — interpret as hex without requiring `0x` prefix
- In TUI mode, `--input` sets the initial input mode

### TUI interactive toggling (open design question)

- It would be useful to toggle individual outputs on/off while in the TUI
- Design TBD — need a good UX for this (keybinding scheme, visual feedback, etc.)
- The `--show` flag covers the initial state; interactive toggling is a stretch goal

## Subtasks

1. [ ] Add `--show` and `--input` flag parsing (e.g. using `flag` stdlib)
2. [ ] Implement one-shot CLI mode (print and exit)
3. [ ] Wire `--show` into TUI to control which output rectangles render
4. [ ] Wire `--input` into TUI to set initial input mode
5. [ ] Design and implement interactive output toggling in TUI (stretch)

## Testing

- Run `bh 255` and verify output shows all bases
- Run `bh 0xFF` and verify it detects hex
- Run `bh` with no args, verify TUI launches
- Run `bh --show hex,bin 255` and verify only hex and binary are printed
- Run `bh --show dec 0xFF` and verify only decimal is printed
- Run `bh --input bin 1010` and verify it interprets as binary
- Run `bh --input hex FF` and verify it interprets as hex
- Run `bh --show hex` (no value) and verify TUI launches with only hex output visible
- Run `bh --show invalid` and verify a helpful error message is shown
