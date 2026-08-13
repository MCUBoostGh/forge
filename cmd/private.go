package cmd

import (
	"log"
	"os"
	"os/exec"
	"forge/contents"
)

var logError = log.New(os.Stderr, "ERROR: ", 0)

var listFolders = []string{"src", "include"}

var listFilesContentMap = map[string]string{
	"Forge.toml":     contents.ForgeTOMLContent,
	"main.c":         contents.MainCContent,
	"CMakeLists.txt": contents.CMakeListsContent,
	"CMakePresets.json": contents.CMakePresetsContent,
}

const buildDir = "build"

var projectDir string

func runCommand(command string, args ...string) error {
	cmd := exec.Command(command, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	return cmd.Run()
}

func setProjectDir(path string) {
	projectDir = path
}
