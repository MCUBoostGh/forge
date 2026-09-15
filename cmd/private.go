package cmd

import (
	"os"
	"os/exec"
)

var listFolders = []string{"src", "include", "cmake"}


var listFilesContentMap = map[string]string{
	"main.c":                        "main.txt.tmpl",
	"CMakeLists.txt":                "host/CMakeLists.txt.tmpl",
	"CMakePresets.json":             "host/CMakePresets.json.tmpl",
	// "cmake/gcc-arm-none-eabi.cmake": contents.GccArmNoneEabiCmakeContent,
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



func runCommand(command string, args ...string) error {
	cmd := exec.Command(command, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	return cmd.Run()
}

