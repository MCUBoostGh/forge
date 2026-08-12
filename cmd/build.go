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
	err = runCommand("cmake", "-S", ".", "-B", "build")
	if err != nil {
		return err
	}
	err = runCommand("cmake", "--build", "build")
	if err != nil {
		return err
	}
	return nil	
}