# Commands

CLI reference for Forge **v0.3.0**.

```text
forge <command> [arguments]
```

| Command | Status | Description |
|---------|--------|-------------|
| `new` | Available | Create a project directory and `Forge.toml` (device optional) |
| `init` | Available | Scaffold files. STM32 projects also fetch packages and generate `cmake/Package.cmake` |
| `build` | Available | Configure and build a CMake preset (`forge build debug`) |
| `toolchain` | Available | Install or check catalog tools from `tools.yaml` |
| `help` | Available | Show usage |
| `version` | Available | Print the version string |
| `list` | Planned | List supported devices |
| `run` | Planned | Run the project |
| `test` | Planned | Run tests |

## `forge new`

```bash
forge new <name>
forge new <name> <device>
forge new <name> --device <device>
```

Creates `<name>/` and writes `Forge.toml`. `<name>` may be a nested path such as `examples/blink`; the project name written into `Forge.toml` and CMake is the last path element (`blink`).

| Argument / flag | Description |
|-----------------|-------------|
| `<name>` | Project directory (required) |
| `<device>` | Catalog id or alias (for example `stm32f103r8`, `bluepill`, `stm32f746zg`, `nucleo-g431rb`) |
| `--device <device>` | Same as positional `<device>` |

With no device, the project is a host project: `target.kind` is `host`, the compiler is `gcc`, and `dependencies` is empty. Host builds use the system `gcc`. That compiler is not in the toolchain catalog.

With a device, `Forge.toml` `dependencies` come from the device catalog family: F1 uses `cmsis5@5.9.0`, `stm32f1-cmsis-device@4.3.5`, `stm32f1-hal@1.1.10`; F7 and G4 use their matching CMSIS-Device and HAL packs.

Unknown devices fail with `unknown device "..."`. See [Device catalog](device-catalog.md).

## `forge init`

```bash
forge init
```

Must be run in a directory that already contains `Forge.toml`.

A host project gets `src/`, `include/`, `cmake/`, a plain `main.c`, `CMakeLists.txt`, and `CMakePresets.json` with `host-debug` and `host-release`. It does not download packages and does not write a linker script, `cmake/Package.cmake`, or `cmake/gcc-arm-none-eabi.cmake`.

An STM32 project scaffolds the same folders plus `LinkerScript.ld` and `cmake/Package.cmake`. It downloads each `dependencies` entry (`name@version`) into the package cache when that version is not already extracted, copies `include/*hal_conf.h` from the HAL template when missing, and copies device `system_*.c` and the matching GCC `startup_*.s` into the project root.

Package cache (`os.UserCacheDir()`), separate from the toolchain cache:

- Windows: `%LOCALAPPDATA%\forge\packages\`
- Linux: `~/.cache/forge/packages/`

STM32 `CMakePresets.json` uses the Ninja generator. Presets are `debug`, `release`, `relwithdebinfo`, and `minsizerel`. A project generated before that change still says `Unix Makefiles` until `forge init` is run again. Generated `cmake_minimum_required` stays `3.20`.

## `forge build`

```bash
forge build <preset>
```

Runs `cmake --preset <preset>` then `cmake --build --preset <preset>`. STM32 presets include `debug` and `release`. Host presets are `host-debug` and `host-release`. Output is `build/<preset>/`. Build type is not stored in `Forge.toml`.

Forge starts the cached `cmake` by its absolute path. It also prepends each installed toolchain `bin` directory to that CMake process only, so CMake can find `ninja` and `arm-none-eabi-gcc`. The user `PATH` is not changed. If `cmake` is not in the cache or on `PATH`, the error is `cmake not found; run forge toolchain install cmake`.

Generated projects ask for CMake 3.20 or newer. The CMake binary Forge downloads is the version in `internal/toolchain/tools.yaml`.

## `forge toolchain`

```bash
forge toolchain install <name|all>
forge toolchain install <name|all> --dry-run
forge toolchain check
forge toolchain check <name>
```

Reads the embedded catalog [internal/toolchain/tools.yaml](../internal/toolchain/tools.yaml). It does not read or write `Forge.toml`, does not call apt, and does not change the user `PATH`.

| Name | Version | Binary |
|------|---------|--------|
| `cmake` | 4.5.1 | `bin/cmake` |
| `ninja` | 1.13.1-1 | `bin/ninja` |
| `gcc-arm-none-eabi` | 15.2.1-1.1 | `bin/arm-none-eabi-gcc` |
| `openocd` | 0.12.0-7 | `bin/openocd` |

`install all` installs every key in that file. Host `gcc` is not an entry. URLs, sha256 sums, and archive roots live in the YAML. `{version}` is substituted from the version field. Linux archives are `.tar.gz`; Windows archives are `.zip`. Hosts are Windows and Linux amd64. macOS is out of scope.

Install root (`os.UserCacheDir()`):

- Windows: `%LOCALAPPDATA%\forge\toolchains\<name>\<version>\`
- Linux: `~/.cache/forge/toolchains/<name>/<version>/`

`install` skips a tool whose binary is already in that directory. `--dry-run` prints the name, version, URL, and destination and does not download.

`check` probes the cache first, then `PATH`, and runs `<tool> --version`. A version mismatch in the cache does not fall through to `PATH`. Status is `ok`, `missing`, or `version mismatch`. The command exits non-zero when any reported tool is not `ok`. A version matches when the tool output contains the full catalog version or its leading `x.y.z` (`1.13.1-1` matches `1.13.1`).

## `forge help`

```bash
forge help
```

Prints usage text.

## `forge version`

```bash
forge version
```

Prints the Forge version (currently `0.3.0`).

## Planned commands

These appear in help or the roadmap but are not implemented yet: `run`, `test`, `list`, and later `flash`, `monitor`, and others. See [ROADMAP.md](../ROADMAP.md).

## See also

- [Getting started](getting-started.md)
- [Create a project](create-a-project.md)
- [Device catalog](device-catalog.md)
