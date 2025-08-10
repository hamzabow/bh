# BH - Binary Hex Converter

A Terminal User Interface (TUI) for converting between Hexadecimal, Binary, and Decimal numbers built with Go and Charm Bracelet's Bubbletea framework.

## Features

- **Multi-base conversion**: Convert between decimal, hexadecimal, and binary formats
- **Bit size customization**: Support for 8, 16, 32, and 64-bit numbers
- **Signed/Unsigned modes**: Handle both signed and unsigned integer representations
- **Visual binary formatting**: Binary numbers displayed with Unicode decorations showing hex values under byte groups
- **Overflow detection**: Alerts when values exceed the selected bit range
- **Interactive TUI**: Clean, responsive terminal interface with keyboard shortcuts

## Installation

```bash
git clone <repository-url>
cd bh
go build
```

## Usage

Run the application:

```bash
./bh
```

### Controls

- **Tab**: Cycle between input types (Decimal → Hex → Binary)
- **Shift+Tab**: Cycle between bit sizes (8 → 16 → 32 → 64)
- **S**: Toggle between Signed/Unsigned modes
- **q** or **Ctrl+C**: Quit

### Input Formats

- **Decimal**: Standard decimal numbers (e.g., `255`, `1024`)
- **Hexadecimal**: Hex digits without prefix (e.g., `FF`, `A0B1`)
- **Binary**: Binary digits (e.g., `11111111`, `1010`)

## Example

When you input a number, the application displays:
- The original input
- Conversions to all three formats
- Current bit size and signed/unsigned mode
- Valid range for the current configuration
- Overflow warnings when applicable

Binary output includes visual formatting with Unicode box drawing characters showing hex values under each byte group.

## Dependencies

- [Bubbletea](https://github.com/charmbracelet/bubbletea) - TUI framework
- [Lipgloss](https://github.com/charmbracelet/lipgloss) - Styling library

## License

This project is open source and available under the MIT License.
