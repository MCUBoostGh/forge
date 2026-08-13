package cmd

import (
	"log"
	"os"
	"os/exec"
)

var logError = log.New(os.Stderr, "ERROR: ", 0)

var listFolders = []string {"src","include"}
var listFiles = []string {"Forge.toml","main.c","CMakeLists.txt","CMakePresets.json"}

const buildDir = "build"

const forgeTOMLContent = `
[project]
name = "MyProject"
version = "0.1.0"

[build]
system = "cmake"
build_dir = buildDir

`



func runCommand(command string, args ...string) error {
	cmd := exec.Command(command, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	return cmd.Run()
}

var projectDir string

func setProjectDir(path string) {
	projectDir = path
}




