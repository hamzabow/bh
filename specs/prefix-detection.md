---
title: Prefix Detection
status: not started
priority: low
category: feature
---

## Description

Auto-detect the input base when the user types a prefix like 0x, 0b, or 0o.

## Requirements

- Typing "0x" switches to hex mode automatically
- Typing "0b" switches to binary mode
- Typing "0o" switches to octal mode
- The prefix is consumed and not shown as part of the input value

## Testing

- Type "0xFF" in decimal mode, verify it switches to hex and shows FF
- Type "0b1010", verify it switches to binary
- Type "0o77", verify it switches to octal
