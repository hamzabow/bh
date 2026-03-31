---
title: Prefix Detection
status: done
priority: low
category: feature
---

## Description

Support standard number prefixes (`0x`, `0b`, `0o`) as first-class input. When a recognized prefix is typed, the input mode switches automatically. The prefix remains visible in the input field — it is part of the number, not a hidden shortcut.

## Requirements

### Prefix recognition

- `0x` or `0X` — switch to hex mode
- `0b` or `0B` — switch to binary mode
- `0o` or `0O` — switch to octal mode

### Works from any mode

Prefixes are detected regardless of the current input mode. Typing `0x` while in binary mode switches to hex. The rationale: prefixes are a standard way to express numbers and should behave consistently.

### Input preservation

The prefix stays in the input field as typed. `0xFF` is displayed as `0xFF`, not `FF`. Parsing strips the prefix internally.

### Character validation after prefix

Once a prefix is active, subsequent characters are validated against the target mode (e.g., after `0x`, hex digits are accepted).

### Backspace behavior

Deleting the second character of a prefix (e.g., backspacing `x` from `0x`) reverts the mode to decimal (the default). The remaining `0` is valid in any numeric mode.

## Subtasks

- [ ] Detect prefix when second character is typed at position 1 after `0`
- [ ] Switch input mode on prefix detection
- [ ] Allow valid characters for the new mode after the prefix
- [ ] Strip prefix before parsing in updateConversions
- [ ] Revert to decimal mode when prefix is backspaced away
- [ ] Handle case-insensitive prefixes

## Testing

- Type `0xFF` from decimal mode — switches to hex, shows `0xFF`, converts correctly
- Type `0b1010` from hex mode — switches to binary, shows `0b1010`, converts correctly
- Type `0o77` from binary mode — switches to octal, shows `0o77`, converts correctly
- Backspace the `x` from `0x` — reverts to decimal, input is `0`
- Type `0X1A` — uppercase prefix works identically
- Existing non-prefixed input continues to work unchanged
