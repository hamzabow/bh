---
title: Float Mode (IEEE 754)
status: done
priority: high
category: feature
---

## Description

Add a separate float mode for IEEE 754 floating-point visualization. Toggle between integer mode and float mode with a key press.

## Requirements

- Press "f" to toggle between integer mode and float mode
- Float mode supports 32-bit (single) and 64-bit (double) precision
- Display sign bit, exponent, and mantissa separately
- Show the decimal value, hex representation, and binary layout
- Visual breakdown of the binary representation showing each field

## Subtasks

- [ ] Add mode field to model (integer vs float)
- [ ] Implement IEEE 754 encoding/decoding
- [ ] Build float-specific view rendering
- [ ] Handle special values (NaN, Infinity, -0)
- [ ] Update help text to show mode toggle

## Testing

- Enter 1.5 in float mode, verify correct IEEE 754 binary
- Check special values: 0, -0, Infinity, NaN
- Toggle between integer and float mode, verify UI switches cleanly
