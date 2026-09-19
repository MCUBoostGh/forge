package Package

import (
	"strings"
	"testing"
)

func TestParseSpec(t *testing.T) {
	name, version, err := ParseSpec("stm32f1-hal@1.1.10")
	if err != nil {
		t.Fatalf("ParseSpec: %v", err)
	}
	if name != "stm32f1-hal" || version != "1.1.10" {
		t.Fatalf("got %s@%s, want stm32f1-hal@1.1.10", name, version)
	}
}

func TestRegister_STM32F1HAL(t *testing.T) {
	pkg, err := Register("stm32f1-hal", "1.1.10")
	if err != nil {
		t.Fatalf("Register: %v", err)
	}
	if pkg.Name != "stm32f1-hal" {
		t.Errorf("Name = %q", pkg.Name)
	}
	if pkg.ArchiveRoot != "stm32f1xx-hal-driver-1.1.10" {
		t.Errorf("ArchiveRoot = %q, want stm32f1xx-hal-driver-1.1.10", pkg.ArchiveRoot)
	}
	if !strings.Contains(pkg.URL, "/v1.1.10.tar.gz") {
		t.Errorf("URL = %q, want tag v1.1.10", pkg.URL)
	}
	if len(pkg.Include) != 1 || pkg.Include[0] != "Inc" {
		t.Errorf("Include = %v, want [Inc]", pkg.Include)
	}
	if len(pkg.Sources) != 1 || pkg.Sources[0] != "Src" {
		t.Errorf("Sources = %v, want [Src]", pkg.Sources)
	}
}

func TestRegister_CMSISDeviceF1(t *testing.T) {
	pkg, err := Register("stm32f1-cmsis-device", "4.3.5")
	if err != nil {
		t.Fatalf("Register: %v", err)
	}
	if pkg.ArchiveRoot != "cmsis-device-f1-4.3.5" {
		t.Errorf("ArchiveRoot = %q, want cmsis-device-f1-4.3.5", pkg.ArchiveRoot)
	}
	if pkg.StartupDir != "Source/Templates/gcc" {
		t.Errorf("StartupDir = %q", pkg.StartupDir)
	}
	if len(pkg.ProjectSources) != 1 || pkg.ProjectSources[0] != "Source/Templates/system_stm32f1xx.c" {
		t.Errorf("ProjectSources = %v", pkg.ProjectSources)
	}
	if got := pkg.StartupFileName("STM32F103xB"); got != "startup_stm32f103xb.s" {
		t.Errorf("StartupFileName = %q, want startup_stm32f103xb.s", got)
	}
}

func TestRegister_CMSISDeviceF7(t *testing.T) {
	pkg, err := Register("stm32f7-cmsis-device", "1.2.10")
	if err != nil {
		t.Fatalf("Register: %v", err)
	}
	if pkg.ArchiveRoot != "cmsis-device-f7-1.2.10" {
		t.Errorf("ArchiveRoot = %q", pkg.ArchiveRoot)
	}
	if !strings.Contains(pkg.URL, "/v1.2.10.tar.gz") {
		t.Errorf("URL = %q, want tag v1.2.10", pkg.URL)
	}
	if len(pkg.ProjectSources) != 1 || pkg.ProjectSources[0] != "Source/Templates/system_stm32f7xx.c" {
		t.Errorf("ProjectSources = %v", pkg.ProjectSources)
	}
	if got := pkg.StartupFileName("STM32F746xx"); got != "startup_stm32f746xx.s" {
		t.Errorf("StartupFileName = %q, want startup_stm32f746xx.s", got)
	}
}

func TestRegister_CMSISDeviceG4(t *testing.T) {
	pkg, err := Register("stm32g4-cmsis-device", "1.2.6")
	if err != nil {
		t.Fatalf("Register: %v", err)
	}
	if pkg.ArchiveRoot != "cmsis-device-g4-1.2.6" {
		t.Errorf("ArchiveRoot = %q", pkg.ArchiveRoot)
	}
	if got := pkg.StartupFileName("STM32G431xx"); got != "startup_stm32g431xx.s" {
		t.Errorf("StartupFileName = %q, want startup_stm32g431xx.s", got)
	}
}

func TestRegister_STM32F7HAL(t *testing.T) {
	pkg, err := Register("stm32f7-hal", "1.3.3")
	if err != nil {
		t.Fatalf("Register: %v", err)
	}
	if pkg.ArchiveRoot != "stm32f7xx-hal-driver-1.3.3" {
		t.Errorf("ArchiveRoot = %q", pkg.ArchiveRoot)
	}
}

func TestRegister_STM32G4HAL(t *testing.T) {
	pkg, err := Register("stm32g4-hal", "1.2.6")
	if err != nil {
		t.Fatalf("Register: %v", err)
	}
	if pkg.ArchiveRoot != "stm32g4xx-hal-driver-1.2.6" {
		t.Errorf("ArchiveRoot = %q", pkg.ArchiveRoot)
	}
}
