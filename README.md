# Forge

**A unified development workflow for modern embedded C/C++.**

## What is Forge?

Forge is a lightweight CLI for embedded C/C++ projects. It ties project creation, configuration, and CMake-based builds into one consistent workflow so firmware work starts from a clear, repeatable baseline instead of ad-hoc copy-paste setups.

At the center is `Forge.toml`: the project’s source of truth. Forge commands read it to scaffold sources, generate CMake files, and drive builds.

## Why does it exist?

Embedded projects often begin with fragile CMake snippets, unclear toolchain paths, and board-specific one-offs that are hard to reproduce. Forge exists to make that setup a first-class workflow:

- **One config** — project and target settings live in `Forge.toml`
- **One scaffold** — generate a clean, readable CMake layout from that config
- **One build path** — configure and build through CMake presets
- **Device-aware** — start from a named device (for example STM32) instead of a blank folder

The goal is not to hide CMake or replace vendor ecosystems. It is to give embedded engineers a unified path from an empty directory to a buildable project — and eventually through flash and monitor.

## What does it currently do?

Forge is early-stage (**v0.2.0**). Today it supports:

| Capability | Details |
|------------|---------|
| Project creation | `forge new <name> <device>` (or `--device`) writes a project dir and `Forge.toml` |
| Scaffolding | `forge init` generates CMake files, a Cortex-M smoke-test `main.c`, device `startup_*.s` / `system_*.c`, and `cmake/Package.cmake` |
| Builds | `forge build <preset>` runs a CMake preset (`debug` or `release`); output under `build/debug` or `build/release` |
| CMSIS | `cmsis5@5.9.0` is the default Core dependency; unpacked once under `~/.cache/forge/packages/` |
| STM32 HAL | Family HAL from the catalog (`stm32f1-hal`, `stm32f7-hal`, or `stm32g4-hal`) is a cached STATIC library; device headers stay INTERFACE |
| STM32 target | Bare-metal F1, F7, and G4 (`gcc-arm-none-eabi`); host Linux is not a separate `new` flow yet |
| Device catalog | STM32 YAML profiles under `internal/devices/` ([how to add a board](docs/device-catalog.md)) |

**Not ready yet:** `run`, `test`, `list`, `setup`, `flash`, `monitor`, and related roadmap commands.

For command details, see [docs/commands.md](docs/commands.md).

## Where is it going?

**v0.2.0** is the current STM32 bootstrap (`new` → `init` → `build debug`). Next is **v0.3.0**: host tools (`forge setup` / `sync`) and a host Linux path.

Toward **v1.0.0**:

- Stable host Linux + STM32 bare-metal support
- `setup`, `sync`, `flash`, `monitor`, and related workflow commands
- CMake library targets and Docker-based reproducible environments

**v2.0.0** looks at embedded Linux targets (for example Raspberry Pi).

Full plan: [ROADMAP.md](ROADMAP.md). Task breakdown: [TODO.md](TODO.md).

## How do I try it?

Requires **Go 1.26+**, **CMake 3.20+**, **`gcc-arm-none-eabi`**, and network access the first time CMSIS is fetched.

```bash
go build -o forge .
sudo install -m 755 forge /usr/local/bin/forge

forge new blink stm32f103r8
cd blink
forge init
forge build debug
```

### Tutorials and docs

| Guide | Description |
|-------|-------------|
| [Getting started](docs/getting-started.md) | Install Forge and run your first project |
| [Create a project](docs/create-a-project.md) | Walkthrough of `new` → `init` → `build` |
| [Commands](docs/commands.md) | CLI reference |
| [Device catalog](docs/device-catalog.md) | YAML keys and how to add a board |
| [Changelog](CHANGELOG.md) | User-visible changes |
| [Docs index](docs/README.md) | All documentation |

## License

See [LICENSE](LICENSE).
