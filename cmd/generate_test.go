package cmd

import (
	"os"
	"path/filepath"
	"testing"

	"forge/internal/config"

	"github.com/pelletier/go-toml/v2"
)

func resetCmdConfig() {
	config.Set(config.Config{})
}

func chdirTemp(t *testing.T) {
	t.Helper()
	dir := t.TempDir()
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	if err := os.Chdir(dir); err != nil {
		t.Fatalf("chdir: %v", err)
	}
	t.Cleanup(func() {
		_ = os.Chdir(cwd)
	})
}

func TestCMSISCoreHeader(t *testing.T) {
	cases := map[string]string{
		"cortex-m3":     "core_cm3.h",
		"CORTEX-M4":     "core_cm4.h",
		"cortex-m7":     "core_cm7.h",
		"cortex-m0plus": "core_cm0plus.h",
		"unknown":       "core_cm3.h",
	}
	for cpu, want := range cases {
		if got := cmsisCoreHeader(cpu); got != want {
			t.Errorf("cmsisCoreHeader(%q) = %q, want %q", cpu, got, want)
		}
	}
}

func TestCompilersMap_CortexM(t *testing.T) {
	for _, cpu := range []string{"cortex-m3", "cortex-m4", "cortex-m7"} {
		if compilersMap[cpu] != "gcc-arm-none-eabi" {
			t.Errorf("compilersMap[%q] = %q, want gcc-arm-none-eabi", cpu, compilersMap[cpu])
		}
	}
}

func TestNew_WritesFamilyDependencies(t *testing.T) {
	f1 := []string{"cmsis5@5.9.0", "stm32f1-cmsis-device@4.3.5", "stm32f1-hal@1.1.10"}
	f7 := []string{"cmsis5@5.9.0", "stm32f7-cmsis-device@1.2.10", "stm32f7-hal@1.3.3"}
	g4 := []string{"cmsis5@5.9.0", "stm32g4-cmsis-device@1.2.6", "stm32g4-hal@1.2.6"}

	cases := []struct {
		name, device, wantID string
		wantDeps             []string
	}{
		{"blinkf1", "stm32f103r8", "stm32f103r8", f1},
		{"blinkf7", "nucleo-f746zg", "stm32f746zg", f7},
		{"blinkg4", "stm32g431rb", "stm32g431rb", g4},
	}
	for _, tc := range cases {
		t.Run(tc.device, func(t *testing.T) {
			chdirTemp(t)
			resetCmdConfig()
			if err := New(tc.name, tc.device); err != nil {
				t.Fatalf("New: %v", err)
			}
			cfg := config.Get()
			if cfg.Target.Device != tc.wantID {
				t.Errorf("Target.Device = %q, want %q", cfg.Target.Device, tc.wantID)
			}
			if cfg.Toolchain.Compiler != "gcc-arm-none-eabi" {
				t.Errorf("Compiler = %q", cfg.Toolchain.Compiler)
			}
			if len(cfg.Dependencies) != len(tc.wantDeps) {
				t.Fatalf("Dependencies = %v, want %v", cfg.Dependencies, tc.wantDeps)
			}
			for i, spec := range tc.wantDeps {
				if cfg.Dependencies[i] != spec {
					t.Fatalf("Dependencies = %v, want %v", cfg.Dependencies, tc.wantDeps)
				}
			}

			raw, err := os.ReadFile(filepath.Join(tc.name, "Forge.toml"))
			if err != nil {
				t.Fatalf("read Forge.toml: %v", err)
			}
			var fromFile config.Config
			if err := toml.Unmarshal(raw, &fromFile); err != nil {
				t.Fatalf("unmarshal: %v", err)
			}
			if fromFile.Target.Device != tc.wantID {
				t.Errorf("file device = %q, want %q", fromFile.Target.Device, tc.wantID)
			}
			if len(fromFile.Dependencies) != len(tc.wantDeps) || fromFile.Dependencies[len(tc.wantDeps)-1] != tc.wantDeps[len(tc.wantDeps)-1] {
				t.Errorf("file Dependencies = %v, want %v", fromFile.Dependencies, tc.wantDeps)
			}
		})
	}
}

func TestNew_UnknownDevice(t *testing.T) {
	chdirTemp(t)
	resetCmdConfig()
	if err := New("blink", "not-a-board"); err == nil {
		t.Fatal("New: expected error for unknown device")
	}
}

func TestNew_ExistingDirectory(t *testing.T) {
	chdirTemp(t)
	resetCmdConfig()
	if err := os.Mkdir("blink", 0755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	if err := New("blink", "stm32f103r8"); err == nil {
		t.Fatal("New: expected error when project directory exists")
	}
}

func TestRelCachePaths(t *testing.T) {
	got := relCachePaths("/cache", []string{"/cache/a/b", "/other/x"})
	if len(got) != 2 || got[0] != "a/b" {
		t.Fatalf("relCachePaths = %v", got)
	}
}
