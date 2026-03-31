---
title: Negative Number Input
status: not started
priority: medium
category: feature
---

## Description

Allow typing a minus sign in signed mode to input negative numbers directly.

## Requirements

- In signed mode, accept "-" as the first character in decimal input
- Only allow one minus sign, and only at position 0
- Not applicable for hex/octal/binary input modes

## Testing

- Switch to signed mode, type "-128" in 8-bit, verify it converts correctly
- Verify minus sign is rejected in unsigned mode
- Verify minus sign cannot be typed in the middle of a number
