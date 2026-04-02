---
title: Grouping Display Modes
status: done
priority: medium
category: ui
---

# Grouping Display Modes

## Description

Replace the current on/off grouping toggle (F4) with a 4-mode cycle that controls how digit groups are visually separated. The modes are: Off, Brackets, Spaces, Both.

This supersedes the toggle behavior from the [Input Digit Grouping](input-grouping.md) spec. The grouping logic and group boundaries remain the same; only the visual rendering changes.

## Modes

### Off

No grouping. Digits displayed as a flat string. Same as the current "Off" state.

### Brackets

Unicode box-drawing brackets around groups, with annotations where applicable. No spaces between groups. This is equivalent to the current "On" state.

- **Binary**: Nibble brackets (`╭──╮`) on top every 4 bits, byte brackets (`╰──XX──╯`) on bottom every 8 bits with hex annotation.
- **Hex**: Byte brackets (`╰╯`) on bottom every 2 digits.
- **Decimal/Octal**: Bottom brackets (`╰─╯`) every 3 digits.

### Spaces

Spaces between groups at the same boundaries as brackets. No brackets, no annotations.

- **Binary**: A space every 4 bits.
- **Hex**: A space every 2 digits.
- **Decimal/Octal**: A space every 3 digits.

### Both

Brackets and spaces combined. Spaces are inserted between groups, and brackets are drawn around each group. Annotations appear where applicable.

- **Binary**: A space every 4 bits. Nibble brackets (`╭──╮`) span each 4-bit group. Byte brackets (`╰──XX──╯`) span 9 characters (8 digits + 1 internal nibble space) with hex annotation. No bracket spans across the gap between bytes.
- **Hex**: A space every 2 digits. Byte brackets (`╰╯`) under each pair.
- **Decimal/Octal**: A space every 3 digits. Bottom brackets (`╰─╯`) under each group.

## F4 Cycle

F4 cycles through: **Off → Brackets → Spaces → Both → Off**

The tab bar indicator updates to show the current mode:

```
[F4] Grouped: Off | Brackets | Spaces | Both
```

## Requirements

- The grouping boundaries (nibble, byte, 3-digit) remain unchanged from the existing implementation.
- Cursor navigation must account for inserted spaces in Spaces and Both modes. Spaces are visual-only and not part of the raw input string.
- Partial leading groups (fewer digits than a full group) are displayed without brackets or trailing spaces, same as current behavior.
- The mode state persists while the app is running (not saved across sessions).
- Default mode on startup: Off.

## Subtasks

- [ ] Replace `groupedInput bool` with a grouping mode field (enum/int cycling through Off, Brackets, Spaces, Both)
- [ ] Update F4 handler to cycle through the four modes
- [ ] Update tab bar to display the current mode name
- [ ] Update binary renderer to support Spaces mode (space every 4 bits, no brackets/annotations)
- [ ] Update binary renderer to support Both mode (spaces + brackets, byte bracket spans 9 chars)
- [ ] Update hex renderer to support Spaces mode (space every 2 digits, no brackets)
- [ ] Update hex renderer to support Both mode (spaces + brackets)
- [ ] Update decimal/octal renderer to support Spaces mode (space every 3 digits, no brackets)
- [ ] Update decimal/octal renderer to support Both mode (spaces + brackets)
- [ ] Update cursor position mapping to account for inserted spaces in Spaces and Both modes
- [ ] Ensure arrow keys skip over visual spaces naturally

## Testing

- F4 cycles through Off, Brackets, Spaces, Both, and back to Off
- Tab bar shows the correct mode label at each step
- Brackets mode matches current behavior exactly (no regression)
- Spaces mode: binary shows spaces every 4 bits with no brackets or annotations
- Spaces mode: hex shows spaces every 2 digits with no brackets
- Spaces mode: decimal/octal show spaces every 3 digits with no brackets
- Both mode: binary shows spaces every 4 bits, nibble brackets span 4 chars, byte brackets span 9 chars (8 digits + 1 space) with hex annotation
- Both mode: hex shows spaces every 2 digits with byte brackets
- Both mode: decimal/octal show spaces every 3 digits with bottom brackets
- Cursor navigation works correctly in all modes and all bases
- Partial leading groups display correctly in all modes
- No bracket spans across the gap between byte groups in Both mode for binary
