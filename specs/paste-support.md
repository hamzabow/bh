---
title: Paste Support
status: not started
priority: low
category: feature
---

## Description

Allow pasting values from the clipboard into the input field.

## Requirements

- Ctrl+V pastes clipboard content into the input at cursor position
- Invalid characters for the current input type are filtered out
- Conversion updates immediately after paste

## Testing

- Copy "FF" to clipboard, switch to hex mode, paste, verify conversion
- Paste a string with invalid characters, verify only valid ones are inserted
