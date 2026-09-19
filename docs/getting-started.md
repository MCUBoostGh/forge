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
forge new blink stm32f103r8
cd blink
forge init
forge build debug
```

What happens:

1. `new` creates `blink/` and writes `Forge.toml` (CMSIS Core plus the family’s CMSIS-Device and HAL, for example STM32F1)
2. `init` downloads those packages into `~/.cache/forge/packages/` if needed, then scaffolds sources and CMake files
3. `build debug` configures and builds the `debug` CMake preset

`forge new` requires a catalog id or alias (`stm32f103r8`, `bluepill`, `stm32f746zg`, `nucleo-g431rb`). A host-only `new` path is not implemented.

## Tests from a checkout

Unit tests:

```bash
go test ./...
```

End-to-end `new` → `init` → `build debug` (STM32F1 and STM32G4). Requires CMake, `gcc-arm-none-eabi`, `make`, and network on the first package download:

```bash
go test -tags=integration ./cmd -count=1 -timeout 20m
```

## Next steps

- [Create a project](create-a-project.md) — layout, `Forge.toml` dependencies, and package cache
- [Device catalog](device-catalog.md) — YAML keys and how to add a board
- [Commands](commands.md) — full CLI reference
- [ROADMAP.md](../ROADMAP.md) — where Forge is headed
