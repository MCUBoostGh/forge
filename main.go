package main

import (
	"fmt"
	"forge/cmd"
	"forge/internal/logger"
	"os"
)

func main() {
	code, err := run(os.Args[1:])
	if err != nil {
		logger.Fatal(err)
	}
	os.Exit(code)
}

func run(args []string) (int, error) {
	if len(args) < 1 {
		cmd.Help()
		return 1, nil
	}

	switch args[0] {

	case "new":
		if len(args) < 2 {
			logger.Error("Invalid arguments")
			return 1, fmt.Errorf("Failed to generate new project.")
		}
		if err := cmd.New(args[1:]...); err != nil {
			return 1, err
		}
		logger.Info("Use forge init to generate project structure and files.")
		return 0, nil

	case "init":
		if err := cmd.Init(); err != nil {
			return 1, err
		}
		return 0, nil

	case "build":
		if err := cmd.Build(args[1:]...); err != nil {
			return 1, err
		}
		return 0, nil

	case "help":
		cmd.Help()
		return 0, nil

	case "version":
		fmt.Println("Forge version 0.2.0")
		return 0, nil

	default:
		return 1, fmt.Errorf("Unknown command. Use 'forge help' to see available commands.")
	}
}
