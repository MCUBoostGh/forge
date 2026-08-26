package cmd

import (
	"errors"
	"fmt"
	"forge/internal/devices"
	"forge/internal/logger"
	"os"
	"path/filepath"
	"strings"

	"github.com/pelletier/go-toml/v2"
)

func updateProjectData(cfg *tomlConfig) (projectData, error) {
	pd := projectData{
		Name:      cfg.Project.Name,
		Version:   cfg.Project.Version,
		BuildType: cfg.Build.Type,
	}

	switch {
	case cfg.Target.Kind == "host-x64":
		pd.Architecture = "x86_64"
		pd.Compiler = compilersMap["x86_64"]
		pd.Target = nil
	case cfg.Target.Kind == "mcu" && strings.HasPrefix(strings.ToUpper(cfg.Target.Device), "STM32"):
		pd.Architecture = "cortexM"
		pd.Compiler = compilersMap["CortexM"]
		pd.Target = &targetData{
			Device: cfg.Target.Device,
			// CPU:    stm32CPUMap[strings.ToUpper(cfg.Target.Device)],
			Vendor: "STMicroelectronics",
			Family: "STM32",
			Series: strings.ToUpper(cfg.Target.Device),
		}
	default:
		return pd, fmt.Errorf("invalid target kind")
	}

	return pd, nil
}

func New(args ...string) error {

	msgErr := fmt.Errorf("Falied to generate new project.")

	if len(args) < 2 {
		logger.Error("Device not specified. Use --device <device> to specify the target device.")
		logger.Info("To list supported devices, run 'forge list'.")
		return msgErr
	}
	setConfigDefaults()
	nameProject := args[0]
	setConfig(nameProject)

	logger.Infof("Initializing a new project: %s", nameProject)

	_, err := os.Stat(nameProject)

	// If no error, the path already exists
	if err == nil {
		logger.Error("Project directory already exists.")
		return msgErr
	}

	// If we got an error that is not "not exist", return it
	if err != nil && !os.IsNotExist(err) {
		logger.Error(err)
		return msgErr
	}

	err = os.MkdirAll(nameProject, 0755)
	if err != nil {
		logger.Error(err)
		return msgErr
	}

	if args[1] == "--device" {
		if len(args) > 2 {
			device := args[2]
			catalog, err := devices.Lookup(device)
			if err != nil {
				logger.Error(err)
				return msgErr
			}

			config.Target.Device = catalog.ID
			config.Target.Kind = "mcu"
		}
	}
	_, err = os.Create(filepath.Join(config.Project.Name, forgeTOMLName))
	if err != nil {
		logger.Error(err)
		return msgErr
	}

	data, err := toml.Marshal(config)
	if err != nil {
		logger.Error(err)
		return msgErr
	}
	err = os.WriteFile(filepath.Join(config.Project.Name, forgeTOMLName), data, 0644)
	if err != nil {
		logger.Error(err)
		return msgErr
	}
	return nil
}

func Init() error {
	msgErr := fmt.Errorf("Failed to initializing project.")
	_, err := os.Stat(forgeTOMLName)
	if errors.Is(err, os.ErrNotExist) {
		logger.Error(err)
		return msgErr
	}
	for _, nameFolder := range listFolders {
		err := os.MkdirAll(filepath.Join(config.Project.Name, nameFolder), 0755)
		if err != nil {
			logger.Error("Failed to create folder:", nameFolder, err)
			return msgErr
		}
	}

	for name, Content := range listFilesContentMap {
		_, err := os.Create(filepath.Join(config.Project.Name, name))
		if err != nil {
			logger.Error("Failed to create ", name, " file.", err)
			return err
		}
		err = os.WriteFile(filepath.Join(config.Project.Name, name), []byte(Content), 0644)
		if err != nil {
			logger.Error("Failed to write to ", name, " file.", err)
			return err
		}
	}
	return nil
}
