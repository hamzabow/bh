---
title: README Demo Video
status: not started
priority: medium
category: docs
---

## Description

Add a demo video to the README showing the TUI in action. The video should walk through all major features with playback controls (pause, rewind, scrub) — not a GIF, which forces passive viewing with no way to pause or rewind.

## Requirements

- A video (MP4, or hosted on YouTube/asciinema) showing a full feature walkthrough
- Placed near the top of the README for immediate visual impact
- Covers both current features and upcoming features (see demo script below)
- Consider using keystroke overlay software (e.g., KeyCastr, Screenkey, or ShowMeTheKey) to display key presses on screen during recording

## Demo script

The recording should walk through the following scenes in order. Each scene lists what to type and which keys to press.

### Scene 1 — Basic decimal conversion

1. App launches, default state: Decimal mode, 32-bit, Unsigned
2. Type `255`
3. Pause to show all four conversion outputs (Decimal, Hex, Octal, Binary)
4. Pause on the binary visualization with Unicode bracket decorations

### Scene 2 — Hex input with prefix detection

1. Press `Ctrl+W` to clear input
2. Type `0xFF` — tab bar auto-switches to Hex mode
3. Pause to show the conversion results

### Scene 3 — Binary input with prefix detection

1. Press `Ctrl+W` to clear
2. Type `0b11001010` — tab bar auto-switches to Binary mode
3. Pause to show results

### Scene 4 — Octal input with prefix detection

1. Press `Ctrl+W` to clear
2. Type `0o777` — tab bar auto-switches to Octal mode
3. Pause to show results

### Scene 5 — Cycling input base with F1

1. Press `F1` to cycle to Hex mode (input clears)
2. Type `DEADBEEF`
3. Pause to show the binary visualization for a large value
4. Press `F1` to cycle to Octal, then `F1` again to Binary, showing the tab bar updating

### Scene 6 — Bit size cycling with F2

1. Switch back to Decimal mode (`F1` until Decimal)
2. Type `200`
3. Press `F2` to cycle through 8-bit → 16-bit → 32-bit → 64-bit
4. Show the overflow warning appearing at 8-bit (200 > 255 for signed, valid for unsigned)
5. Show how the binary visualization grows/shrinks with bit size

### Scene 7 — Signed mode and negative numbers

1. Press `Ctrl+W` to clear
2. Press `F3` to toggle to Signed mode
3. Press `F2` until 8-bit
4. Type `-42`
5. Pause to show two's complement representation in binary and hex
6. Show the range info updating in the footer
7. Press `F3` to toggle back to Unsigned, show the mode change

### Scene 8 — Cursor navigation and editing

1. Press `Ctrl+W` to clear, switch to Decimal
2. Type `12345`
3. Press `Home` to jump to beginning, then `Right` a couple times
4. Press `Delete` to delete a character at cursor
5. Press `End` to jump to end
6. Press `Backspace` to delete from end

### Scene 9 — Help page

1. Press `?` to open the help page
2. Pause to show the full keyboard shortcuts reference
3. Press `?` to close

### Scene 10 — Clipboard copy (upcoming)

1. Type a value, e.g. `0xCAFE`
2. Press `c` to copy — show the confirmation message
3. Brief note: "Copy any conversion result with a single keypress"

### Scene 11 — Input digit grouping (upcoming)

1. Switch to Decimal mode, type a long number like `1234567890`
2. Press `F4` to toggle digit grouping on
3. Show the input reformatting with visual separators
4. Switch to Binary mode, type a long binary value, show bracket grouping on input
5. Press `F4` to toggle off

### Scene 12 — Float mode / IEEE 754 (upcoming)

1. Press `f` to enter Float mode
2. Type `3.14`
3. Show the IEEE 754 breakdown: sign bit, exponent, mantissa
4. Press `F2` to toggle between 32-bit (single) and 64-bit (double)
5. Show how the binary layout changes

### Scene 13 — CLI one-shot mode (upcoming)

1. Exit the TUI
2. Run `bh 255` in the terminal — show the one-shot output
3. Run `bh 0xFF` — show prefix detection working from CLI
4. Run `bh 0b11111111` — same value, different input base

### Scene 14 — Quit

1. Press `q` to quit

## Notes

- Scenes 10–13 cover features that are specced but not yet implemented — record these after implementation
- For keystroke overlay, Linux options include Screenkey or ShowMeTheKey; macOS has KeyCastr
- Consider asciinema for terminal recording with built-in playback controls, or OBS for a more polished video with overlays
- Keep each scene short (5–10 seconds) with brief pauses on the results
- Total video target: under 2 minutes
