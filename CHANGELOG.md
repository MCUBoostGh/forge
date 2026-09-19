# Changelog

All notable user-visible changes to Forge are listed here.

## [0.2.0] — in progress

### Added

- Default `Forge.toml` dependencies `cmsis5@5.9.0`, `stm32f1-cmsis-device@4.3.5`, and `stm32f1-hal@1.1.10` on `forge new`.
- `forge init` downloads CMSIS and STM32F1 HAL once into `~/.cache/forge/packages/` and generates `cmake/Package.cmake` with `set(cache_dir ...)` (INTERFACE for CMSIS Core and CMSIS-Device headers, STATIC for HAL).
- Copy device `system_stm32f1xx.c` and the matching GCC `startup_*.s` into the project root and add them to the firmware executable.
- Generated `main.c` is a Cortex-M CMSIS smoke test (CMSIS version + `SCB->CPUID`), not Hello World.
- Copy `stm32f1xx_hal_conf.h` into `include/` from the HAL template when missing.

### Changed

- CMSIS Core, CMSIS-Device headers, and HAL stay in the shared cache; CMake uses `${cache_dir}/...`.
- Device `system_*.c` and GCC `startup_*.s` are copied into the project root (not compiled as a `cmsis-device` CMake library).
- `forge init` reads `dependencies` from `Forge.toml` (`name@version`) before scaffolding.

### Notes

- `forge new` still requires `--device` (for example `stm32f103r8`).
