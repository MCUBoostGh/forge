# Getting started

This guide walks you from a clean checkout to a built Forge project on Windows or Linux.

## Requirements

- **Go** 1.26 or later
- Network access the first time toolchain binaries, CMSIS, or STM32 HAL are downloaded
- A system **`gcc`** only for a host project (`forge new <name>` with no device). That compiler is not downloaded.

`forge toolchain install` downloads CMake, Ninja, `gcc-arm-none-eabi`, and OpenOCD. You do not need those tools on `PATH` before the first STM32 build. macOS is out of scope.

## Build and install Forge

From the repository root:

```bash
go install .
forge version
```

On Linux you can instead build a binary and copy it onto `PATH`:

```bash
go build -o forge .
sudo install -m 755 forge /usr/local/bin/forge
```

You should see `Forge version 0.2.0`.

If you prefer not to install, run commands with `go run .` from the repo root.

## Install the toolchain

```bash
forge toolchain install all
forge toolchain check
```

`install` reads [internal/toolchain/tools.yaml](../internal/toolchain/tools.yaml) and extracts each tool under the user cache. It does not read `Forge.toml` and does not change your `PATH`. `check` should report `ok` for `cmake`, `ninja`, `gcc-arm-none-eabi`, and `openocd`.

Cache root (`os.UserCacheDir()`):

- Windows: `%LOCALAPPDATA%\forge\toolchains\<name>\<version>\`
- Linux: `~/.cache/forge/toolchains/<name>/<version>/`

`--dry-run` prints the URL and destination without downloading. CMSIS and HAL packages use a different directory, `forge/packages/`, under the same user cache.

## First project (STM32)

```bash
forge new blink stm32f103r8
cd blink
forge init
forge build debug
```

What happens:

1. `new` creates `blink/` and writes `Forge.toml` (CMSIS Core plus the family’s CMSIS-Device and HAL, for example STM32F1)
2. `init` downloads those packages into the package cache if needed, then scaffolds sources and CMake files. Presets use the Ninja generator
3. `build debug` starts the cached `cmake` and prepends the toolchain `bin` directories to that process only

`forge new` accepts a catalog id or alias (`stm32f103r8`, `bluepill`, `stm32f746zg`, `nucleo-g431rb`). Omit the device for a host project.

## Host project

This needs a `gcc` already on `PATH`.

```bash
forge new app
cd app
forge init
forge build host-debug
```

`init` writes `host-debug` and `host-release`. It does not download CMSIS or HAL.

## Tests from a checkout

Unit tests:

```bash
go test ./...
```

End-to-end `new` → `init` → `build debug` (STM32F1 and STM32G4). The integration test still requires `cmake`, `gcc-arm-none-eabi`, and `make` on `PATH`, plus network on the first package download:

```bash
go test -tags=integration ./cmd -count=1 -timeout 20m
```

GitHub Actions runs these on every push and pull request (`.github/workflows/ci.yml`).

## Next steps

- [Create a project](create-a-project.md) — layout, `Forge.toml` dependencies, and package cache
- [Device catalog](device-catalog.md) — YAML keys and how to add a board
- [Commands](commands.md) — full CLI reference
- [ROADMAP.md](../ROADMAP.md) — where Forge is headed
