package cmd

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"forge/internal/toolchain"
)

var listFolders = []string{"src", "include", "cmake"}

var listFilesContentMap = map[string]string{
	"main.c":                        "main.txt.tmpl",
	"CMakeLists.txt":                "cmake/CMakeLists.txt.tmpl",
	"CMakePresets.json":             "cmake/CMakePresets.json.tmpl",
	"cmake/gcc-arm-none-eabi.cmake": "cmake/gcc-arm-none-eabi.cmake.tmpl",
	"cmake/Package.cmake":           "cmake/Package.cmake.tmpl",
	"LinkerScript.ld":               "ld/LinkerScript.ld.tmpl",
}

var compilersMap = map[string]string{
	"cortex-m0": "gcc-arm-none-eabi",
	"cortex-m3": "gcc-arm-none-eabi",
	"cortex-m4": "gcc-arm-none-eabi",
	"cortex-m7": "gcc-arm-none-eabi",
	"x86_64":    "gcc",
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
	env, err := commandEnv()
	if err != nil {
		return err
	}
	// exec.Command resolves a bare name with the parent PATH and ignores cmd.Env.
	bin, err := toolchain.FindExecutable(command, pathValue(env))
	if err != nil {
		if command == "cmake" && errors.Is(err, exec.ErrNotFound) {
			return fmt.Errorf("cmake not found; run forge toolchain install cmake")
		}
		return err
	}
	cmd := exec.Command(bin, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Env = env
	return cmd.Run()
}

func pathValue(env []string) string {
	for _, entry := range env {
		if len(entry) >= 5 && strings.EqualFold(entry[:5], "PATH=") {
			return entry[5:]
		}
	}
	return ""
}

func commandEnv() ([]string, error) {
	dirs, err := toolchain.BinDirs()
	if err != nil {
		return nil, err
	}
	path := toolchain.PrependPath(os.Getenv("PATH"), dirs)
	env := os.Environ()
	replaced := false
	for i, entry := range env {
		if len(entry) >= 5 && strings.EqualFold(entry[:5], "PATH=") {
			env[i] = "PATH=" + path
			replaced = true
		}
	}
	if !replaced {
		env = append(env, "PATH="+path)
	}
	return env, nil
}
