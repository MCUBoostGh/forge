package cmd

import (
	"fmt"
)

func Build() error{
	// Implement the build logic here
	_,err := fmt.Println("Building the project...")
	if err != nil {
		return err
	}
	err = runCommand("cmake", "--preset", "default")
	if err != nil {
		return err
	}
	err = runCommand("cmake", "--build", "--preset", "default")
	if err != nil {
		return err
	}
	return nil	
}