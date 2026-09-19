package templates

import (
	"bytes"
	"html/template"
	"strings"
	"testing"
)

func TestFS_ContainsScaffoldTemplates(t *testing.T) {
	want := []string{
		"main.txt.tmpl",
		"cmake/CMakeLists.txt.tmpl",
		"cmake/CMakePresets.json.tmpl",
		"cmake/gcc-arm-none-eabi.cmake.tmpl",
		"cmake/Package.cmake.tmpl",
		"ld/LinkerScript.ld.tmpl",
	}
	for _, name := range want {
		if _, err := FS.ReadFile(name); err != nil {
			t.Errorf("missing %s: %v", name, err)
		}
	}
}

func TestRender_MainByCPU(t *testing.T) {
	cases := []struct {
		cpu, header, fpu, wantMacro, wantHeader string
	}{
		{"cortex-m3", "core_cm3.h", "0U", "__CM3_REV", "core_cm3.h"},
		{"cortex-m4", "core_cm4.h", "1U", "__CM4_REV", "core_cm4.h"},
		{"cortex-m7", "core_cm7.h", "1U", "__CM7_REV", "core_cm7.h"},
	}
	raw, err := FS.ReadFile("main.txt.tmpl")
	if err != nil {
		t.Fatal(err)
	}
	tmpl, err := template.New("main.c").Parse(string(raw))
	if err != nil {
		t.Fatal(err)
	}
	for _, tc := range cases {
		t.Run(tc.cpu, func(t *testing.T) {
			data := sampleData()
			data.TargetCPU = tc.cpu
			data.CMSISCoreHeader = tc.header
			data.FPUPresent = tc.fpu
			var buf bytes.Buffer
			if err := tmpl.Execute(&buf, data); err != nil {
				t.Fatalf("execute: %v", err)
			}
			got := buf.String()
			if !strings.Contains(got, tc.wantMacro) {
				t.Errorf("missing %s in:\n%s", tc.wantMacro, got)
			}
			if !strings.Contains(got, `#include "`+tc.wantHeader+`"`) {
				t.Errorf("missing include %s", tc.wantHeader)
			}
			if !strings.Contains(got, "#define __FPU_PRESENT "+tc.fpu) {
				t.Errorf("missing FPU present %s", tc.fpu)
			}
			if tc.cpu == "cortex-m7" && !strings.Contains(got, "__ICACHE_PRESENT") {
				t.Error("M7 template should define cache macros")
			}
		})
	}
}

func TestRender_PackageCMake(t *testing.T) {
	got := mustRender(t, "cmake/Package.cmake.tmpl", sampleData())
	if !strings.Contains(got, `set(cache_dir "/tmp/forge-cache")`) {
		t.Errorf("cache_dir missing:\n%s", got)
	}
	if !strings.Contains(got, "add_library(cmsis5 INTERFACE)") {
		t.Error("expected INTERFACE cmsis5")
	}
	if !strings.Contains(got, "add_library(stm32f1-hal STATIC") {
		t.Error("expected STATIC HAL")
	}
	if !strings.Contains(got, "USE_HAL_DRIVER") || !strings.Contains(got, "STM32F103xB") {
		t.Error("expected HAL defines")
	}
	if !strings.Contains(got, "target_link_libraries(stm32f1-hal PUBLIC cmsis5") {
		t.Error("expected HAL depends")
	}
}

func TestRender_CMakeListsAndPresets(t *testing.T) {
	lists := mustRender(t, "cmake/CMakeLists.txt.tmpl", sampleData())
	if !strings.Contains(lists, "project(blink") || !strings.Contains(lists, "startup_stm32f103xb.s") {
		t.Errorf("CMakeLists:\n%s", lists)
	}
	if !strings.Contains(lists, "stm32f1-hal") {
		t.Error("CMakeLists should link HAL")
	}

	presets := mustRender(t, "cmake/CMakePresets.json.tmpl", sampleData())
	if !strings.Contains(presets, `"CPU": "cortex-m3"`) {
		t.Error("preset CPU")
	}
	if !strings.Contains(presets, `"STM32_DEVICE": "STM32F103xB"`) {
		t.Error("preset STM32_DEVICE")
	}
	if !strings.Contains(presets, `"FLASH_KB": "64"`) || !strings.Contains(presets, `"RAM_KB": "20"`) {
		t.Error("preset flash/ram")
	}
}

func TestRender_ToolchainAndLinker(t *testing.T) {
	tool := mustRender(t, "cmake/gcc-arm-none-eabi.cmake.tmpl", sampleData())
	if !strings.Contains(tool, `set(CPU "cortex-m3")`) {
		t.Error("toolchain CPU fallback")
	}
	if !strings.Contains(tool, `set(FLOAT_ABI "soft")`) {
		t.Error("toolchain float abi")
	}

	ld := mustRender(t, "ld/LinkerScript.ld.tmpl", sampleData())
	if !strings.Contains(ld, "LENGTH = 64K") || !strings.Contains(ld, "LENGTH = 20K") {
		t.Errorf("linker memory: %s", ld)
	}
	if !strings.Contains(ld, "STM32F103XX") {
		t.Error("linker series comment")
	}
}

func mustRender(t *testing.T, name string, data TemplateData) string {
	t.Helper()
	raw, err := FS.ReadFile(name)
	if err != nil {
		t.Fatalf("read %s: %v", name, err)
	}
	tmpl, err := template.New(name).Parse(string(raw))
	if err != nil {
		t.Fatalf("parse %s: %v", name, err)
	}
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		t.Fatalf("execute %s: %v", name, err)
	}
	return buf.String()
}

func sampleData() TemplateData {
	return TemplateData{
		ProjectName:                 "blink",
		ProjectVersion:              "0.1.0",
		ProjectDescription:          "blink",
		CStandard:                   "11",
		CMakeMinimumRequiredVersion: "3.20",
		CMakeMinimumRequiredMajor:   3,
		CMakeMinimumRequiredMinor:   20,
		ToolchainCompiler:           "gcc-arm-none-eabi",
		TargetDevice:                "stm32f103r8",
		TargetCPU:                   "cortex-m3",
		TargetFPU:                   "",
		TargetFloatABI:              "soft",
		TargetVendor:                "ST",
		TargetFamily:                "STM32F1",
		TargetSeries:                "STM32F103XX",
		TargetFlashKB:               64,
		TargetRAMKB:                 20,
		STM32Device:                 "STM32F103xB",
		CMSISCoreHeader:             "core_cm3.h",
		FPUPresent:                  "0U",
		CacheDir:                    "/tmp/forge-cache",
		Packages: []PackageData{
			{
				Name:        "cmsis5",
				Kind:        "INTERFACE",
				IncludeDirs: []string{"CMSIS_5-5.9.0/CMSIS/Core/Include"},
			},
			{
				Name:        "stm32f1-hal",
				Kind:        "STATIC",
				IncludeDirs: []string{"stm32f1xx-hal-driver-1.1.10/Inc"},
				Sources:     []string{"stm32f1xx-hal-driver-1.1.10/Src/stm32f1xx_hal.c"},
				Defines:     []string{"USE_HAL_DRIVER", "STM32F103xB"},
				Depends:     []string{"cmsis5"},
			},
		},
		AppSources:       []string{"system_stm32f1xx.c", "startup_stm32f103xb.s"},
		DebugBuildType:   "Debug",
		ReleaseBuildType: "Release",
	}
}
