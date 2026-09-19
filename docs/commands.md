# Commands

CLI reference for Forge **v0.2.0**.

```text
forge <command> [arguments]
```

| Command | Status | Description |
|---------|--------|-------------|
| `new` | Available | Create a project directory and `Forge.toml` (`--device` required) |
| `init` | Available | Scaffold files, fetch `dependencies` into the user cache, generate `cmake/Package.cmake` |
| `build` | Available | Configure and build with CMake presets |
| `help` | Available | Show usage |
| `version` | Available | Print the version string |
| `list` | Planned | List supported devices |
| `run` | Planned | Run the project |
| `test` | Planned | Run tests |

## `forge new`

```bash
forge new <name> --device <device>
```

Creates `<name>/` and writes `Forge.toml`. `--device` is required.

| Argument / flag | Description |
|-----------------|-------------|
| `<name>` | Project directory name (required) |
| `--device <device>` | Device catalog id (for example `stm32f103r8`) |

Default `Forge.toml` includes `dependencies = ["cmsis5@5.9.0"]`.

## `forge init`

```bash
forge init
```

Must be run in a directory that already contains `Forge.toml`. Scaffolds folders, `main.c`, CMake files, `LinkerScript.ld`, and `cmake/Package.cmake`. Downloads each `dependencies` entry (`name@version`) into `~/.cache/forge/packages/` when that version is not already extracted.

## `forge build`

```bash
forge build
```

Loads `Forge.toml`, selects a CMake preset from the configured build type (falls back to `debug`), then configures and builds.

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
