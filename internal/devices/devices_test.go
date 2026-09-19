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
