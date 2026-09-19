package main

import (
	"fmt"
	"forge/cmd"
	"forge/internal/logger"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		cmd.Help()
		os.Exit(1)
	}

	switch os.Args[1] {

	case "new":
		if len(os.Args) < 3 {
			logger.Error("Invalid arguments")
			logger.Fatal("Failed to generate new project.")
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
		err := cmd.Build(os.Args[2:]...)
		if err != nil {
			logger.Fatal(err)
		}
	case "help":
		cmd.Help()
	case "version":
		fmt.Println("Forge version 0.2.0")
	default:
		logger.Fatal("Unknown command. Use 'forge help' to see available commands.")
	}
}
