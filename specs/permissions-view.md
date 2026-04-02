---
title: Permissions View (Unix rwx)
status: not started
priority: medium
category: ui
---

# Permissions View (Unix rwx)

## Description

A toggle ("P" key) that switches digit grouping to a Unix-permissions-oriented view. When active, binary digits are grouped by 3 bits (octal) instead of 4 bits (nibble) / 8 bits (byte), and octal digits are expanded into their binary and rwx equivalents. This helps users learn how Unix permission digits map to read/write/execute flags.

The P toggle is orthogonal to F4 grouping styles. F4 continues to control *how* groups are displayed (Off/Brackets/Spaces/Both), while P controls *what* grouping scheme is used.

## Keybinding

- **P**: Toggle permissions view on/off. Available at all times regardless of input mode.

## Display Behavior

### Binary (input or output) with P active

Groups by 3 bits instead of 4. Each full group shows:
- Top line: brackets around the 3-bit group (controlled by F4)
- Digit line: the binary digits
- Bottom line: brackets with the octal digit (0-7) for each full group (controlled by F4)
- rwx line: the rwx equivalent for each group, shown whenever brackets are shown

Example (F4 = Brackets):

```
     ╭─╮ ╭─╮ ╭─╮
  0b 111 101 101
     ╰7╯ ╰5╯ ╰5╯
     rwx r-x r-x
```

Example (F4 = Spaces):

```
  0b 111 101 101
```

Example (F4 = Both):

```
     ╭─╮ ╭─╮ ╭─╮
  0b 111 101 101
     ╰7╯ ╰5╯ ╰5╯
     rwx r-x r-x
```

### Octal input with P active

Each octal digit is displayed centered over a 3-character-wide column (space-digit-space), so there is room for the binary expansion underneath. Annotations show the 3-bit binary equivalent and rwx for each digit.

Example (F4 = Brackets):

```
  0o  7   5   5
     111 101 101
     rwx r-x r-x
```

Example (F4 = Spaces):

```
  0o  7   5   5
```

Example (F4 = Both):

```
  0o  7   5   5
     111 101 101
     rwx r-x r-x
```

### Hex and decimal input with P active

No change to their input display. The binary output row still receives the permissions treatment.

### Binary output row

The binary output is always visible. When P is active, the binary output row uses 3-bit permissions grouping regardless of the current input mode.

## rwx Mapping

Each 3-bit group maps to rwx as follows:

| Binary | Octal | rwx |
|--------|-------|-----|
| 000    | 0     | --- |
| 001    | 1     | --x |
| 010    | 2     | -w- |
| 011    | 3     | -wx |
| 100    | 4     | r-- |
| 101    | 5     | r-x |
| 110    | 6     | rw- |
| 111    | 7     | rwx |

## Interaction with F4

The F4 cycle (Off/Brackets/Spaces/Both) controls the visual style of the permissions grouping just as it does for the default grouping:

- **Off**: No grouping, no annotations. P has no visible effect.
- **Brackets**: Brackets and annotations shown, no spaces between groups.
- **Spaces**: Spaces between groups, no brackets or annotations.
- **Both**: Brackets, annotations, and spaces between groups.

## Partial Groups

Partial leading groups (when digit count is not a multiple of 3) are displayed without brackets, same as current behavior for partial groups. The rwx line is only shown for full 3-bit groups.

## Tab Bar

When P is active, show an indicator in the status bar, e.g.:

```
[P] Permissions: On
```

## Requirements

- P toggle state persists while the app is running (not saved across sessions).
- Default state on startup: Off.
- Cursor navigation must account for the different spacing in octal-with-P mode (padded digits) and binary-with-P mode (3-bit groups instead of 4-bit).
- The rwx line appears whenever brackets are visible (Brackets or Both modes).
- Binary output row always uses permissions grouping when P is active.

## Subtasks

- [ ] Add `permissionsView bool` field to the model
- [ ] Add P key handler to toggle permissions view
- [ ] Add permissions indicator to the tab bar / status bar
- [ ] Update binary renderer: 3-bit grouping with octal annotation and rwx line when P active
- [ ] Update octal renderer: centered digits with binary expansion and rwx line when P active
- [ ] Update binary output row to use permissions grouping when P active
- [ ] Update cursor position mapping for 3-bit binary groups
- [ ] Update cursor position mapping for padded octal digits
- [ ] Ensure F4 styles (Off/Brackets/Spaces/Both) apply correctly to permissions grouping
- [ ] Ensure partial leading groups display correctly

## Testing

- P key toggles permissions view on and off
- Tab bar shows permissions indicator when active
- Binary input with P + Brackets: 3-bit groups, octal digits in bottom brackets, rwx line shown
- Binary input with P + Spaces: 3-bit groups with spaces, no annotations
- Binary input with P + Both: 3-bit groups with spaces, brackets, octal digits, rwx line
- Binary input with P + Off: no grouping, no annotations
- Octal input with P + Brackets: centered digits, binary expansion, rwx line
- Octal input with P + Spaces: centered digits (padded), no annotations
- Octal input with P + Both: centered digits, binary expansion, rwx line
- Hex/decimal input with P: input display unchanged, binary output row shows permissions grouping
- Binary output row shows 3-bit permissions grouping when P active regardless of input mode
- Cursor navigation works correctly in binary mode with 3-bit groups
- Cursor navigation works correctly in octal mode with padded digits
- Partial leading groups handled correctly (no brackets, no rwx)
- rwx labels match the mapping table for all 8 values (000-111)
