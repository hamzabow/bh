---
title: go install Support
status: not started
priority: medium
category: docs
---

## Description

Document the `go install` command in the README so users can install without cloning.

## Requirements

- Add `go install github.com/hamzabow/bh@latest` to README
- Verify the module path in go.mod matches the GitHub repo path

## Testing

- Run `go install github.com/hamzabow/bh@latest` from a clean machine
