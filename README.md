# bh - Binary Hex Converter

A Terminal User Interface (TUI) for converting between number bases and visualizing IEEE 754 floats, built with Go and Charm Bracelet's Bubbletea framework.

## Features

- **Multi-base conversion**: Convert between decimal, hexadecimal, octal, and binary formats
- **IEEE 754 float mode**: Visualize 32-bit and 64-bit floating-point representation with sign, exponent, and mantissa breakdown
- **Bit size customization**: Support for 8, 16, 32, and 64-bit integers
- **Signed/Unsigned modes**: Handle both signed and unsigned integer representations
- **Visual binary formatting**: Binary numbers displayed with Unicode decorations showing hex values under byte groups
- **Unix permissions view**: Octal/rwx permission visualization for binary and octal values
- **Digit grouping**: Multiple display modes (brackets, spaces, or both)
- **Overflow detection**: Alerts when values exceed the selected bit range
- **Prefix detection**: Type `0x`, `0b`, or `0o` to auto-switch input base
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

- **F**: Toggle float mode (IEEE 754)
- **T**: Cycle input base (Decimal → Hex → Octal → Binary)
- **W**: Cycle bit size (8 → 16 → 32 → 64)
- **S**: Toggle Signed/Unsigned mode
- **G**: Cycle digit grouping (Off → Brackets → Spaces → Both)
- **P**: Toggle Unix permissions view
- **H**: Show help page
- **q** or **Ctrl+C**: Quit

### Input Formats

- **Decimal**: Standard decimal numbers (e.g., `255`, `1024`)
- **Hexadecimal**: Hex digits (e.g., `FF`, `A0B1`, or prefix with `0x`)
- **Octal**: Octal digits (e.g., `377`, or prefix with `0o`)
- **Binary**: Binary digits (e.g., `11111111`, or prefix with `0b`)
- **Float mode**: Decimal floats (e.g., `1.5`, `3.14e2`, `nan`, `inf`)

## Example

When you input a number, the application displays:
- The original input
- Conversions to all four bases
- Current bit size and signed/unsigned mode
- Valid range for the current configuration
- Overflow warnings when applicable

Binary output includes visual formatting with Unicode box drawing characters showing hex values under each byte group.

In float mode, the binary representation is broken down into color-coded sign, exponent, and mantissa fields with decoded values.

## Dependencies

- [Bubbletea](https://github.com/charmbracelet/bubbletea) - TUI framework
- [Lipgloss](https://github.com/charmbracelet/lipgloss) - Styling library

## License

This project is open source and available under the MIT License.
