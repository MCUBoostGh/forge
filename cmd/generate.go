package cmd

import (
	"fmt"
	"forge/contents"
	"os"
	"path/filepath"
	"strings"
)

func New(nameProject string) error {

	setConfigDefaults()
	setConfig(nameProject)

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

	forgeTOMLName := "Forge.toml"
	forgeTOMLContents := strings.ReplaceAll(contents.ForgeTOMLContent, "MyProject", config.Project.Name)

	_, err = os.Create(filepath.Join(config.Project.Name, forgeTOMLName))
	if err != nil {
		logError.Println("Failed to create ", forgeTOMLName, " file.", err)
		return err
	}
	err = os.WriteFile(filepath.Join(config.Project.Name, forgeTOMLName), []byte(forgeTOMLContents), 0644)
	if err != nil {
		logError.Println("Failed to write to ", forgeTOMLName, " file.", err)
		return err
	}
	return nil
}

func Init() error {

	for _, nameFolder := range listFolders {
		err := os.MkdirAll(filepath.Join(config.Project.Name, nameFolder), 0755)
		if err != nil {
			logError.Println("Failed to create folder:", nameFolder, err)
			return err
		}
	}

	for name, Content := range listFilesContentMap {
		_, err := os.Create(filepath.Join(config.Project.Name, name))
		if err != nil {
			logError.Println("Failed to create ", name, " file.", err)
			return err
		}
		err = os.WriteFile(filepath.Join(config.Project.Name, name), []byte(Content), 0644)
		if err != nil {
			logError.Println("Failed to write to ", name, " file.", err)
			return err
		}
	}
	return nil
}
