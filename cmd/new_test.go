package cmd

import "testing"

func TestParseNewArgs_Flag(t *testing.T) {
	name, device, err := parseNewArgs([]string{"blink", "--device", "stm32f103r8"})
	if err != nil {
		t.Fatalf("parseNewArgs: %v", err)
	}
	if name != "blink" || device != "stm32f103r8" {
		t.Fatalf("got %q %q", name, device)
	}
}

func TestParseNewArgs_Positional(t *testing.T) {
	name, device, err := parseNewArgs([]string{"blink", "bluepill"})
	if err != nil {
		t.Fatalf("parseNewArgs: %v", err)
	}
	if name != "blink" || device != "bluepill" {
		t.Fatalf("got %q %q", name, device)
	}
}

func TestParseNewArgs_MissingDevice(t *testing.T) {
	if _, _, err := parseNewArgs([]string{"blink"}); err == nil {
		t.Fatal("expected error when device is missing")
	}
}

func TestParseNewArgs_UnknownFlag(t *testing.T) {
	if _, _, err := parseNewArgs([]string{"blink", "--board", "x"}); err == nil {
		t.Fatal("expected error for unknown flag")
	}
}

func TestParseNewArgs_RequiresName(t *testing.T) {
	if _, _, err := parseNewArgs(nil); err == nil {
		t.Fatal("expected error when name is missing")
	}
	if _, _, err := parseNewArgs([]string{"--device", "stm32f103r8"}); err == nil {
		t.Fatal("expected error when name is a flag")
	}
}

func TestParseNewArgs_EmptyDeviceFlag(t *testing.T) {
	if _, _, err := parseNewArgs([]string{"blink", "--device"}); err == nil {
		t.Fatal("expected error when --device has no value")
	}
}

func TestParseNewArgs_ExtraArgs(t *testing.T) {
	if _, _, err := parseNewArgs([]string{"blink", "stm32f103r8", "extra"}); err == nil {
		t.Fatal("expected error for extra positional args")
	}
}
