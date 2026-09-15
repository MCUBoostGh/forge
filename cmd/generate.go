package cmd

import (
	"bytes"
	"fmt"
	"forge/internal/config"
	"forge/internal/devices"
	"forge/internal/logger"
	"forge/internal/templates"
	"html/template"
	"os"
	"path/filepath"
	"strings"
)

func createNewProject(nameProject string) error {

	err := os.MkdirAll(nameProject, 0755)
	if err != nil {
		return err
	}

	return nil
}

func updateProjectData(cfg *config.Config) (projectData, error) {
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

func newTemplateData(cfg *config.Config) (templates.TemplateData, error) {
	minVer := cfg.CMake.MinimumRequiredVersion
	if minVer == "" {
		minVer = "3.20"
	}

	desc := cfg.Project.Name // or a real description field later
	cStd := "11"

	return templates.TemplateData{
		ProjectName:                 cfg.Project.Name,
		ProjectVersion:              cfg.Project.Version,
		ProjectDescription:          desc,
		CStandard:                   cStd,
		CMakeMinimumRequiredVersion: minVer,
		ToolchainCompiler:           cfg.Toolchain.Compiler,
		BuildType:                   cfg.Build.Type,
	}, nil
}

func New(args ...string) error {

	msgErr := fmt.Errorf("Falied to generate new project.")

	if len(args) < 3 {
		logger.Error("Device not specified. Use --device <device> to specify the target device.")
		return msgErr
	}

	nameProject := args[0]

	logger.Infof("Initializing a new project: %s", nameProject)
	
	_, err := os.Stat(nameProject)
	// If no error, the path already exists
	if err == nil {
		logger.Error(err)
		return msgErr
	}

	cfg := config.Get()
	cfg.Project.Name = nameProject
	config.Set(cfg)


	subArgs := args[1]

	switch subArgs {
	case "--device":
		devicePartNumber := args[2]
		catalog, err := devices.Lookup(devicePartNumber)
		if err != nil {
			logger.Error(err)
			return msgErr
		}
		configPtr := config.Get()
		configPtr.Target.Device = catalog.ID
		configPtr.Toolchain.Compiler = compilersMap[catalog.CPU]
		updateProjectData(&configPtr)
		config.Set(configPtr)

		err = createNewProject(nameProject)
		if err != nil {
			logger.Error(err)
			return msgErr
		}

		err=config.New(nameProject)
		if err!=nil{
			logger.Error(err)
			return msgErr
		}

	default:
		logger.Error("Invalid arguments")
		return msgErr
	}

	return nil
}

func Init() error {

	msgErr := fmt.Errorf("Failed to initializing project.")

	err := config.Read()
	if err != nil {
		logger.Error("Failed to load config file.")
		logger.Error(err)
		return msgErr
	}

	// Scaffold into the project cwd (where Forge.toml lives), not Project.Name.
	for _, nameFolder := range listFolders {
		err := os.MkdirAll(nameFolder, 0755)
		if err != nil {
			logger.Error("Failed to create folder:", nameFolder, err)
			return msgErr
		}
	}

	for outName, tmplPath := range listFilesContentMap {
		if err := os.MkdirAll(filepath.Dir(outName), 0755); err != nil && filepath.Dir(outName) != "." {
			logger.Error("Failed to create parent dir for ", outName)
			logger.Error(err)
			return msgErr
		}

		raw, err := templates.FS.ReadFile(tmplPath)
		if err != nil {
			logger.Error("Failed to load template files:", outName)
			logger.Error(err)
			return msgErr
		}

		tmpl, err := template.New(outName).Parse(string(raw))
		if err != nil {
			logger.Error("Failed to parse template file:", outName)
			logger.Error(err)
			return msgErr
		}
		cfg := config.Get()
		data, err := newTemplateData(&cfg)
		if err != nil {
			logger.Error("Failed to build template data")
			logger.Error(err)
			return msgErr
		}

		var buf bytes.Buffer
		err = tmpl.Execute(&buf, data)
		if err != nil {
			logger.Error("Failed to execute template file:", outName)
			logger.Error(err)
			return msgErr
		}

		err = os.WriteFile(outName, buf.Bytes(), 0644)
		if err != nil {
			logger.Error("Failed to write to ", outName, " file.", err)
			return err
		}
	}
	return nil
}
