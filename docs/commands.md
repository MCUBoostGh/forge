# Commands

CLI reference for Forge **v0.2.0**.

```text
forge <command> [arguments]
```

| Command | Status | Description |
|---------|--------|-------------|
| `new` | Available | Create a project directory and `Forge.toml` (`<device>` or `--device`) |
| `init` | Available | Scaffold files, fetch `dependencies` into the user cache, generate `cmake/Package.cmake` |
| `build` | Available | Configure and build a CMake preset (`forge build debug`) |
| `help` | Available | Show usage |
| `version` | Available | Print the version string |
| `list` | Planned | List supported devices |
| `run` | Planned | Run the project |
| `test` | Planned | Run tests |

## `forge new`

```bash
forge new <name> <device>
forge new <name> --device <device>
```

Creates `<name>/` and writes `Forge.toml`. A device is required (catalog id or alias).

| Argument / flag | Description |
|-----------------|-------------|
| `<name>` | Project directory name (required) |
| `<device>` | Catalog id or alias (for example `stm32f103r8`, `bluepill`, `stm32f746zg`, `nucleo-g431rb`) |
| `--device <device>` | Same as positional `<device>` |

Unknown devices fail with `unknown device "..."`. See [Device catalog](device-catalog.md).

`Forge.toml` `dependencies` come from the device catalog family: F1 uses `cmsis5@5.9.0`, `stm32f1-cmsis-device@4.3.5`, `stm32f1-hal@1.1.10`; F7 and G4 use their matching CMSIS-Device and HAL packs.

## `forge init`

```bash
forge init
```

Must be run in a directory that already contains `Forge.toml`. Scaffolds folders, `main.c`, CMake files, `LinkerScript.ld`, and `cmake/Package.cmake`. Downloads each `dependencies` entry (`name@version`) into `~/.cache/forge/packages/` when that version is not already extracted. Copies `include/*hal_conf.h` from the HAL template when missing. Copies device `system_*.c` and the matching GCC `startup_*.s` into the project root.

## `forge build`

```bash
forge build <preset>
```

Runs `cmake --preset <preset>` then `cmake --build --preset <preset>`. Typical presets: `debug`, `release`. Output is `build/debug` or `build/release`. Build type is not stored in `Forge.toml`.

Requires CMake 3.20+.

## `forge help`

```bash
forge help
```

Prints usage text.

## `forge version`

```bash
forge version
```

Prints the Forge version (currently `0.2.0`).

## Planned commands

These appear in help or the roadmap but are not implemented yet: `run`, `test`, `list`, and later `setup`, `sync`, `flash`, `monitor`, and others. See [ROADMAP.md](../ROADMAP.md).

## See also

- [Getting started](getting-started.md)
- [Create a project](create-a-project.md)
- [Device catalog](device-catalog.md)
