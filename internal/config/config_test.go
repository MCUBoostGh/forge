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
	if config.Toolchain.Compiler != "gcc" {
		t.Errorf("Toolchain.Compiler = %q, want %q", config.Toolchain.Compiler, "gcc")
	}
	if config.CMake.Version != "3.30" {
		t.Errorf("CMake.Version = %q, want %q", config.CMake.Version, "3.30")
	}
	if config.CMake.MinimumRequiredVersion != "3.20" {
		t.Errorf("CMake.MinimumRequiredVersion = %q, want %q", config.CMake.MinimumRequiredVersion, "3.20")
	}
	wantDeps := []string{"cmsis5@5.9.0", "stm32f1-cmsis-device@4.3.5", "stm32f1-hal@1.1.10"}
	if len(config.Dependencies) != len(wantDeps) {
		t.Errorf("Dependencies = %v, want %v", config.Dependencies, wantDeps)
	} else {
		for i, dep := range wantDeps {
			if config.Dependencies[i] != dep {
				t.Errorf("Dependencies = %v, want %v", config.Dependencies, wantDeps)
				break
			}
		}
	}

	var fromFile Config
	if err := toml.Unmarshal(raw, &fromFile); err != nil {
		t.Fatalf("unmarshal Forge.toml: %v", err)
	}
	if fromFile.Project.Name != projectDir {
		t.Errorf("file Project.Name = %q, want %q", fromFile.Project.Name, projectDir)
	}
	if fromFile.Build.System != "cmake" {
		t.Errorf("file Build.System = %q, want %q", fromFile.Build.System, "cmake")
	}
}

func TestNew_PreservesTargetDevice(t *testing.T) {
	root := chdirTemp(t)
	resetConfig()
	config.Target.Kind = "mcu"
	config.Target.Device = "stm32f103r8"
	config.Toolchain.Compiler = "gcc-arm-none-eabi"

	projectDir := filepath.Join(root, "blink")
	if err := os.Mkdir(projectDir, 0755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	if err := New(projectDir); err != nil {
		t.Fatalf("New: %v", err)
	}

	if config.Target.Device != "stm32f103r8" {
		t.Errorf("Target.Device = %q, want stm32f103r8", config.Target.Device)
	}
	if config.Toolchain.Compiler != "gcc-arm-none-eabi" {
		t.Errorf("Toolchain.Compiler = %q, want gcc-arm-none-eabi", config.Toolchain.Compiler)
	}

	raw, err := os.ReadFile(filepath.Join(projectDir, forgeTOMLName))
	if err != nil {
		t.Fatalf("read Forge.toml: %v", err)
	}
	var fromFile Config
	if err := toml.Unmarshal(raw, &fromFile); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if fromFile.Target.Device != "stm32f103r8" {
		t.Errorf("file Target.Device = %q, want stm32f103r8", fromFile.Target.Device)
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

	config.Build.System = "ninja"
	config.Project.Version = "1.2.3"
	if err := Write(); err != nil {
		t.Fatalf("Write: %v", err)
	}

	resetConfig()
	if err := Read(); err != nil {
		t.Fatalf("re-Read: %v", err)
	}
	if config.Build.System != "ninja" {
		t.Errorf("Build.System = %q, want %q", config.Build.System, "ninja")
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
