package cmd

import "testing"

func TestParseCMakeVersion(t *testing.T) {
	major, minor := parseCMakeVersion("3.20")
	if major != 3 || minor != 20 {
		t.Errorf("parseCMakeVersion(3.20) = %d.%d", major, minor)
	}
	major, minor = parseCMakeVersion("")
	if major != 3 || minor != 20 {
		t.Errorf("parseCMakeVersion empty = %d.%d", major, minor)
	}
}

func TestBuild_RequiresPreset(t *testing.T) {
	if err := Build(); err == nil {
		t.Fatal("Build: expected error when preset is missing")
	}
}
