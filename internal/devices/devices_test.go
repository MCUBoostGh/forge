package devices

import "testing"

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
	if len(deps) != len(want) {
		t.Fatalf("packages = %v, want %v", deps, want)
	}
	for i, spec := range want {
		if deps[i] != spec {
			t.Fatalf("packages = %v, want %v", deps, want)
		}
	}
}
