---
title: Keyboard Shortcuts & Help Page
status: done
priority: low
category: feature
---

## Description

Add familiar keyboard shortcuts for cursor navigation and input editing, plus a help page that documents all keybindings.

## Requirements

### Cursor movement

- Home moves cursor to the beginning of input
- End moves cursor to the end of input
- Ctrl+Left moves cursor to the beginning of input (equivalent to Home — single word)
- Ctrl+Right moves cursor to the end of input (equivalent to End — single word)

### Editing

- Delete (forward delete) removes the character under the cursor
- Ctrl+W clears all input

### Help page

- `?` toggles a help page overlay showing all keybindings
- The help page lists: F1/F2/F3 mode switches, arrow keys, Home/End, Ctrl+Left/Right, Delete, Ctrl+W, `?` for help, q/Ctrl+C to quit
- The main view's bottom help line should hint that `?` opens help (e.g. `q: Quit · ?: Help`)

## Subtasks

- [ ] Add Home/End key handling
- [ ] Add Ctrl+Left / Ctrl+Right handling (same as Home/End)
- [ ] Add Delete (forward delete) handling
- [ ] Add Ctrl+W to clear input
- [ ] Add `showHelp` state to model
- [ ] Add help page view with all keybindings
- [ ] Toggle help page with `?`
- [ ] Update bottom help line to show `?` hint

## Testing

- Type a number, press Home, verify cursor is at position 0
- Press End, verify cursor is at the end
- Press Ctrl+Left / Ctrl+Right, verify same behavior as Home/End
- Place cursor in the middle, press Delete, verify character at cursor is removed
- Type a number, press Ctrl+W, verify input is cleared
- Press `?`, verify help page appears with all keybindings listed
- Press `?` again, verify help page closes and main view returns
- While help page is shown, press q or Ctrl+C, verify app quits
