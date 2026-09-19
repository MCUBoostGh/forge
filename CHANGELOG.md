# Changelog

All notable user-visible changes to Forge are listed here.

## [0.2.0] — 2026-09-19

### Added

- Default `Forge.toml` dependencies `cmsis5@5.9.0`, `stm32f1-cmsis-device@4.3.5`, and `stm32f1-hal@1.1.10` on `forge new`.
- `forge init` downloads CMSIS and STM32F1 HAL once into `~/.cache/forge/packages/` and generates `cmake/Package.cmake` with `set(cache_dir ...)` (INTERFACE for CMSIS Core and CMSIS-Device headers, STATIC for HAL).
- Copy device `system_stm32f1xx.c` and the matching GCC `startup_*.s` into the project root and add them to the firmware executable.
- Generated `main.c` is a Cortex-M CMSIS smoke test (CMSIS version + `SCB->CPUID`), not Hello World.
- Copy `stm32f1xx_hal_conf.h` into `include/` from the HAL template when missing.
- Document device catalog YAML keys and how to add a board ([docs/device-catalog.md](docs/device-catalog.md)).
- `CMakePresets.json` entries `debug` / `release` (plus relwithdebinfo/minsizerel), with catalog CPU, float ABI, FPU, `STM32_DEVICE`, `FLASH_KB`, and `RAM_KB`. Output directories are `build/debug` and `build/release`.
- `forge build <preset>` runs that CMake preset (`forge build debug`). Build type is not stored in `Forge.toml`.
- `forge new` writes family CMSIS/HAL `dependencies` from the device catalog (F1, F7, G4).
- STM32F7 (`stm32f746zg`, `stm32f767zi`) and STM32G4 (`stm32g431rb`, `stm32g474re`) board profiles, CMSIS-Device packs, and HAL.
- Integration tests: `go test -tags=integration ./cmd` runs `new` → `init` → `build debug` for STM32F1 and STM32G4.

### Changed

- CMSIS Core, CMSIS-Device headers, and HAL stay in the shared cache; CMake uses `${cache_dir}/...`.
- Device `system_*.c` and GCC `startup_*.s` are copied into the project root (not compiled as a `cmsis-device` CMake library).
- `forge init` reads `dependencies` from `Forge.toml` (`name@version`) before scaffolding.
- Remove `build.type` from `Forge.toml`. `CMAKE_BUILD_TYPE` lives only in CMake presets.

### Notes

- `forge new` requires a catalog id or alias (for example `stm32f103r8`, `bluepill`, `stm32f746zg`, `nucleo-g431rb`).
