//go:build integration

package cmd

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func requireSmokeTools(t *testing.T) {
	t.Helper()
	for _, bin := range []string{"cmake", "arm-none-eabi-gcc", "make"} {
		if _, err := exec.LookPath(bin); err != nil {
			t.Skipf("skipping integration test: %s not on PATH", bin)
		}
	}
}

func TestIntegration_NewInitBuild(t *testing.T) {
	requireSmokeTools(t)

	cases := []struct {
		name, device string
	}{
		{"blinkf1", "stm32f103r8"},
		{"blinkg4", "stm32g431rb"},
	}
	for _, tc := range cases {
		t.Run(tc.device, func(t *testing.T) {
			chdirTemp(t)
			resetCmdConfig()

			if err := New(tc.name, tc.device); err != nil {
				t.Fatalf("New: %v", err)
			}
			if err := os.Chdir(tc.name); err != nil {
				t.Fatalf("chdir project: %v", err)
			}
			if err := Init(); err != nil {
				t.Fatalf("Init: %v", err)
			}
			if err := Build("debug"); err != nil {
				t.Fatalf("Build: %v", err)
			}
			elf := filepath.Join("build", "debug", tc.name)
			if _, err := os.Stat(elf); err != nil {
				t.Fatalf("firmware %s: %v", elf, err)
			}
		})
	}
}
