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
