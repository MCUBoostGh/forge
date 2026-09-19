package main

import (
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"forge/internal/config"
)

func captureStdout(t *testing.T, fn func()) string {
	t.Helper()
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("pipe: %v", err)
	}
	old := os.Stdout
	os.Stdout = w
	defer func() { os.Stdout = old }()

	fn()
	_ = w.Close()
	b, err := io.ReadAll(r)
	if err != nil {
		t.Fatalf("read: %v", err)
	}
	return string(b)
}

func TestRun_NoArgsShowsHelp(t *testing.T) {
	out := captureStdout(t, func() {
		code, err := run(nil)
		if err != nil {
			t.Errorf("err = %v", err)
		}
		if code != 1 {
			t.Errorf("code = %d, want 1", code)
		}
	})
	if !strings.Contains(out, "Usage:") || !strings.Contains(out, "forge <command>") {
		t.Errorf("help output = %q", out)
	}
}

func TestRun_Help(t *testing.T) {
	out := captureStdout(t, func() {
		code, err := run([]string{"help"})
		if err != nil {
			t.Errorf("err = %v", err)
		}
		if code != 0 {
			t.Errorf("code = %d, want 0", code)
		}
	})
	if !strings.Contains(out, "version") {
		t.Errorf("help = %q", out)
	}
}

func TestRun_Version(t *testing.T) {
	out := captureStdout(t, func() {
		code, err := run([]string{"version"})
		if err != nil {
			t.Errorf("err = %v", err)
		}
		if code != 0 {
			t.Errorf("code = %d, want 0", code)
		}
	})
	if !strings.Contains(out, "Forge version 0.2.0") {
		t.Errorf("version = %q", out)
	}
}

func TestRun_UnknownCommand(t *testing.T) {
	code, err := run([]string{"not-a-command"})
	if code != 1 {
		t.Errorf("code = %d, want 1", code)
	}
	if err == nil || !strings.Contains(err.Error(), "Unknown command") {
		t.Errorf("err = %v", err)
	}
}

func TestRun_NewMissingArgs(t *testing.T) {
	code, err := run([]string{"new"})
	if code != 1 {
		t.Errorf("code = %d, want 1", code)
	}
	if err == nil || !strings.Contains(err.Error(), "Failed to generate new project") {
		t.Errorf("err = %v", err)
	}
}

func TestRun_NewCreatesProject(t *testing.T) {
	dir := t.TempDir()
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chdir(cwd) })
	config.Set(config.Config{})

	code, err := run([]string{"new", "blink", "stm32f103r8"})
	if err != nil {
		t.Fatalf("run new: %v", err)
	}
	if code != 0 {
		t.Fatalf("code = %d", code)
	}
	if _, err := os.Stat(filepath.Join("blink", "Forge.toml")); err != nil {
		t.Fatalf("Forge.toml: %v", err)
	}
}

func TestRun_BuildRequiresPreset(t *testing.T) {
	code, err := run([]string{"build"})
	if code != 1 {
		t.Errorf("code = %d, want 1", code)
	}
	if err == nil {
		t.Fatal("expected error when preset is missing")
	}
}

func TestRun_InitMissingForgeToml(t *testing.T) {
	dir := t.TempDir()
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chdir(cwd) })
	config.Set(config.Config{})

	code, err := run([]string{"init"})
	if code != 1 || err == nil {
		t.Fatalf("code=%d err=%v, want failure", code, err)
	}
}
