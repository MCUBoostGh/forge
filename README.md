# Forge

Forge is a simple Go-based command-line tool for managing embedded C/C++ projects.

## Overview

This repository contains the `forge` CLI application, built with Go 1.26.0.
It currently supports a minimal command set and is designed to be extended with embedded project workflows such as initialization, building, running, and testing.

## Project Structure

- `main.go` - the CLI entry point and command dispatcher.
- `cmd/build.go` - the build command implementation.
- `go.mod` - Go module definition.

## Available Commands

- `forge build` - run the build command.
- `forge help` - display help information.
- `forge version` - print the current version.

## Example

```bash
go run main.go build
```

## Notes

The build logic is currently a placeholder and should be implemented with project-specific build steps.
