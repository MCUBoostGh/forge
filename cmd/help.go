package cmd

import "fmt"

func Help() {
	fmt.Print(`Forge is a tool for managing your embedded C/C++ projects.

Usage:
  forge <command> [arguments]

Commands:
  new       Create a new project: forge new <name> <device>
  init      Generate project structure and files
  build     Build the project: forge build <preset>
  run       Run the project
  test      Run tests for the project
  list      List supported devices
  help      Show this help message
  version   Print the version
`)
}
