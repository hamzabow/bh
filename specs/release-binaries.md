---
title: Release Binaries
status: not started
priority: low
category: infra
---

## Description

Set up GitHub Actions to automatically build and publish binaries for each platform when a release is tagged.

## Requirements

- Build for Linux, macOS, and Windows (amd64 and arm64)
- Trigger on new git tags (e.g. v1.0.0)
- Upload binaries as GitHub Release assets

## Subtasks

- [ ] Create GitHub Actions workflow
- [ ] Configure GoReleaser or equivalent
- [ ] Test with a tag push

## Testing

- Push a tag, verify binaries appear on the GitHub Releases page
- Download and run on each target platform
