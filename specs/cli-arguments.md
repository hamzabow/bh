---
title: CLI Arguments
status: not started
priority: medium
category: feature
---

## Description

Support one-shot conversions from the command line without entering the TUI.

## Requirements

- `bh 255` converts decimal 255 and prints all bases
- `bh 0xFF` detects hex prefix and converts
- `bh 0b1010` detects binary prefix and converts
- `bh 0o77` detects octal prefix and converts
- If no arguments, launch TUI as usual

## Testing

- Run `bh 255` and verify output shows hex, octal, binary
- Run `bh 0xFF` and verify it detects hex
- Run `bh` with no args, verify TUI launches
