package cmd

import (
	"fmt"
	"forge/internal/logger"
)

func Build(args ...string) error {
	if len(args) < 1 {
		logger.Error("Preset not specified. Use forge build <preset> (debug or release).")
		return fmt.Errorf("Failed to build project.")
	}

	preset := args[0]
	logger.Infof("Building with CMake preset %s", preset)

	if err := runCommand("cmake", "--preset", preset); err != nil {
		return err
	}
	if err := runCommand("cmake", "--build", "--preset", preset); err != nil {
		return err
	}
	return nil
}
