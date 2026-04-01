---
title: Input Digit Grouping
status: not started
priority: medium
category: ui
---

# Input Digit Grouping

## Description

Add an optional visual grouping mode for the input field that formats digits as they are typed, improving readability for long numbers. Toggled with F4, off by default.

The grouping style varies by input base to match the natural structure of each number system. Binary and hex get full bracket treatment (similar to the binary output display), while decimal and octal get lightweight space separators.

## Requirements

### Binary input grouping

Full two-level bracket treatment matching the existing binary output format:

- **Nibble grouping (4 bits)**: top-level brackets `╭──╮` above each group of 4 digits
- **Byte grouping (8 bits)**: bottom-level brackets `╰──╯` below each group of 8 digits, annotated with the 2-digit hex value of that byte

Example with 16 bits typed (`1010001111001101`):

```
  ╭──╮╭──╮ ╭──╮╭──╮
  10100011 11001101
  │      │ │      │
  ╰──A3──╯ ╰──CD──╯
```

### Hex input grouping

Bracket treatment at the byte level (every 2 hex digits):

- Top bracket `╭╮` above each pair
- Bottom bracket `╰╯` below each pair, annotated with the decimal byte value (0-255)

Example with 4 hex digits typed (`A3CD`):

```
  ╭──╮ ╭──╮
  A3   CD
  │  │ │  │
  ╰163╯ ╰205╯
```

### Decimal input grouping

Space separators every 3 digits from the right (standard thousands grouping). No brackets.

Example: `1 234 567`

### Octal input grouping

Space separators every 3 digits from the right, same as decimal. No brackets — octal digits don't align to byte boundaries so bracket annotations would be misleading.

Example: `17 777 777 777`

### Dynamic behavior

Grouping adapts to the current digit count:

- Brackets and separators only appear where there are enough digits to form a complete group
- Partial groups at the leading (leftmost) end are displayed without brackets
- As the user types or deletes, the grouping updates in real time

### Cursor and editing

The input remains fully editable while grouping is active:

- Visual separators and brackets are not part of the actual input string
- The cursor position maps from raw string index to display position, accounting for inserted separator characters
- Left/right arrow keys skip over visual separators naturally
- Backspace and character insertion work on the raw string; the display re-renders with updated grouping

### Toggle

- **F4** toggles grouped input display on/off
- Off by default
- Add to the tab bar area: `[F4] Grouped: On | Off`
- State persists while the app is running (not saved across sessions)

## Subtasks

- [ ] Add `groupedInput` bool field to model
- [ ] Handle F4 keypress to toggle grouping
- [ ] Implement raw-to-display position mapping for cursor placement
- [ ] Implement binary input grouping renderer (nibble + byte brackets with hex annotation)
- [ ] Implement hex input grouping renderer (byte brackets with decimal annotation)
- [ ] Implement decimal input grouping renderer (space every 3 digits)
- [ ] Implement octal input grouping renderer (space every 3 digits)
- [ ] Render the grouped input in the View when toggled on
- [ ] Add F4 indicator to the tab bar / help area
- [ ] Ensure left/right arrow keys move correctly through grouped display
- [ ] Ensure backspace/insert work correctly with grouping active

## Testing

- Binary grouping shows correct nibble and byte brackets with accurate hex annotations
- Hex grouping shows correct byte brackets with accurate decimal annotations
- Decimal and octal show space separators at correct positions
- Cursor navigation works correctly in all bases with grouping on
- Typing and deleting digits updates grouping in real time
- Grouping handles partial leading groups correctly (e.g., 5 binary digits: no full byte bracket)
- F4 toggles grouping on and off, display updates immediately
- Grouping off by default on app start
- No interference with prefix detection (0x, 0b, 0o)
- Works correctly across all bit sizes (8, 16, 32, 64)
