package cmd

import (
	"fmt"
	"forge/internal/config"
)

func Build() error {
	_, err := fmt.Println("Building the project...")
	if err != nil {
		return err
	}

	

	preset := config.Get().Build.Type
	
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