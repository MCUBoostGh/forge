# Create a project

Forge uses a two-step setup, then a CMake-backed build.

```text
forge new  →  Forge.toml
forge init →  CMSIS cache + sources + CMake
forge build → ELF / HEX / BIN
```

## 1. Create the project

```bash
forge new <name> --device <device>
```

Examples:

```bash
forge new blink --device stm32f103r8
```

`--device` is required. `Lookup` is by catalog id (`stm32f103r8`); aliases such as `bluepill` are not wired yet.

`new` creates `<name>/` and writes `Forge.toml` with project metadata and a default dependency:

```toml
dependencies = ["cmsis5@5.9.0"]
```

Change the spec (for example `cmsis6@6.1.0`) before `init` if you want the CMSIS 6 catalog entry. Version is taken from this list; URLs live in Forge’s package YAML, not in the project.

Then enter the project:

```bash
cd <name>
```

## 2. Initialize the tree

```bash
forge init
```

`init` reads `Forge.toml` from the current directory, fetches each `name@version` into the user cache if missing, and generates:

```text
.
├── src/
├── include/
├── cmake/
│   ├── gcc-arm-none-eabi.cmake
│   └── Package.cmake
├── Forge.toml
├── LinkerScript.ld
├── main.c
├── CMakeLists.txt
└── CMakePresets.json
```

`main.c` is a Cortex-M CMSIS smoke test (not Hello World). `cmake/Package.cmake` defines an INTERFACE library per dependency and points include paths at the cache. `CMakeLists.txt` links those targets with `PRIVATE`.

CMSIS is **not** copied into the project. Shared layout:

```text
~/.cache/forge/packages/CMSIS_5-5.9.0/CMSIS/Core/Include
```

A later project that uses `cmsis5@5.9.0` reuses that directory.

The linker may warn about a missing `Reset_Handler` until startup files exist.

## 3. Build

```bash
forge build
```

Forge loads the build type from `Forge.toml` (defaulting to `debug` when unset) and runs the matching CMake preset.

Build output lands under the preset’s binary directory (typically under `build/`).

For clangd / IDEs, point the compilation database at `build/debug/compile_commands.json` (or symlink it to the project root).

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
- Prefer inspecting generated CMake — Forge does not hide the build system.
- Device support is still early; see [ROADMAP.md](../ROADMAP.md).

## See also

- [Getting started](getting-started.md)
- [Commands](commands.md)
