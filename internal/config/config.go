package config

import (
	"errors"
	"fmt"
	"forge/internal/logger"
	"os"
	"path/filepath"

	"github.com/pelletier/go-toml/v2"
)

type Config struct {
	Project struct {
		Name    string `toml:"name"`
		Version string `toml:"version"`
	} `toml:"project"`

	Target struct {
		Kind   string `toml:"kind"`   // mcu, fpga, etc.
		Device string `toml:"device"` // stm32f407vg, etc.
		Board  string `toml:"board"`  // stm32f407vg, etc.
	} `toml:"target"`

	Build struct {
		System string `toml:"system"` //cmake
		Type   string `toml:"type"`   // debug, release, relwithdebinfo, minsizerel
	} `toml:"build"`

	Toolchain struct {
		Compiler string `toml:"compiler"` // gcc, clang, etc.
	} `toml:"toolchain"`

	Dependencies []string `toml:"dependencies"`

	CMake struct {
		Version                string `toml:"version"`
		MinimumRequiredVersion string `toml:"minimum_required_version"`
	} `toml:"cmake"`
}

var forgeTOMLName string = "Forge.toml"
var config = Config{}

func setConfigDefaults() {
	config.Project.Name = "MyProject"
	config.Project.Version = "0.1.0"
	config.Build.System = "cmake"
	config.Build.Type = "debug"
	config.Toolchain.Compiler = "gcc"
	config.Dependencies = []string{"cmsis5@5.9.0"}
	config.CMake.Version = "3.30"
	config.CMake.MinimumRequiredVersion = "3.20"
}

func New(path string) error {

	msgErr := fmt.Errorf("Failed to generate config file.")

	_, err := os.Create(filepath.Join(path, forgeTOMLName))
	if err != nil {
		logger.Error(err)
		return msgErr
	}

	setConfigDefaults()
	config.Project.Name = path

	data, err := toml.Marshal(config)
	if err != nil {
		logger.Error(err)
		return msgErr
	}

	err = os.WriteFile(filepath.Join(path, forgeTOMLName), data, 0644)
	if err != nil {
		return err
	}
	return nil
}

func Read() error {

	msgErr := fmt.Errorf("Failed to read Config file.")

	if err := IsExist(); err != nil {
		logger.Error(err)
		return msgErr
	}

	raw, err := os.ReadFile(forgeTOMLName)
	if err != nil {
		logger.Error(err)
		return msgErr
	}

	if err := toml.Unmarshal(raw, &config); err != nil {
		return err
	}

	return nil
}

func Write() error {
	msgErr := fmt.Errorf("Failed to write Config file.")

	data, err := toml.Marshal(config)
	if err != nil {
		logger.Error(err)
		return msgErr
	}

	if err := os.WriteFile(forgeTOMLName, data, 0644); err != nil {
		logger.Error(err)
		return msgErr
	}
	return nil
}

func Get() Config {
	return config
}

func Set(cfg Config) {
	config = cfg
}

func IsExist() error {

	_, err := os.Stat(forgeTOMLName)
	if errors.Is(err, os.ErrNotExist) {
		return err
	}
	return nil
}
