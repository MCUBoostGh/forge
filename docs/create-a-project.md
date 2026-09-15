# Create a project

Forge uses a two-step setup, then a CMake-backed build.

```text
forge new  →  Forge.toml
forge init →  sources + CMake
forge build → binaries
```

## 1. Create the project

```bash
forge new <name> [--device <device>]
```

Examples:

```bash
forge new blink --device stm32f103r8
forge new blink --device bluepill
forge new demo_app
```

`new` creates a directory named `<name>` and writes `Forge.toml` with project metadata (name, version, toolchain, CMake settings). With `--device`, Forge records a target device for STM32-oriented workflows.

Then enter the project:

```bash
cd <name>
```

## 2. Initialize the tree

```bash
forge init
```

`init` reads `Forge.toml` from the current directory and generates:

```text
.
├── src/
├── include/
├── cmake/
├── Forge.toml
├── main.c
├── CMakeLists.txt
└── CMakePresets.json
```

Generated projects include a starter `main.c` and a minimal CMake + presets setup.

## 3. Build

```bash
forge build
```

Forge loads the build type from `Forge.toml` (defaulting to `debug` when unset) and runs the matching CMake preset.

Build output lands under the preset’s binary directory (typically under `build/`).

## Without installing Forge

From a Forge checkout:

```bash
go run . new my_app --device stm32f103r8
cd my_app
go run ../. init
go run ../. build
```

Or use the repo helper:

```bash
./run.sh my_app --device stm32f103r8
```

## Tips

- Run `init` and `build` from the project directory (where `Forge.toml` lives).
- Prefer inspecting and editing the generated CMake files — Forge aims to stay transparent, not hide the build system.
- Device support is still early; see [ROADMAP.md](../ROADMAP.md) for the STM32 bootstrap plan.

## See also

- [Getting started](getting-started.md)
- [Commands](commands.md)
