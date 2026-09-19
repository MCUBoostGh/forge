# Create a project

Forge uses a two-step setup, then a CMake-backed build.

```text
forge new  →  Forge.toml
forge init →  package cache + sources + CMake
forge build → ELF / HEX / BIN
```

## 1. Create the project

```bash
forge new <name> <device>
forge new <name> --device <device>
```

Examples:

```bash
forge new blink stm32f103r8
forge new blink --device stm32f103r8
forge new blink --device bluepill
forge new blink stm32f746zg
forge new blink stm32g431rb
```

A device is required. `Resolve` accepts the catalog id (`stm32f103r8`, `stm32f746zg`, `stm32g431rb`, …) or an alias (`bluepill`, `nucleo-f746zg`). Unknown names fail with a clear error. Catalog keys and how to add a board: [Device catalog](device-catalog.md).

`new` creates `<name>/` and writes `Forge.toml` with project metadata and **family** default dependencies from the catalog. STM32F1 example:

```toml
dependencies = [
  "cmsis5@5.9.0",
  "stm32f1-cmsis-device@4.3.5",
  "stm32f1-hal@1.1.10",
]
```

STM32F7 uses `stm32f7-cmsis-device@1.2.10` and `stm32f7-hal@1.3.3`. STM32G4 uses `stm32g4-cmsis-device@1.2.6` and `stm32g4-hal@1.2.6`.

Change a spec before `init` if you want another catalog entry (for example `cmsis6@6.1.0`, or another `stm32*-hal` family from the HAL YAML). Version is taken from this list; URLs live in Forge’s package YAML, not in the project.

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
├── system_stm32f1xx.c
├── startup_stm32f103xb.s
├── CMakeLists.txt
└── CMakePresets.json
```

(F7/G4 projects get `system_stm32f7xx.c` / `system_stm32g4xx.c` and the matching `startup_*.s` instead.)

`main.c` is a Cortex-M CMSIS smoke test (not Hello World). `cmake/Package.cmake` sets `cache_dir` to the shared package cache and references headers/HAL sources as `${cache_dir}/...`. CMSIS Core and CMSIS-Device are INTERFACE libraries; the family HAL is STATIC. `system_*.c` and the device GCC `startup_*.s` are copied into the project root and added to the firmware executable with `main.c`. `include/*hal_conf.h` is copied from the HAL template if it is not already present.

HAL/CMSIS archives stay in the cache. Shared layout:

```text
~/.cache/forge/packages/CMSIS_5-5.9.0/CMSIS/Core/Include
~/.cache/forge/packages/cmsis-device-f1-4.3.5/Include
~/.cache/forge/packages/stm32f1xx-hal-driver-1.1.10/Inc
```

A later project that uses `cmsis5@5.9.0` reuses that directory.

## 3. Build

```bash
forge build debug
```

`forge build <preset>` runs the CMake preset of that name (`debug`, `release`, …). `CMAKE_BUILD_TYPE` and catalog CPU/float ABI/`STM32_DEVICE`/`FLASH_KB`/`RAM_KB` come from `CMakePresets.json`, not from `Forge.toml`. The linker script already uses flash/RAM from the device catalog.

Build output lands under `build/<preset>/` (for example `build/debug/`).

For clangd / IDEs, point the compilation database at `build/debug/compile_commands.json` (or symlink it to the project root).

## Without installing Forge

From a Forge checkout:

```bash
go run . new my_app --device stm32f103r8
cd my_app
go run ../. init
go run ../. build debug
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
- [Device catalog](device-catalog.md)
- [Commands](commands.md)
