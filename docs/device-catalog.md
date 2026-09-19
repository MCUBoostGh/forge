# Device catalog

Board profiles live in YAML under `internal/devices/st/` and are embedded into the Forge binary (`//go:embed st/*.yaml`). `forge new` and `forge init` resolve them through `devices.Resolve` (catalog id or alias).

Shipped ids today: `stm32f103r8`, `stm32f103c8`. Alias example: `bluepill` → `stm32f103r8`.

`forge list` is not implemented. Unknown names fail at `forge new` / `init` with `unknown device "..."`.

## Catalog keys

Each selectable board is a YAML map key. That key is the catalog id (`device.ID` is set from the key at load time). Keys that start with `x-` are skipped (anchor-only bases, not valid `--device` values).

| Key | YAML field | Used by `new` / `init` today |
|-----|------------|------------------------------|
| Catalog id | map key | `--device` / positional `<device>`; written to `Forge.toml` as `target.device` |
| `id` | `id` | Optional in YAML; overwritten by the map key |
| `vendor` | `vendor` | Copied into templates (`TargetVendor`) |
| `family` | `family` | Copied into templates (`TargetFamily`) |
| `series` | `series` | Linker-script comment (`TargetSeries`) |
| `architecture` | `architecture` | Stored only |
| `cpu` | `cpu` | `-mcpu` via CMake `CPU` (preset + toolchain fallback); CMSIS core header; toolchain name via `compilersMap` |
| `fpu` | `fpu` | CMake `FPU` on `stm32-base` (`-mfpu` when non-empty) |
| `float_abi` | `float_abi` | `-mfloat-abi` via CMake `FLOAT_ABI` |
| `flash_kb` | `flash_kb` | `FLASH` length in `LinkerScript.ld`; CMake `FLASH_KB` on the STM32 preset |
| `ram_kb` | `ram_kb` | `RAM` length in `LinkerScript.ld`; CMake `RAM_KB` on the STM32 preset |
| `pins` | `pins` | Stored only |
| `package` | `package` | Stored only |
| `stm32_device` | `stm32_device` | `USE_HAL_DRIVER` + device define; GCC startup file `startup_<lowercase>.s`; CMake `STM32_DEVICE` |
| `linker_script` | `linker_script` | Stored only; init always emits `LinkerScript.ld` from the Forge template |
| `openocd_target` | `openocd_target` | Stored only (`forge flash` is not implemented) |
| `presets` | `presets` | `debug` / `release` `build_type` → `CMAKE_BUILD_TYPE` on the `debug` / `release` CMake presets |
| `aliases` | `aliases` | Accepted by `forge new` via `Resolve` (for example `bluepill`); `Forge.toml` stores the canonical id |

`cpu` must be a key in `compilersMap` (`cmd/private.go`). Currently that is `cortex-m0` and `cortex-m3` for `gcc-arm-none-eabi` (plus `x86_64` → `gcc`, unused by STM32 boards). An unknown `cpu` leaves the compiler unset and `Forge.toml` falls back to `gcc`.

`stm32_device` must match a CMSIS-Device GCC startup in the cached pack, for example `STM32F103xB` → `startup_stm32f103xb.s` under `Source/Templates/gcc/`. Default dependencies still pull **STM32F1** CMSIS-Device and HAL (`stm32f1-cmsis-device@4.3.5`, `stm32f1-hal@1.1.10`). Adding a non-F1 part also requires matching packages; that is not automatic.

## File layout and YAML merge

`internal/devices/st/f1.yaml` uses YAML anchors so family fields are not repeated:

```yaml
x-stm32f1: &stm32f1
  vendor: ST
  family: STM32F1
  architecture: cortex-m3
  cpu: cortex-m3
  fpu: ""
  float_abi: soft
  openocd_target: stm32f1x
  presets:
    debug: { build_type: Debug }
    release: { build_type: Release }

x-stm32f103: &stm32f103
  <<: *stm32f1
  series: STM32F103XX

stm32f103r8:
  <<: *stm32f103
  id: stm32f103r8
  flash_kb: 64
  ram_kb: 20
  pins: 64
  package: LQFP64
  stm32_device: STM32F103xB
  linker_script: STM32F103R8Tx_FLASH.ld
  aliases:
    - stm32f103r
    - bluepill
```

New family files must match `st/*.yaml` (for example `st/f4.yaml`) so the embed glob picks them up.

## How to add a board

1. Choose a lowercase catalog id. That string is the positional device / `--device` value.
2. Add an entry in the right `internal/devices/st/*.yaml` file. Reuse `x-` anchors when the MCU is in an existing family.
3. Set at least `cpu`, `float_abi`, `flash_kb`, `ram_kb`, and `stm32_device` so init can emit a toolchain, linker script, and startup file.
4. Rebuild Forge (`go build`). The catalog is compiled into the binary.
5. Create a project with the new id:

```bash
forge new blink stm32f103c8
```

Unknown ids fail at `forge new` (`Resolve`). `init` also resolves `target.device` from `Forge.toml`.

### Example: another STM32F103 part

```yaml
stm32f103rb:
  <<: *stm32f103
  id: stm32f103rb
  flash_kb: 128
  ram_kb: 20
  pins: 64
  package: LQFP64
  stm32_device: STM32F103xB
```

Same `stm32_device` as `stm32f103r8`, so the existing F1 HAL, CMSIS-Device headers, and `startup_stm32f103xb.s` still apply. `flash_kb: 128` changes the generated `LinkerScript.ld` FLASH length.

## See also

- [Create a project](create-a-project.md)
- [Commands](commands.md)
