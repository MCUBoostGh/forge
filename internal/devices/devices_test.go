package devices

import (
	"slices"
	"testing"
)

func equalSpecs(a, b []string) bool {
	return slices.Equal(a, b)
}

func TestResolve_CatalogID(t *testing.T) {
	dev, err := Resolve("stm32f103r8")
	if err != nil {
		t.Fatalf("Resolve: %v", err)
	}
	if dev.ID != "stm32f103r8" {
		t.Errorf("ID = %q", dev.ID)
	}
	if dev.STM32Device != "STM32F103xB" {
		t.Errorf("STM32Device = %q", dev.STM32Device)
	}
}

func TestResolve_Alias(t *testing.T) {
	dev, err := Resolve("bluepill")
	if err != nil {
		t.Fatalf("Resolve alias: %v", err)
	}
	if dev.ID != "stm32f103r8" {
		t.Errorf("ID = %q, want stm32f103r8", dev.ID)
	}
}

func TestResolve_CaseInsensitive(t *testing.T) {
	dev, err := Resolve("STM32F103C8")
	if err != nil {
		t.Fatalf("Resolve: %v", err)
	}
	if dev.ID != "stm32f103c8" {
		t.Errorf("ID = %q, want stm32f103c8", dev.ID)
	}
}

func TestResolve_Unknown(t *testing.T) {
	_, err := Resolve("not-a-board")
	if err == nil {
		t.Fatal("expected error for unknown device")
	}
}

func TestLookup_UsesResolve(t *testing.T) {
	dev, err := Lookup("bluepill")
	if err != nil {
		t.Fatalf("Lookup alias: %v", err)
	}
	if dev.ID != "stm32f103r8" {
		t.Errorf("ID = %q, want stm32f103r8", dev.ID)
	}
}

func TestResolve_STM32F7(t *testing.T) {
	dev, err := Resolve("nucleo-f746zg")
	if err != nil {
		t.Fatalf("Resolve F7 alias: %v", err)
	}
	if dev.ID != "stm32f746zg" {
		t.Errorf("ID = %q, want stm32f746zg", dev.ID)
	}
	if dev.CPU != "cortex-m7" {
		t.Errorf("CPU = %q, want cortex-m7", dev.CPU)
	}
	if dev.STM32Device != "STM32F746xx" {
		t.Errorf("STM32Device = %q", dev.STM32Device)
	}
	if dev.FPU != "fpv5-sp-d16" {
		t.Errorf("FPU = %q", dev.FPU)
	}
	want := []string{"cmsis5@5.9.0", "stm32f7-cmsis-device@1.2.10", "stm32f7-hal@1.3.3"}
	deps, err := dev.DefaultDependencies()
	if err != nil {
		t.Fatalf("DefaultDependencies: %v", err)
	}
	if len(deps) != len(want) {
		t.Fatalf("packages = %v, want %v", deps, want)
	}
	for i, spec := range want {
		if deps[i] != spec {
			t.Fatalf("packages = %v, want %v", deps, want)
		}
	}
}

func TestResolve_STM32G4(t *testing.T) {
	dev, err := Resolve("stm32g474re")
	if err != nil {
		t.Fatalf("Resolve G4: %v", err)
	}
	if dev.CPU != "cortex-m4" {
		t.Errorf("CPU = %q, want cortex-m4", dev.CPU)
	}
	if dev.STM32Device != "STM32G474xx" {
		t.Errorf("STM32Device = %q", dev.STM32Device)
	}
	if dev.FPU != "fpv4-sp-d16" {
		t.Errorf("FPU = %q", dev.FPU)
	}
	want := []string{"cmsis5@5.9.0", "stm32g4-cmsis-device@1.2.6", "stm32g4-hal@1.2.6"}
	deps, err := dev.DefaultDependencies()
	if err != nil {
		t.Fatalf("DefaultDependencies: %v", err)
	}
	if !equalSpecs(deps, want) {
		t.Fatalf("packages = %v, want %v", deps, want)
	}
}

func TestCatalog_ShippedFamilies(t *testing.T) {
	f1 := []string{"cmsis5@5.9.0", "stm32f1-cmsis-device@4.3.5", "stm32f1-hal@1.1.10"}
	f7 := []string{"cmsis5@5.9.0", "stm32f7-cmsis-device@1.2.10", "stm32f7-hal@1.3.3"}
	g4 := []string{"cmsis5@5.9.0", "stm32g4-cmsis-device@1.2.6", "stm32g4-hal@1.2.6"}

	cases := []struct {
		input, id, family, cpu, stm32, fpu string
		flash, ram                         int
		packages                           []string
	}{
		{"stm32f103r8", "stm32f103r8", "STM32F1", "cortex-m3", "STM32F103xB", "", 64, 20, f1},
		{"stm32f103c8", "stm32f103c8", "STM32F1", "cortex-m3", "STM32F103xB", "", 64, 20, f1},
		{"stm32f746zg", "stm32f746zg", "STM32F7", "cortex-m7", "STM32F746xx", "fpv5-sp-d16", 1024, 320, f7},
		{"nucleo-f767zi", "stm32f767zi", "STM32F7", "cortex-m7", "STM32F767xx", "fpv5-d16", 2048, 512, f7},
		{"nucleo-g431rb", "stm32g431rb", "STM32G4", "cortex-m4", "STM32G431xx", "fpv4-sp-d16", 128, 32, g4},
		{"stm32g474re", "stm32g474re", "STM32G4", "cortex-m4", "STM32G474xx", "fpv4-sp-d16", 512, 128, g4},
	}
	for _, tc := range cases {
		t.Run(tc.input, func(t *testing.T) {
			dev, err := Resolve(tc.input)
			if err != nil {
				t.Fatalf("Resolve(%q): %v", tc.input, err)
			}
			if dev.ID != tc.id || dev.Family != tc.family || dev.CPU != tc.cpu {
				t.Errorf("id/family/cpu = %s %s %s", dev.ID, dev.Family, dev.CPU)
			}
			if dev.STM32Device != tc.stm32 || dev.FPU != tc.fpu {
				t.Errorf("stm32/fpu = %s %s", dev.STM32Device, dev.FPU)
			}
			if dev.FlashKB != tc.flash || dev.RAMKB != tc.ram {
				t.Errorf("flash/ram = %d/%d", dev.FlashKB, dev.RAMKB)
			}
			deps, err := dev.DefaultDependencies()
			if err != nil {
				t.Fatalf("DefaultDependencies: %v", err)
			}
			if !equalSpecs(deps, tc.packages) {
				t.Errorf("packages = %v, want %v", deps, tc.packages)
			}
		})
	}
}

func TestResolve_Empty(t *testing.T) {
	if _, err := Resolve("  "); err == nil {
		t.Fatal("expected error for empty device")
	}
}

func TestResolve_AnchorIDSkipped(t *testing.T) {
	if _, err := Resolve("x-stm32f1"); err == nil {
		t.Fatal("expected error for YAML anchor id")
	}
}

func TestDefaultDependencies_Empty(t *testing.T) {
	if _, err := (Catalog{ID: "none"}).DefaultDependencies(); err == nil {
		t.Fatal("expected error when packages are missing")
	}
}

func TestDefaultDependencies_ReturnsCopy(t *testing.T) {
	dev, err := Resolve("stm32f103r8")
	if err != nil {
		t.Fatalf("Resolve: %v", err)
	}
	deps, err := dev.DefaultDependencies()
	if err != nil {
		t.Fatalf("DefaultDependencies: %v", err)
	}
	deps[0] = "mutated"
	again, err := dev.DefaultDependencies()
	if err != nil {
		t.Fatalf("DefaultDependencies: %v", err)
	}
	if again[0] == "mutated" {
		t.Fatal("DefaultDependencies must not expose catalog slice")
	}
}

func TestList_FiltersFamily(t *testing.T) {
	all := List(Filter{})
	if len(all) < 6 {
		t.Fatalf("List all = %d, want at least 6 shipped boards", len(all))
	}
	f7 := List(Filter{Family: "STM32F7"})
	if len(f7) != 2 {
		t.Fatalf("List F7 = %d, want 2", len(f7))
	}
	g4 := List(Filter{Family: "stm32g4"})
	if len(g4) != 2 {
		t.Fatalf("List G4 = %d, want 2", len(g4))
	}
	f1 := List(Filter{Family: "STM32F1", Series: "STM32F103XX"})
	if len(f1) != 2 {
		t.Fatalf("List F103 = %d, want 2", len(f1))
	}
	for _, d := range all {
		if d.ID == "x-stm32f1" || d.ID == "x-stm32f7" || d.ID == "x-stm32g4" {
			t.Fatalf("List included anchor %s", d.ID)
		}
	}
}
