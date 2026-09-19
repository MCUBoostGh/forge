# Getting started

This guide walks you from a clean checkout to a built Forge project.

## Requirements

- **Go** 1.26 or later
- **CMake** 3.20 or later
- **`gcc-arm-none-eabi`** on `PATH` (STM32 / Cortex-M)
- Network access the first time CMSIS or STM32 HAL is downloaded

## Build and install Forge

From the repository root:

```bash
go build -o forge .
sudo install -m 755 forge /usr/local/bin/forge
forge version
```

You should see `Forge version 0.2.0` (or newer).

If you prefer not to install globally, run commands with `go run .` from the repo root, or use `./forge` after `go build`.

## First project (STM32)

```bash
forge new blink --device stm32f103r8
cd blink
forge init
forge build
```

What happens:

1. `new` creates `blink/` and writes `Forge.toml` (including CMSIS Core, STM32F1 CMSIS-Device, and `stm32f1-hal`)
2. `init` downloads those packages into `~/.cache/forge/packages/` if needed, then scaffolds sources and CMake files
3. `build` configures and builds using CMake presets

`forge new` currently requires `--device`. A host-only `new` path is not implemented.

## Helper script

From the Forge repo root, `run.sh` builds Forge and runs the full create/init/build flow under `examples/`:

```bash
./run.sh blink --device stm32f103r8
```

```bash
./run.sh --help
```

## Next steps

- [Create a project](create-a-project.md) — layout, `Forge.toml` dependencies, and package cache
- [Commands](commands.md) — full CLI reference
- [ROADMAP.md](../ROADMAP.md) — where Forge is headed
