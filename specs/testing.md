---
title: Testing
status: not started
priority: high
category: infra
---

## Description

Testing strategy for bh, covering three layers: unit tests for pure functions, model tests for Bubbletea state transitions, and CLI integration tests. Includes prerequisite refactoring to make the single-file codebase testable.

## Requirements

### Refactoring (prerequisite)

- Extract pure functions from `main.go` into a separate file (e.g. `conversion.go`)
- Export functions that need testing (capitalize names)
- Keep Bubbletea model, `Update`, and `View` logic in `main.go`
- Functions to extract and export:
  - `hasPrefix` — prefix detection (`0x`, `0b`, `0o`)
  - `isValidChar` — character validation per input type
  - `getBitLimits` — returns max/min values for a given bit size
  - `toTwosComplement` / `fromTwosComplement` — signed number encoding
  - `groupDigits` — splits digit string into groups
  - `formatBinaryWithBytes` / `formatHexWithBytes` — output formatting with brackets
  - `getNumFromBinary` — binary string to int64 parsing

### Layer 1: Unit tests (pure functions)

Use table-driven tests for each exported function. File: `conversion_test.go`.

- `HasPrefix` — all prefix variants (`0x`, `0X`, `0b`, `0B`, `0o`, `0O`), non-prefix input, empty string
- `IsValidChar` — valid/invalid chars for each input type (decimal, hex, octal, binary), minus sign only allowed at position 0 in signed mode
- `GetBitLimits` — correct (maxUnsigned, maxSigned, minSigned) for 8, 16, 32, 64 bit
- `ToTwosComplement` / `FromTwosComplement` — round-trip correctness, boundary values (-1, -128, min signed for each bit size)
- `GroupDigits` — various group sizes (3 for decimal, 4 for binary), partial groups, single digit, empty input
- `FormatBinaryWithBytes` / `FormatHexWithBytes` — output matches expected bracket formatting
- `GetNumFromBinary` — correct parsing for various binary strings

### Layer 2: Model tests (Bubbletea Update)

Test the model as a state machine by calling `Update()` with synthetic `tea.KeyMsg` values and asserting on the resulting model fields. File: `model_test.go`.

- **Character input:** valid chars accepted and appended to `input`, invalid chars rejected (model unchanged)
- **Prefix detection:** typing `0` then `x` switches `inputType` to hex; `0b` → binary; `0o` → octal
- **Cursor navigation:** Left/Right move by 1 and respect bounds, Home jumps to 0, End jumps to `len(input)`
- **Delete/Backspace:** character removal at correct position, prefix reversion when prefix is broken by deletion
- **Ctrl+W:** clears input back to prefix boundary
- **Mode toggling:** F1 cycles input type (dec→hex→oct→bin→dec), F2 cycles bit size (8→16→32→64→8), F3 toggles signed/unsigned, F4 toggles grouping
- **Conversion accuracy:** spot-check known values across all bases and bit sizes (e.g. input `255` in decimal → hex `FF`, binary `11111111`, octal `377`)
- **Overflow detection:** `overflow` flag set at bit boundaries (e.g. `256` in 8-bit unsigned, `128` in 8-bit signed)
- **Signed mode:** negative values produce correct two's complement display
- **Help toggle:** `?` key toggles `showHelp`

### Layer 3: CLI integration tests

Depends on the [CLI Arguments](cli-arguments.md) spec being implemented. File: `cli_test.go`.

- Build the binary, then run it with arguments and check stdout
- `bh 255` — output contains hex `FF`, binary `11111111`, octal `377`
- `bh 0xFF` — detects hex prefix, output contains decimal `255`
- `bh 0b1010` — detects binary, output contains decimal `10`
- `bh 0o77` — detects octal, output contains decimal `63`
- Invalid input (e.g. `bh zzz`) — prints an error message, exits non-zero

## Subtasks

- [ ] Extract pure functions to `conversion.go` and export them
- [ ] Create `conversion_test.go` with table-driven unit tests (Layer 1)
- [ ] Create `model_test.go` with Bubbletea Update tests (Layer 2)
- [ ] Create `cli_test.go` with integration tests (Layer 3, after CLI Arguments spec is done)

## Testing

- `go test ./...` runs all layers
- Layer 3 requires the binary to be built first (`go build`)
- All tests should pass in CI with no TUI or terminal required
