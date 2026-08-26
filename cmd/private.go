package cmd

import (
	"forge/contents"
	"os"
	"os/exec"
)

var listFolders = []string{"src", "include", "cmake"}

var listFilesContentMap = map[string]string{
	"main.c":                        contents.MainCContent,
	"CMakeLists.txt":                contents.CMakeListsContent,
	"CMakePresets.json":             contents.CMakePresetsContent,
	"cmake/gcc-arm-none-eabi.cmake": contents.GccArmNoneEabiCmakeContent,
}

var compilersMap = map[string]string{
	"cortex-m0": "gcc-arm-none-eabi",
	"cortex-m3": "gcc-arm-none-eabi",
	"x86_64":  "gcc",
}

type projectData struct {
	Name         string
	Version      string
	BuildType    string
	Architecture string
	Target       *targetData
	Compiler     string
	BuildDir     string
}

type targetData struct {
	Device string
	CPU    string
	Vendor string
	Family string
	Series string
}

type tomlConfig struct {
	Project struct {
		Name    string `toml:"name"`
		Version string `toml:"version"`
	} `toml:"project"`

	Target struct {
		Kind   string `toml:"kind"`   // mcu, fpga, etc.
		Device string `toml:"device"` // stm32f407vg, etc.
		Board  string `toml:"board"`  // stm32f407vg, etc.
	} `toml:"target"`

	Build struct {
		System string `toml:"system"` //cmake
		Type   string `toml:"type"`   // debug, release, relwithdebinfo, minsizerel
	} `toml:"build"`

	Toolchain struct {
		Compiler string `toml:"compiler"` // gcc, clang, etc.
	} `toml:"toolchain"`

	Dependencies []string `toml:"dependencies"`

	CMake struct {
		Version                string `toml:"version"`
		MinimumRequiredVersion string `toml:"minimum_required_version"`
	} `toml:"cmake"`
}

var forgeTOMLName string = "Forge.toml"
var config = tomlConfig{}
func getConfig() *tomlConfig {
	return &config
}
func setConfigDefaults() {
	config.Project.Name = "MyProject"
	config.Project.Version = "0.1.0"
	config.Build.System = "cmake"
	config.Build.Type = "debug"
	config.Toolchain.Compiler = "gcc"
	config.Dependencies = []string{}
	config.CMake.Version = "3.30"
}
func runCommand(command string, args ...string) error {
	cmd := exec.Command(command, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	return cmd.Run()
}

func setConfig(path string) {

	config.Project.Name = path

}
