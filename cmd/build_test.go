package cmd

import (
	"os"
	"strings"
	"testing"

	"forge/internal/toolchain"
)

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

func TestRunCommand_UsesCachedCMake(t *testing.T) {
	tool, err := toolchain.Lookup("cmake")
	if err != nil {
		t.Fatalf("Lookup cmake: %v", err)
	}
	bin, err := toolchain.BinaryPath(tool)
	if err != nil {
		t.Fatalf("BinaryPath: %v", err)
	}
	if _, err := os.Stat(bin); err != nil {
		t.Skip("cmake is not installed in the toolchain cache")
	}
	if err := runCommand("cmake", "--version"); err != nil {
		t.Fatal(err)
	}
}

func TestCommandEnv_PrependsInstalledBins(t *testing.T) {
	dirs, err := toolchain.BinDirs()
	if err != nil {
		t.Fatalf("BinDirs: %v", err)
	}
	env, err := commandEnv()
	if err != nil {
		t.Fatalf("commandEnv: %v", err)
	}
	var paths []string
	for _, entry := range env {
		if len(entry) >= 5 && strings.EqualFold(entry[:5], "PATH=") {
			paths = append(paths, entry[5:])
		}
	}
	if len(paths) != 1 {
		t.Fatalf("PATH entries = %d, want 1", len(paths))
	}
	got := paths[0]
	for _, dir := range dirs {
		if !strings.Contains(got, dir) {
			t.Fatalf("PATH %q missing %q", got, dir)
		}
	}
	if len(dirs) > 0 && !strings.HasPrefix(got, dirs[0]) {
		t.Fatalf("PATH %q does not start with %q", got, dirs[0])
	}
}
