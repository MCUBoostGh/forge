# Forge

Forge is a lightweight Go CLI for scaffolding simple embedded C/C++ projects and generating a basic CMake-based build setup.

## Overview

This project creates a minimal project structure for new embedded or bare-metal style applications. It currently supports project initialization and a basic build step through CMake.

## Current functionality

- `forge init <project_name>` creates a new project directory with the common starter files.
- `forge build` configures and builds the project using CMake.
- `forge help` prints the CLI usage information.
- `forge version` prints the current version string.

## Project structure

- `main.go` - CLI entry point and command dispatch.
- `cmd/init.go` - project scaffolding logic.
- `cmd/build.go` - build command implementation.
- `cmd/private.go` - shared internal helpers.
- `contents/` - template content used to generate project files.
- `go.mod` - Go module definition.

## Generated project layout

Running:

```bash
forge init my_project
```

creates a structure similar to:

```text
my_project/
├── src/
├── include/
├── Forge.toml
├── main.c
├── CMakeLists.txt
├── CMakePresets.json
```

## Usage

### Build the CLI

```bash
go build -o forge .
```

### Install globally

```bash
sudo install -m 755 forge /usr/local/bin/forge
```

### Initialize a new project

```bash
forge init demo_app
cd demo_app
```

### Build the generated project

```bash
forge build
```

This runs CMake in the current project directory and then builds the output in `build/`.

## Example

```bash
go run . init my_app
cd my_app
forge build
```

## Notes

This repository is still a small starter project and is intended to evolve with additional tooling for embedded workflows, test support, and more advanced build configuration.
