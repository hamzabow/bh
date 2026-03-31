---
title: Octal Base
status: done
priority: "-"
category: feature
---

## Description

Add octal (base 8) as a number base alongside decimal, hex, and binary.

## Requirements

- Octal appears in the tab cycle: Decimal -> Hex -> Octal -> Binary
- Input validation accepts digits 0-7
- Octal conversion displayed in output between hex and binary

## Testing

- Switch to octal input mode, type "77", verify decimal shows 63 and hex shows 3F
- Enter a decimal like 255, verify octal shows 377
