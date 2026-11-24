# Create a project

Forge uses a two-step setup, then a CMake-backed build.

```text
forge toolchain install all  →  toolchain cache
forge new                    →  Forge.toml
forge init                   →  sources + CMake
forge build                  →  ELF / HEX / BIN, or a host executable
```

`forge toolchain` is optional if CMake, Ninja, and `gcc-arm-none-eabi` are already installed. It does not read `Forge.toml`. See [Commands](commands.md) for cache paths and `check`.

## 1. Create the project

```bash
forge new <name>
forge new <name> <device>
forge new <name> --device <device>
```

Examples:

```bash
forge new app
forge new blink stm32f103r8
forge new blink --device stm32f103r8
forge new blink --device bluepill
forge new blink stm32f746zg
forge new blink stm32g431rb
forge new examples/blink stm32f103r8
```

`forge new <name>` with no device writes a host project (`target.kind` is `host`, compiler `gcc`, empty `dependencies`). A device is a catalog id (`stm32f103r8`, `stm32f746zg`, `stm32g431rb`, …) or an alias (`bluepill`, `nucleo-f746zg`). Unknown names fail with a clear error. Catalog keys and how to add a board: [Device catalog](device-catalog.md).

`new` creates `<name>/` and writes `Forge.toml`. The project name is the last path element, so `examples/blink` creates `examples/blink/` and names the project `blink`.

STM32 `new` writes **family** default dependencies from the catalog. STM32F1 example:

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

`init` reads `Forge.toml` from the current directory.

A host project looks like this:

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

`main.c` is a small `puts("hello")` program. Presets are `host-debug` (Debug) and `host-release` (Release), both with the Ninja generator. There is no linker script, `cmake/Package.cmake`, or `cmake/gcc-arm-none-eabi.cmake`. The compiler is the system `gcc`.

An STM32 project also downloads each `name@version` into the package cache if it is missing, and generates:

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

STM32 `main.c` is a Cortex-M CMSIS smoke test (not Hello World). Presets are `debug`, `release`, `relwithdebinfo`, and `minsizerel`, generated with Ninja. `cmake_minimum_required` is `3.20`. A tree created when the generator was `Unix Makefiles` keeps that generator until you run `forge init` again.

`cmake/Package.cmake` sets `cache_dir` to the shared package cache and references headers/HAL sources as `${cache_dir}/...`. CMSIS Core and CMSIS-Device are INTERFACE libraries; the family HAL is STATIC. `system_*.c` and the device GCC `startup_*.s` are copied into the project root and added to the firmware executable with `main.c`. `include/*hal_conf.h` is copied from the HAL template if it is not already present.

HAL/CMSIS archives stay in the package cache, not the toolchain cache. Shared layout:

```text
~/.cache/forge/packages/CMSIS_5-5.9.0/CMSIS/Core/Include
~/.cache/forge/packages/cmsis-device-f1-4.3.5/Include
~/.cache/forge/packages/stm32f1xx-hal-driver-1.1.10/Inc
```

On Windows the same tree is under `%LOCALAPPDATA%\forge\packages\`. A later project that uses `cmsis5@5.9.0` reuses that directory.

## 3. Build

```bash
forge build debug
forge build host-debug
```

`forge build <preset>` runs the CMake preset of that name. Forge launches the cached `cmake` executable and prepends installed toolchain `bin` directories to that process, so Ninja and `arm-none-eabi-gcc` are visible without editing your `PATH`. `CMAKE_BUILD_TYPE` and catalog CPU/float ABI/`STM32_DEVICE`/`FLASH_KB`/`RAM_KB` come from `CMakePresets.json`, not from `Forge.toml`. The linker script already uses flash/RAM from the device catalog.

Build output lands under `build/<preset>/` (for example `build/debug/` or `build/host-debug/`).

For clangd / IDEs, point the compilation database at `build/debug/compile_commands.json` (or symlink it to the project root).

## Without installing Forge

From a Forge checkout:

```bash
go run . toolchain install all
go run . new my_app --device stm32f103r8
cd my_app
go run ../. init
go run ../. build debug
```

To smoke-test the same flow from the repo (F1 and G4), the integration test still expects `cmake`, `gcc-arm-none-eabi`, and `make` on `PATH`:

```bash
go test -tags=integration ./cmd -count=1 -timeout 20m
```

## Tips

- Run `init` and `build` from the project directory (where `Forge.toml` lives).
- Prefer inspecting generated CMake — Forge does not hide the build system.
- Device support is still early; see [ROADMAP.md](../ROADMAP.md).

## See also

- [Getting started](getting-started.md)
- [Device catalog](device-catalog.md)
- [Commands](commands.md)
