# Changelog

All notable user-visible changes to Forge are listed here.

## [0.2.0] — in progress

### Added

- Default `Forge.toml` dependency `cmsis5@5.9.0` on `forge new`.
- `forge init` downloads CMSIS once into `~/.cache/forge/packages/` and generates `cmake/Package.cmake` as a CMake INTERFACE library.
- Generated `main.c` is a Cortex-M CMSIS smoke test (CMSIS version + `SCB->CPUID`), not Hello World.

### Changed

- Third-party sources are not copied into the project tree; CMake include paths point at the shared cache.
- `forge init` reads `dependencies` from `Forge.toml` (`name@version`) before scaffolding.

### Notes

- `forge new` still requires `--device` (for example `stm32f103r8`).
- Linking currently warns about a missing `Reset_Handler` until a startup file is generated.
