package cmd

import (
	"fmt"
)

func Build() error {
	_, err := fmt.Println("Building the project...")
	if err != nil {
		return err
	}

	if err := loadConfig(forgeTOMLName); err != nil {
		return err
	}

	preset := config.Build.Type
	if preset == "" {
		preset = "debug"
	}

	err = runCommand("cmake", "--preset", preset)
	if err != nil {
		return err
	}
	err = runCommand("cmake", "--build", "--preset", preset)
	if err != nil {
		return err
	}
	return nil
}