# Arm bare-metal toolchain (gcc-arm-none-eabi) for STM32 / Cortex-M.
# Configure with -DCPU=..., -DFPU=..., -DFLOAT_ABI=..., -DSTM32_DEVICE=...
# Example: cmake --preset stm32 -DCPU=cortex-m3 -DSTM32_DEVICE=STM32F103xB

set(CMAKE_SYSTEM_NAME Generic)
set(CMAKE_SYSTEM_PROCESSOR arm)

# Skip executable link checks; bare-metal needs a linker script.
set(CMAKE_TRY_COMPILE_TARGET_TYPE STATIC_LIBRARY)

# Toolchain prefix: "arm-none-eabi-" or an absolute path like "/usr/bin/arm-none-eabi-"
if(NOT DEFINED TOOLCHAIN_PREFIX)
	set(TOOLCHAIN_PREFIX "arm-none-eabi-" CACHE STRING "gcc-arm-none-eabi toolchain prefix")
endif()

set(CMAKE_C_COMPILER   "${TOOLCHAIN_PREFIX}gcc"   CACHE FILEPATH "" FORCE)
set(CMAKE_CXX_COMPILER "${TOOLCHAIN_PREFIX}g++"   CACHE FILEPATH "" FORCE)
set(CMAKE_ASM_COMPILER "${TOOLCHAIN_PREFIX}gcc"   CACHE FILEPATH "" FORCE)
set(CMAKE_AR           "${TOOLCHAIN_PREFIX}ar"    CACHE FILEPATH "" FORCE)
set(CMAKE_RANLIB       "${TOOLCHAIN_PREFIX}ranlib" CACHE FILEPATH "" FORCE)
set(CMAKE_OBJCOPY      "${TOOLCHAIN_PREFIX}objcopy" CACHE FILEPATH "" FORCE)
set(CMAKE_OBJDUMP      "${TOOLCHAIN_PREFIX}objdump" CACHE FILEPATH "" FORCE)
set(CMAKE_SIZE         "${TOOLCHAIN_PREFIX}size"  CACHE FILEPATH "" FORCE)
set(CMAKE_STRIP        "${TOOLCHAIN_PREFIX}strip" CACHE FILEPATH "" FORCE)

# CPU / FPU (override per board; defaults suit STM32F1 Cortex-M3)
if(NOT DEFINED CPU)
	set(CPU "cortex-m3" CACHE STRING "ARM CPU core (e.g. cortex-m0, cortex-m3, cortex-m4, cortex-m7)")
endif()

set(_ARM_FLAGS "-mcpu=${CPU} -mthumb")

if(DEFINED FPU AND NOT "${FPU}" STREQUAL "")
	string(APPEND _ARM_FLAGS " -mfpu=${FPU}")
endif()

if(DEFINED FLOAT_ABI AND NOT "${FLOAT_ABI}" STREQUAL "")
	string(APPEND _ARM_FLAGS " -mfloat-abi=${FLOAT_ABI}")
endif()

set(CMAKE_C_FLAGS_INIT   "${_ARM_FLAGS}")
set(CMAKE_CXX_FLAGS_INIT "${_ARM_FLAGS}")
set(CMAKE_ASM_FLAGS_INIT "${_ARM_FLAGS}")
set(CMAKE_EXE_LINKER_FLAGS_INIT "${_ARM_FLAGS} --specs=nano.specs --specs=nosys.specs")

# Optional STM32 device define (e.g. STM32F103xB, STM32F411xE)
if(DEFINED STM32_DEVICE AND NOT "${STM32_DEVICE}" STREQUAL "")
	string(TOUPPER "${STM32_DEVICE}" _STM32_DEVICE_UPPER)
	add_compile_definitions(${_STM32_DEVICE_UPPER})
endif()

set(CMAKE_FIND_ROOT_PATH_MODE_PROGRAM NEVER)
set(CMAKE_FIND_ROOT_PATH_MODE_LIBRARY ONLY)
set(CMAKE_FIND_ROOT_PATH_MODE_INCLUDE ONLY)
set(CMAKE_FIND_ROOT_PATH_MODE_PACKAGE ONLY)

# After linking: write .hex / .bin and print section sizes.
function(forge_post_build_hex_size target)
	add_custom_command(TARGET ${target} POST_BUILD
		COMMAND ${CMAKE_OBJCOPY} -O ihex
			$<TARGET_FILE:${target}>
			$<TARGET_FILE_DIR:${target}>/$<TARGET_FILE_BASE_NAME:${target}>.hex
		COMMAND ${CMAKE_OBJCOPY} -O binary
			$<TARGET_FILE:${target}>
			$<TARGET_FILE_DIR:${target}>/$<TARGET_FILE_BASE_NAME:${target}>.bin
		COMMAND ${CMAKE_SIZE} --format=berkeley $<TARGET_FILE:${target}>
		COMMENT "HEX/BIN + size: ${target}"
		VERBATIM
	)
endfunction()
