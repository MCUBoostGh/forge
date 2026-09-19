# Forge TODO

Work breakdown derived from [README.md](README.md) and [ROADMAP.md](ROADMAP.md).

**Convention:** each version is one **GitHub Issue**. Nested items are tasks and subtasks for that issue. Check boxes as work lands; close the issue when every task for that version is done.

| Status | Meaning |
|--------|---------|
| `[x]` | Done in tree / shipped |
| `[ ]` | Not done |
| **In progress** | Partially implemented on branch (e.g. `V0.2.0`) |

---

## Issue: v0.1.0 — Prototype CLI skeleton

> Basic CLI and generic project scaffolding (README current commands: `new`, `init`, `build`, `help`, `version`).

### Tasks

- [x] **T1 — CLI entry and dispatch**
  - [x] `main.go` command switch (`new`, `init`, `build`, `help`, `version`)
  - [x] `forge help` usage text
  - [x] `forge version` string
  - [ ] Remove or clearly mark stub help entries for unimplemented `run` / `test` (README notes they appear in help but are not implemented)

- [x] **T2 — `forge new <project_name>`**
  - [x] Create project directory
  - [x] Write generic `Forge.toml` (now via `internal/config.New`; legacy `contents/` removed)
  - [x] Project metadata defaults (name, version, toolchain, CMake settings)

- [x] **T3 — `forge init`**
  - [x] Generate `src/`, `include/`
  - [x] Generate Hello World `main.c` (replaced in v0.2.0 by a Cortex-M CMSIS smoke test)
  - [x] Generate minimal `CMakeLists.txt`
  - [x] Generate `CMakePresets.json`
  - [x] Fix known limitation: init must read `Forge.toml` from disk (completed under v0.2.0)

- [x] **T4 — `forge build`**
  - [x] `cmake -S . -B build`
  - [x] `cmake --build build`
  - [x] Document CMake 3.20+ requirement

- [x] **T5 — Docs for prototype**
  - [x] README overview, commands table, usage, example
  - [x] ROADMAP v0.1.0 section

**Done when:** Issue closed; README prototype workflow works end-to-end.

---

## Issue: v0.2.0 — STM32 project bootstrap

> Phase 1 start. Device-aware creation for STM32 (first board: `stm32f103r8`). **In progress** — catalog anchors, `--device` new, `internal/config`, arm-none-eabi + `LinkerScript.ld`, CMSIS cache + `cmake/Package.cmake` exist; startup, STM32 presets, host path, and install seeding still open.

### Tasks

- [ ] **T1 — Device catalog**
  - [x] `internal/devices/` package skeleton
  - [x] Embed / load ST YAML catalog
  - [x] `Lookup` / `Resolve` / `List` APIs
  - [x] First complete board profile: `stm32f103r8` (CPU, flash/RAM, linker, OpenOCD target, aliases) — plus YAML anchors + `stm32f103c8`
  - [ ] Document catalog keys and how to add a board

- [ ] **T2 — Device-aware `forge new <name> <device>`**
  - [x] Accept device argument / `--device` path
  - [ ] Resolve device via catalog; fail clearly on unknown device (`Lookup` only today; aliases / `Resolve` not wired; positional `<device>` not supported)
  - [ ] Write device-aware `Forge.toml` (`[target]`, board, toolchain hints) — `config.New` still resets via `setConfigDefaults()` and drops device fields
  - [ ] Seed `[install].packages` from target kind when applicable
  - [x] Usage/examples: docs/README use `forge new blink --device stm32f103r8` (flag form; positional form still open)

- [ ] **T3 — `forge init` reads `Forge.toml`**
  - [x] Load and validate `Forge.toml` from cwd (replace in-memory-only path)
  - [ ] Generate STM32-aware tree (startup/linker refs as planned) — `LinkerScript.ld` emitted; startup `.s` not yet
  - [ ] Board-specific CMake generation from config (CPU/FloatABI + flash/RAM templating; CMSIS INTERFACE via `cmake/Package.cmake`; STM32 presets / startup not done)
  - [ ] Keep host/generic path working for non-STM32 projects (`syncTemplateData` requires MCU/STM32; presets still arm-toolchain on `host-base`)

- [ ] **T4 — gcc-arm-none-eabi + CMake presets**
  - [x] Toolchain template content (`internal/templates/cmake/gcc-arm-none-eabi.cmake.tmpl`)
  - [x] Emit CMake toolchain file for arm-none-eabi (`cmake/gcc-arm-none-eabi.cmake` on `init`)
  - [ ] `CMakePresets.json` entries: `stm32-debug` / `stm32-release` (still host-named `debug` / `release`)
  - [ ] Wire presets to board catalog fields (CPU/FloatABI in toolchain + flash/RAM in linker; catalog `presets:` unused)

- [ ] **T5 — Config foundation**
  - [x] Harden `Forge.toml` load/save toward `internal/config/` (`New` / `Read` / `Write` / `Get` / `Set` + tests; little validation yet)
  - [ ] Align schema fields with ROADMAP (`target.kind`, `board`, toolchain, `[install]` / `[tools]`)

- [ ] **T6 — Docs**
  - [x] Update README for device-aware `new` / STM32 init (getting-started + create-a-project docs)
  - [x] Document CMSIS cache, `cmake/Package.cmake`, and default `dependencies`
  - [ ] Mark v0.2.0 delivered in ROADMAP when complete

**Done when:** `forge new blink stm32f103r8` → `init` → CMake STM32 presets work; issue closed.

---

## Issue: v0.3.0 — Host tools, setup, and multi-target build

> Install/detect toolchains; host Linux presets; smarter `build`.

### Tasks

- [ ] **T1 — `internal/setup/` + `internal/install/`**
  - [ ] Map target kind (`host` / `stm32`) to apt package lists
  - [ ] Package manifest per target kind
  - [ ] Detect already-installed tools vs missing

- [ ] **T2 — `forge setup`**
  - [ ] Read `Forge.toml` (`[target]`, `[install]`, `[tools]`)
  - [ ] Install missing packages via apt (Ubuntu/Debian only for v1)
  - [ ] Flags: `--dry-run`, `--yes`
  - [ ] On success, run `forge sync`

- [ ] **T3 — `forge sync`**
  - [ ] Detect installed tool paths/versions
  - [ ] Write paths into `Forge.toml` (`[toolchain]`, `[tools]`)
  - [ ] Clear errors when tools missing

- [ ] **T4 — Host Linux + preset-driven `forge build`**
  - [ ] Host preset (`host-debug`)
  - [ ] Select preset from `Forge.toml`
  - [ ] Enhance `forge build` beyond bare configure/build

- [ ] **T5 — Docs**
  - [ ] Document setup/sync workflow and apt-only host constraint

**Done when:** setup → sync → build works for host and STM32 toolchains; issue closed.

---

## Issue: v0.4.0 — STM32 on-device workflow

> Flash and serial monitor driven by `Forge.toml` and board defaults.

### Tasks

- [ ] **T1 — Flash defaults in device catalog**
  - [ ] ST-Link + SWD defaults for `stm32f103r8`
  - [ ] `[flash]` schema: interface, transport

- [ ] **T2 — `forge flash`**
  - [ ] Locate `.elf` / `.bin` from build output
  - [ ] Flash via OpenOCD (ST-Link)
  - [ ] Fallback to `st-flash`
  - [ ] Read flash settings from `Forge.toml`

- [ ] **T3 — `forge monitor`**
  - [ ] Attach to serial port; stream logs
  - [ ] Port/baud from `[monitor]` in `Forge.toml`
  - [ ] Sensible USART defaults per board in catalog

- [ ] **T4 — Docs**
  - [ ] Document flash/monitor hardware prerequisites

**Done when:** build → flash → monitor works on `stm32f103r8`; issue closed.

---

## Issue: v0.5.0 — Libraries and third-party code

> CMake library targets and third-party integration.

### Tasks

- [ ] **T1 — `forge lib <name>`**
  - [ ] Scaffold `lib/<name>/` with `CMakeLists.txt`, `include/`, `src/`
  - [ ] Support static (and shared if in scope) library layout

- [ ] **T2 — Wire libraries into root build**
  - [ ] Update root `CMakeLists.txt` to consume `forge lib` output
  - [ ] CMake library targets as first-class build graph nodes

- [ ] **T3 — `forge add <package>`**
  - [x] Integrate CMSIS Core on `forge init` from `Forge.toml` `dependencies` (cache + INTERFACE library; not a `forge add` command yet)
  - [ ] Integrate STM32 HAL, FreeRTOS, and a dedicated `forge add` command
  - [ ] Update root `CMakeLists.txt` for additional packages
  - [x] Record CMSIS under `dependencies` in new `Forge.toml` (`cmsis5@5.9.0`)

- [ ] **T4 — Docs**
  - [ ] Examples for lib + add workflows

**Done when:** lib/add produce a linkable multi-target CMake project; issue closed.

---

## Issue: v0.6.0 — Reproducible builds and testing

> Docker env generation and host/unit test scaffolding.

### Tasks

- [ ] **T1 — `forge docker`**
  - [ ] Generate `Dockerfile` pinned to tool versions from `Forge.toml`
  - [ ] Optional `docker-compose.yml`
  - [ ] Document how to build inside the container

- [ ] **T2 — `forge test`**
  - [ ] Create `tests/` layout
  - [ ] Fetch Unity + CMock via CMake `FetchContent`
  - [ ] Implement the `test` command (replace README stub)
  - [ ] Decide fate of undocumented `run` stub in help

- [ ] **T3 — Complete CMake presets matrix**
  - [ ] Unified presets: host/stm32 × debug/release
  - [ ] Align preset names with `Forge.toml` and docs

- [ ] **T4 — Docs**
  - [ ] Reproducible-build and test sections in README

**Done when:** docker + test + full preset matrix work; issue closed.

---

## Issue: v1.0.0 — First stable release

> Stable Linux host + STM32 release with full v1 command set and polish.

### Tasks

- [ ] **T1 — `forge doc`**
  - [ ] Generate Doxygen config from `Forge.toml` / sources
  - [ ] Scaffold `docs/`

- [ ] **T2 — `forge package`**
  - [ ] Produce host distributable (`.tar.gz` and/or `.deb`)
  - [ ] Include requirements manifest

- [ ] **T3 — v1.0.0 checklist (from ROADMAP)**
  - [ ] Toolchains: `gcc`, `gcc-arm-none-eabi`
  - [ ] Infrastructure: Docker, CMake libs, CMake presets, `forge setup`
  - [ ] Targets: host Linux, STM32 on Linux host
  - [ ] Commands: `new`, `init`, `setup`, `sync`, `docker`, `build`, `doc`, `flash`, `monitor`, `test`, `lib`, `add`, `package`

- [ ] **T4 — Release polish**
  - [ ] Reference example: `blink` on STM32F103R8
  - [ ] README rewritten for v1 workflow
  - [ ] Integration smoke tests for all v1 commands
  - [ ] Tag release `v1.0.0`

**Done when:** checklist complete, example works, release tagged; issue closed.

---

## Issue: v2.0.0 — Embedded Linux (future)

> Deferred from v1. Raspberry Pi / `gcc-arm-linux-gnueabihf`.

### Tasks

- [ ] **T1 — ARM Linux toolchain module**
  - [ ] CMake toolchain file for `gcc-arm-linux-gnueabihf`
  - [ ] Sync/setup package mapping for arm-linux kind

- [ ] **T2 — Raspberry Pi target preset**
  - [ ] `target.kind = "arm-linux"`
  - [ ] Board profile (e.g. `raspberry-pi-4`)

- [ ] **T3 — Cross-compile `forge build` path**
  - [ ] Build Linux ARM binaries on host using synced toolchain paths
  - [ ] Presets for arm-linux debug/release

- [ ] **T4 — Docs**
  - [ ] Document embedded Linux workflow and scope vs STM32

**Done when:** Pi-style cross-build documented and working; issue closed.

---

## Out of scope (do not file as v1 issues)

Per ROADMAP — track separately only if needed later:

- Windows/macOS host support
- Full STM32 family catalog beyond F103 (expand post-1.0)
- IDE plugins (VS Code extension)
- Cloud CI templates

---

## Suggested GitHub issue titles

Copy-paste when opening issues:

1. `v0.1.0: Prototype CLI skeleton`
2. `v0.2.0: STM32 project bootstrap`
3. `v0.3.0: Host tools, setup, and multi-target build`
4. `v0.4.0: STM32 on-device workflow (flash/monitor)`
5. `v0.5.0: Libraries and third-party code (lib/add)`
6. `v0.6.0: Reproducible builds and testing (docker/test)`
7. `v1.0.0: First stable release`
8. `v2.0.0: Embedded Linux (Raspberry Pi)`

Paste the matching section from this file into each issue body. Optionally add each issue to a GitHub Project board (Todo → In Progress → Done).
