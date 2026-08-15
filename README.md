# Forge

Forge is a lightweight Go CLI for scaffolding simple embedded C/C++ projects and generating a basic CMake-based build setup.

## Overview

Forge creates a minimal project structure for new embedded or bare-metal style applications. It currently supports a two-step project setup (`new` then `init`) and a basic CMake build step.

## Commands

| Command | Description |
|---------|-------------|
| `forge new <project_name>` | Create a new project directory and `Forge.toml` |
| `forge init` | Generate starter folders and source/build files |
| `forge build` | Configure and build the project with CMake |
| `forge help` | Print CLI usage information |
| `forge version` | Print the current version string |

`run` and `test` are listed in the help output but are not implemented yet.

## Repository layout

- `main.go` — CLI entry point and command dispatch
- `cmd/generate.go` — `new` and `init` scaffolding logic
- `cmd/build.go` — build command implementation
- `cmd/private.go` — shared config and helpers
- `contents/` — template content used to generate project files
- `go.mod` — Go module definition

## Generated project layout

After running `forge new` and `forge init`, a project looks like this:

```text
my_project/
├── src/
├── include/
├── Forge.toml
├── main.c
├── CMakeLists.txt
└── CMakePresets.json
```

Generated projects include a Hello World `main.c`, a minimal `CMakeLists.txt`, and a `Forge.toml` with project metadata (name, version, toolchain, CMake settings).

## Usage

### Build the CLI

Requires Go 1.26 or later.

```bash
go build -o forge .
```

### Install globally

```bash
sudo install -m 755 forge /usr/local/bin/forge
```

### Create a new project

Project setup is a two-step process:

```bash
forge new demo_app
cd demo_app
forge init
```

1. `forge new` creates the project directory and writes `Forge.toml`.
2. `forge init` generates `src/`, `include/`, `main.c`, `CMakeLists.txt`, and `CMakePresets.json`.

### Build the generated project

From the project directory (requires CMake 3.20+):

```bash
forge build
```

This runs `cmake -S . -B build` and then `cmake --build build`. Output is placed in `build/`.

## Example

```bash
go run . new my_app
cd my_app
go run ../. init
go run ../. build
```

Or, with the binary installed:

```bash
forge new my_app
cd my_app
forge init
forge build
```

## Notes

This repository is an early-stage starter project. Planned additions include embedded workflow tooling, test support, and more advanced build configuration.
