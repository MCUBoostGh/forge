package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

func Init(nameProject string) error {

	logError := log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime)
	logSuccess := log.New(os.Stderr, "Success: ", log.Ldate|log.Ltime)

	fmt.Println("Initializing a new project:", nameProject)

	_, err := os.Stat(nameProject)

	// If no error, the path already exists
	if err == nil {
		logError.Println("Project directory already exists.")
		return fmt.Errorf("project directory already exists")
	}
	// If we got an error that is not "not exist", return it
	if err != nil && !os.IsNotExist(err) {
		logError.Println("Failed to check project directory:", err)
		return err
	}

	err = os.MkdirAll(nameProject, 0755)
	if err != nil {
		logError.Println("Failed to create project directory.", err)
		return err
	}

	err = os.MkdirAll(filepath.Join(nameProject, "src"), 0755)
	if err != nil {
		logError.Println("Failed to create project directory.", err)
		return err
	}

	err = os.MkdirAll(filepath.Join(nameProject, "include"), 0755)
	if err != nil {
		logError.Println("Failed to create project directory.", err)
		return err
	}

	_, err = os.Create(filepath.Join(nameProject, "main.c"))
	if err != nil {
		logError.Println("Failed to create main.c file.", err)
		return err
	}

	_, err = os.Create(filepath.Join(nameProject, "CMakeLists.txt"))
	if err != nil {
		logError.Println("Failed to create CMakeLists.txt file.", err)
		return err
	}

	_, err = os.Create(filepath.Join(nameProject, "Forge.toml"))
	if err != nil {
		logError.Println("Failed to create Forge.toml file.", err)
		return err
	}

	_, err = os.Create(filepath.Join(nameProject, "CMakePresets.json"))
	if err != nil {
		logError.Println("Failed to create CMakePresets.json file.", err)
		return err
	}

	logSuccess.Println("Project directory created successfully at", nameProject)

	return nil
}
