package main

import (
	"forge/cmd"
	"forge/internal/logger"
	"log"
	"os"
)

func help() {
	logger.Println("Forge is a tool for managing your embedded c/c++ projects.")
	logger.Println("Usage: forge <command> [arguments]")
	logger.Println("Commands:")
	logger.Println("  init      Initialize a new project")
	logger.Println("  build     Build the project")
	logger.Println("  run       Run the project")
	logger.Println("  test      Run tests for the project")
	logger.Println("  help      Show this help message")
	logger.Println("  list 	 Show list of supported devices")
}

func main() {

	if len(os.Args) < 2 {
		help()
		logger.Fatal("Usage: forge <command> [arguments]")
	}

	switch os.Args[1] {
	case "new":
		if len(os.Args) < 3 {
			log.Fatal("usage: forge new <project> --device <device>")
		}
		args := os.Args[2:]
		err := cmd.New(args...)
		if err != nil {
			logger.Fatal(err)
		}
		log.Println("Use forge init to generate project structure and files.")

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
		logger.Println("Forge version 0.1.0")
	default:
		logger.Fatal("Unknown command. Use 'forge help' to see available commands.")
	}
}
