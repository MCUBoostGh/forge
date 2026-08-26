package main

import (
	"forge/cmd"
	"log"
	"os"
)

func help() {
	log.Println("Forge is a tool for managing your embedded c/c++ projects.")
	log.Println("Usage: forge <command> [arguments]")
	log.Println("Commands:")
	log.Println("  init      Initialize a new project")
	log.Println("  build     Build the project")
	log.Println("  run       Run the project")
	log.Println("  test      Run tests for the project")
	log.Println("  help      Show this help message")
	log.Println("  list 	 Show list of supported devices")
}

func main() {

	log.SetPrefix("")
	log.SetFlags(0)

	if len(os.Args) < 2 {
		help()
		log.Fatal("Usage: forge <command> [arguments]")
	}

	switch os.Args[1] {
	case "new":
		if len(os.Args) < 3 {
			log.Fatalln("usage: forge new <project> --device <device>")
		}
		args := os.Args[2:]
		err := cmd.New(args...)
		if err != nil {
			log.Println(err)
			log.Fatalln("Failed to generate new project.")
		}
		log.Println("Use forge init to generate project structure and files.")

	case "init":
		err := cmd.Init()
		if err != nil {
			log.Fatal(err)
		}
	case "build":
		err := cmd.Build()
		if err != nil {
			log.Fatal(err)
		}
	case "help":
		help()
	case "version":
		log.Println("Forge version 0.1.0")
	default:
		log.Fatal("Unknown command. Use 'forge help' to see available commands.")
	}
}
