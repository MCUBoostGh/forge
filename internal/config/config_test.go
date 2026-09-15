package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/pelletier/go-toml/v2"
)

func resetConfig() {
	config = Config{}
}

func chdirTemp(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	if err := os.Chdir(dir); err != nil {
		t.Fatalf("chdir %s: %v", dir, err)
	}
	t.Cleanup(func() {
		_ = os.Chdir(cwd)
	})
	return dir
}

func TestNew_GeneratesDefaultConfig(t *testing.T) {
	root := chdirTemp(t)
	resetConfig()

	projectDir := filepath.Join(root, "demo")
	if err := os.Mkdir(projectDir, 0755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}

	if err := New(projectDir); err != nil {
		t.Fatalf("New: %v", err)
	}

	tomlPath := filepath.Join(projectDir, forgeTOMLName)
	raw, err := os.ReadFile(tomlPath)
	if err != nil {
		t.Fatalf("read generated Forge.toml: %v", err)
	}

	if config.Project.Name != projectDir {
		t.Errorf("Project.Name = %q, want %q", config.Project.Name, projectDir)
	}
	if config.Project.Version != "0.1.0" {
		t.Errorf("Project.Version = %q, want %q", config.Project.Version, "0.1.0")
	}
	if config.Build.System != "cmake" {
		t.Errorf("Build.System = %q, want %q", config.Build.System, "cmake")
	}
	if config.Build.Type != "debug" {
		t.Errorf("Build.Type = %q, want %q", config.Build.Type, "debug")
	}
	if config.Toolchain.Compiler != "gcc" {
		t.Errorf("Toolchain.Compiler = %q, want %q", config.Toolchain.Compiler, "gcc")
	}
	if config.CMake.Version != "3.30" {
		t.Errorf("CMake.Version = %q, want %q", config.CMake.Version, "3.30")
	}
	if config.CMake.MinimumRequiredVersion != "3.20" {
		t.Errorf("CMake.MinimumRequiredVersion = %q, want %q", config.CMake.MinimumRequiredVersion, "3.20")
	}

	var fromFile Config
	if err := toml.Unmarshal(raw, &fromFile); err != nil {
		t.Fatalf("unmarshal Forge.toml: %v", err)
	}
	if fromFile.Project.Name != projectDir {
		t.Errorf("file Project.Name = %q, want %q", fromFile.Project.Name, projectDir)
	}
	if fromFile.Build.Type != "debug" {
		t.Errorf("file Build.Type = %q, want %q", fromFile.Build.Type, "debug")
	}
}

func TestRead_LoadsConfigBuffer(t *testing.T) {
	root := chdirTemp(t)
	resetConfig()

	projectDir := filepath.Join(root, "demo")
	if err := os.Mkdir(projectDir, 0755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	if err := New(projectDir); err != nil {
		t.Fatalf("New: %v", err)
	}

	if err := os.Chdir(projectDir); err != nil {
		t.Fatalf("chdir project: %v", err)
	}

	resetConfig()
	if err := Read(); err != nil {
		t.Fatalf("Read: %v", err)
	}

	if config.Project.Name != projectDir {
		t.Errorf("Project.Name = %q, want %q", config.Project.Name, projectDir)
	}
	if config.Project.Version != "0.1.0" {
		t.Errorf("Project.Version = %q, want %q", config.Project.Version, "0.1.0")
	}
	if config.Build.System != "cmake" {
		t.Errorf("Build.System = %q, want %q", config.Build.System, "cmake")
	}
	if config.Build.Type != "debug" {
		t.Errorf("Build.Type = %q, want %q", config.Build.Type, "debug")
	}
	if config.Toolchain.Compiler != "gcc" {
		t.Errorf("Toolchain.Compiler = %q, want %q", config.Toolchain.Compiler, "gcc")
	}
}

func TestWrite_PersistsConfigBuffer(t *testing.T) {
	root := chdirTemp(t)
	resetConfig()

	projectDir := filepath.Join(root, "demo")
	if err := os.Mkdir(projectDir, 0755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	if err := New(projectDir); err != nil {
		t.Fatalf("New: %v", err)
	}
	if err := os.Chdir(projectDir); err != nil {
		t.Fatalf("chdir project: %v", err)
	}

	resetConfig()
	if err := Read(); err != nil {
		t.Fatalf("Read: %v", err)
	}

	config.Build.Type = "release"
	config.Project.Version = "1.2.3"
	if err := Write(); err != nil {
		t.Fatalf("Write: %v", err)
	}

	resetConfig()
	if err := Read(); err != nil {
		t.Fatalf("re-Read: %v", err)
	}
	if config.Build.Type != "release" {
		t.Errorf("Build.Type = %q, want %q", config.Build.Type, "release")
	}
	if config.Project.Version != "1.2.3" {
		t.Errorf("Project.Version = %q, want %q", config.Project.Version, "1.2.3")
	}
}

func TestRead_MissingFile(t *testing.T) {
	chdirTemp(t)
	resetConfig()

	if err := Read(); err == nil {
		t.Fatal("Read: expected error when Forge.toml is missing")
	}
}
