package main

import (
	"fmt"
	"os"

	"forge/cmd"
	"forge/internal/logger"
)

func help() {
	fmt.Print(`Forge is a tool for managing your embedded C/C++ projects.

Usage:
  forge <command> [arguments]

Commands:
  new       Create a new project: forge new <name> --device <device>
  init      Generate project structure and files
  build     Build the project
  run       Run the project
  test      Run tests for the project
  list      List supported devices
  help      Show this help message
  version   Print the version
`)
}

func main() {
	if len(os.Args) < 2 {
		help()
		os.Exit(1)
	}

	switch os.Args[1] {
	case "new":
		if len(os.Args) < 4 {
			logger.Error("Invalid argument")
			help()
			os.Exit(1)
		}
		args := os.Args[2:]
		err := cmd.New(args...)
		if err != nil {
			logger.Fatal(err)
		}
		logger.Info("Use forge init to generate project structure and files.")

	case "init":
		err := cmd.Init()
		if err != nil {
			logger.Fatal(err)
		}
	case "build":
		err := cmd.Build()
		if err != nil {
			logger.Fatal(err)
		}
	case "help":
		help()
	case "version":
		logger.Info("Forge version 0.1.0")
	default:
		logger.Fatal("Unknown command. Use 'forge help' to see available commands.")
	}
}
