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

func TestParseSpec_Invalid(t *testing.T) {
	for _, spec := range []string{"", "stm32f1-hal", "stm32f1-hal@", "@1.1.10", "  "} {
		if _, _, err := ParseSpec(spec); err == nil {
			t.Errorf("ParseSpec(%q): expected error", spec)
		}
	}
}

func TestRegister_Unknown(t *testing.T) {
	if _, err := Register("not-a-package", "1.0.0"); err == nil {
		t.Fatal("Register: expected error for unknown package")
	}
}

func TestRegister_FamilyCatalog(t *testing.T) {
	cases := []struct {
		name, version, archiveRoot, urlTag string
		startupDev, startupFile            string
		system                             string
		depends                            []string
	}{
		{
			name: "cmsis5", version: "5.9.0", archiveRoot: "CMSIS_5-5.9.0",
			urlTag: "5.9.0.tar.gz",
		},
		{
			name: "stm32f1-cmsis-device", version: "4.3.5", archiveRoot: "cmsis-device-f1-4.3.5",
			urlTag: "/v4.3.5.tar.gz", startupDev: "STM32F103xB", startupFile: "startup_stm32f103xb.s",
			system: "Source/Templates/system_stm32f1xx.c", depends: []string{"cmsis5"},
		},
		{
			name: "stm32f1-hal", version: "1.1.10", archiveRoot: "stm32f1xx-hal-driver-1.1.10",
			urlTag: "/v1.1.10.tar.gz", depends: []string{"cmsis5", "stm32f1-cmsis-device"},
		},
		{
			name: "stm32f7-cmsis-device", version: "1.2.10", archiveRoot: "cmsis-device-f7-1.2.10",
			urlTag: "/v1.2.10.tar.gz", startupDev: "STM32F746xx", startupFile: "startup_stm32f746xx.s",
			system: "Source/Templates/system_stm32f7xx.c", depends: []string{"cmsis5"},
		},
		{
			name: "stm32f7-hal", version: "1.3.3", archiveRoot: "stm32f7xx-hal-driver-1.3.3",
			urlTag: "/v1.3.3.tar.gz", depends: []string{"cmsis5", "stm32f7-cmsis-device"},
		},
		{
			name: "stm32g4-cmsis-device", version: "1.2.6", archiveRoot: "cmsis-device-g4-1.2.6",
			urlTag: "/v1.2.6.tar.gz", startupDev: "STM32G431xx", startupFile: "startup_stm32g431xx.s",
			system: "Source/Templates/system_stm32g4xx.c", depends: []string{"cmsis5"},
		},
		{
			name: "stm32g4-hal", version: "1.2.6", archiveRoot: "stm32g4xx-hal-driver-1.2.6",
			urlTag: "/v1.2.6.tar.gz", depends: []string{"cmsis5", "stm32g4-cmsis-device"},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name+"@"+tc.version, func(t *testing.T) {
			pkg, err := Register(tc.name, tc.version)
			if err != nil {
				t.Fatalf("Register: %v", err)
			}
			if pkg.Name != tc.name {
				t.Errorf("Name = %q", pkg.Name)
			}
			if pkg.Version != tc.version {
				t.Errorf("Version = %q", pkg.Version)
			}
			if pkg.ArchiveRoot != tc.archiveRoot {
				t.Errorf("ArchiveRoot = %q, want %q", pkg.ArchiveRoot, tc.archiveRoot)
			}
			if !strings.Contains(pkg.URL, tc.urlTag) {
				t.Errorf("URL = %q, want substring %q", pkg.URL, tc.urlTag)
			}
			if tc.system != "" {
				if len(pkg.ProjectSources) != 1 || pkg.ProjectSources[0] != tc.system {
					t.Errorf("ProjectSources = %v, want [%s]", pkg.ProjectSources, tc.system)
				}
			}
			if tc.startupDev != "" {
				if got := pkg.StartupFileName(tc.startupDev); got != tc.startupFile {
					t.Errorf("StartupFileName = %q, want %q", got, tc.startupFile)
				}
			}
			if tc.depends != nil && !slicesEqual(pkg.Depends, tc.depends) {
				t.Errorf("Depends = %v, want %v", pkg.Depends, tc.depends)
			}
		})
	}
}

func slicesEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func TestRegister_TagAlreadyPrefixed(t *testing.T) {
	pkg, err := Register("stm32g4-hal", "v1.2.6")
	if err != nil {
		t.Fatalf("Register: %v", err)
	}
	if strings.Contains(pkg.URL, "/vv1.2.6.tar.gz") {
		t.Errorf("URL double prefix: %q", pkg.URL)
	}
	if !strings.Contains(pkg.URL, "/v1.2.6.tar.gz") {
		t.Errorf("URL = %q, want tag v1.2.6", pkg.URL)
	}
	if pkg.ArchiveRoot != "stm32g4xx-hal-driver-v1.2.6" {
		t.Errorf("ArchiveRoot = %q (version string is used as-is)", pkg.ArchiveRoot)
	}
}

func TestStartupFileName_Empty(t *testing.T) {
	pkg := Package{StartupDir: "Source/Templates/gcc"}
	if got := pkg.StartupFileName(""); got != "" {
		t.Errorf("StartupFileName empty device = %q", got)
	}
	pkg.StartupDir = ""
	if got := pkg.StartupFileName("STM32F103xB"); got != "" {
		t.Errorf("StartupFileName empty dir = %q", got)
	}
}

func TestProjectFileNames_IncludesStartup(t *testing.T) {
	pkg, err := Register("stm32f7-cmsis-device", "1.2.10")
	if err != nil {
		t.Fatalf("Register: %v", err)
	}
	names := pkg.ProjectFileNames("STM32F767xx")
	want := []string{"system_stm32f7xx.c", "startup_stm32f767xx.s"}
	if !slicesEqual(names, want) {
		t.Errorf("ProjectFileNames = %v, want %v", names, want)
	}
}
