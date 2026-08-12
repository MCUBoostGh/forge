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

## Install Forge to `/usr/bin`

To install the `forge` binary system-wide, build it and move it into `/usr/bin`.

```bash
go build -o forge .
sudo mv forge /usr/bin/
```

If your Go environment is set up for module-aware installs, you can also use:

```bash
go install .
```

Then copy the resulting binary to `/usr/bin`:

```bash
sudo cp $(go env GOPATH)/bin/forge /usr/bin/
```

After installation, run `forge` from any directory:

```bash
forge help
```

## Notes

The build logic is currently a placeholder and should be implemented with project-specific build steps.
