# Getting started

This guide walks you from a clean checkout to a built Forge project.

## Requirements

- **Go** 1.26 or later
- **CMake** 3.20 or later
- A C toolchain on `PATH` (host `gcc`, or `gcc-arm-none-eabi` for STM32 work)

## Build and install Forge

From the repository root:

```bash
go build -o forge .
sudo install -m 755 forge /usr/local/bin/forge
forge version
```

You should see `Forge version 0.2.0` (or newer).

If you prefer not to install globally, run commands with `go run .` from the repo root, or use `./forge` after `go build`.

## First project (STM32-oriented)

```bash
forge new blink --device stm32f103r8
cd blink
forge init
forge build
```

What happens:

1. `new` creates `blink/` and writes `Forge.toml`
2. `init` scaffolds sources and CMake files from that config
3. `build` configures and builds using CMake presets

## First project (host / generic)

```bash
forge new demo_app
cd demo_app
forge init
forge build
```

## Helper script

From the Forge repo root, `run.sh` builds Forge and runs the full create/init/build flow under `examples/`:

```bash
./run.sh blink --device stm32f103r8
```

```bash
./run.sh --help
```

## Next steps

- [Create a project](create-a-project.md) — layout and workflow details
- [Commands](commands.md) — full CLI reference
- [ROADMAP.md](../ROADMAP.md) — where Forge is headed
