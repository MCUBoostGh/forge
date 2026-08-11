package main


import(
	"os"
	"log"
	"forge/cmd"
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
}	

func main(){
	
	log.SetPrefix("")
	log.SetFlags(0)


	if len(os.Args) < 2 {
		help()
		log.Fatal("Usage: forge <command> [arguments]")
	}

	switch os.Args[1] {

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