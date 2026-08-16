package cmd

import (
	"forge/contents"
	"log"
	"os"
	"os/exec"
)

var logError = log.New(os.Stderr, "ERROR: ", 0)

var listFolders = []string{"src", "include"}

var listFilesContentMap = map[string]string{
	"main.c":            contents.MainCContent,
	"CMakeLists.txt":    contents.CMakeListsContent,
	"CMakePresets.json": contents.CMakePresetsContent,
}

const buildDir = "build"

type Config struct {
	Project struct {
		Name    string `toml:"name"`
		Version string `toml:"version"`
	} `toml:"project"`

	Target struct {
		Device string `toml:"device"`
	} `toml:"target"`

	Build struct {
		System string `toml:"system"` //cmake
		Type   string `toml:"type"` // debug, release, relwithdebinfo, minsizerel
	} `toml:"build"`

	Toolchain struct {
		Compiler string `toml:"compiler"` // gcc, clang, etc.
	} `toml:"toolchain"`

	Dependencies []string `toml:"dependencies"`

	CMake struct {
		Version string `toml:"version"`
	} `toml:"cmake"`
}

var config = Config{}

func setConfigDefaults() {
	config.Project.Name = "MyProject"
	config.Project.Version = "0.1.0"
	config.Build.System = "cmake"
	config.Build.Type = "debug"
	config.Toolchain.Compiler = "gcc"
	config.Dependencies = []string{}
	config.CMake.Version = "4.2.0"
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
