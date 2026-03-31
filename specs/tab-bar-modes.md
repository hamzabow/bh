---
title: Tab Bar Mode Display
status: done
priority: medium
category: ui
---

# Tab Bar Mode Display

## Description

Replace the current single-label mode indicators with tab-bar-style displays so users can see all available options, know which is active, and understand where cycling will take them.

Currently, the input type (Tab), bit size (Shift+Tab), and signed/unsigned (S) modes display only the active value with no indication of what other options exist. This makes cycling feel blind — the user doesn't know how many presses it takes to reach a desired mode.

## Requirements

### Tab bars for all three mode types

Each mode type gets a horizontal tab bar showing all options with the active one highlighted:

1. **Input type**: Decimal | Hex | Octal | Binary
2. **Bit size**: 8-bit | 16-bit | 32-bit | 64-bit
3. **Signed mode**: Unsigned | Signed

### Keybindings

Replace Tab/Shift+Tab/S with function keys to avoid conflicts with character input (hex digits like `b` and future prefix detection using `0x`, `0b`, `0o`):

- **F1** — cycle input type
- **F2** — cycle bit size
- **F3** — cycle signed/unsigned

Note: these keybindings are easy to change later if a better scheme is found.

### Visual style

- Segment/pill style with separators: `Decimal | Hex | Octal | Binary`
- Active segment: bold + accent color (from existing palette)
- Inactive segments: dimmed foreground (e.g. color 241)

### Layout

- Input type tab bar at the top (below the title)
- Bit size tab bar + signed/unsigned tab bar on the next line
- Range info below that
- Input field below range info
- Help line at the bottom shows the F1/F2/F3 keybinds

## Subtasks

- [ ] Create styles for active tab, inactive tab, and separator
- [ ] Render input type tab bar replacing the `Input (Type):` label
- [ ] Render bit size tab bar replacing the `N-bit` badge
- [ ] Render signed/unsigned as a tab bar replacing the inline text
- [ ] Change keybindings from Tab/Shift+Tab/S to F1/F2/F3
- [ ] Update help line at the bottom with new keybinds
- [ ] Verify input clears when switching input type (existing behavior preserved)
- [ ] Verify conversions update when switching bit size or signed mode (existing behavior preserved)

## Testing

- All options visible at all times for all three mode types, with correct one highlighted
- Cycling with F1/F2/F3 updates the highlight correctly
- Existing conversion behavior unchanged
- F-keys don't interfere with normal character input
- Visual appearance reasonable at typical terminal widths
