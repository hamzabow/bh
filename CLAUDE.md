# CLAUDE.md

## Project overview

bh is a Terminal User Interface (TUI) for converting between number bases (decimal, hexadecimal, octal, binary), built with Go and Bubbletea.

## Build and run

```bash
go build          # build the binary
./bh              # run the TUI
```

## Code structure

- `main.go` — the entire application (single-file TUI)

## Specs directory

Feature specs live in `specs/`. Each spec is a markdown file with YAML frontmatter (title, status, priority, category) followed by sections for description, requirements, subtasks, and testing.

`specs/README.md` contains an overview table of all features and their status.

When the user asks for new features:

1. Switch to plan mode before creating or updating specs
2. Create or update spec files in `specs/` first
3. Ask clarifying questions to fill in requirements and testing details
4. Update `specs/README.md` with the new entries
5. Only implement after the spec is agreed upon

When a feature is completed, update its frontmatter status to `done` and update the overview table in `specs/README.md`.
