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

func syncTemplateData(cfg *config.Config) (templates.TemplateData, error) {

	catalog, err := devices.Lookup(cfg.Target.Device)
	if err != nil {
		logger.Error(err)
		return templates.TemplateData{}, fmt.Errorf("Failed to load device catalog.")
	}

	minVer := cfg.CMake.MinimumRequiredVersion
	if minVer == "" {
		minVer = "3.20"
	}

	desc := cfg.Project.Name // or a real description field later
	cStd := "11"

	var tmplData = templates.TemplateData{
		ProjectName:                 cfg.Project.Name,
		ProjectVersion:              cfg.Project.Version,
		ProjectDescription:          desc,
		CStandard:                   cStd,
		CMakeMinimumRequiredVersion: minVer,
		ToolchainCompiler:           cfg.Toolchain.Compiler,
		BuildType:                   cfg.Build.Type,
	}

	tmplData.TargetDevice = catalog.ID
	tmplData.TargetVendor = catalog.Vendor
	tmplData.TargetFamily = catalog.Family
	tmplData.TargetSeries = strings.ToUpper(catalog.ID)
	tmplData.TargetCPU = catalog.CPU
	tmplData.TargetFPU = catalog.FPU
	tmplData.TargetSeries = catalog.Series
	tmplData.TargetFlashKB = catalog.FlashKB
	tmplData.TargetRAMKB = catalog.RAMKB
	tmplData.TargetFloatABI = catalog.FloatABI

	return tmplData, nil
}

func generateFiles(cfg *config.Config) error {

	msgErr := fmt.Errorf("Falied to generate project files.")

	for outName, tmplPath := range listFilesContentMap {
		if err := os.MkdirAll(filepath.Dir(outName), 0755); err != nil && filepath.Dir(outName) != "." {
			logger.Error(err)
			return msgErr
		}

		raw, err := templates.FS.ReadFile(tmplPath)
		if err != nil {
			logger.Error(err)
			return msgErr
		}

		tmpl, err := template.New(outName).Parse(string(raw))
		if err != nil {
			logger.Error(err)
			return msgErr
		}

		data, err := syncTemplateData(cfg)
		if err != nil {
			logger.Error(err)
			return msgErr
		}

		var buf bytes.Buffer
		err = tmpl.Execute(&buf, data)
		if err != nil {
			logger.Error(err)
			return msgErr
		}

		err = os.WriteFile(outName, buf.Bytes(), 0644)
		if err != nil {
			logger.Error(err)
			return msgErr
		}
	}
	return nil
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

		cfg.Target.Device = catalog.ID
		cfg.Toolchain.Compiler = compilersMap[catalog.CPU]
		cfg.Target.Kind = "mcu"
		config.Set(cfg)

		err = createNewProject(nameProject)
		if err != nil {
			logger.Error(err)
			return msgErr
		}

		err = config.New(nameProject)
		if err != nil {
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

	if err := config.Read(); err != nil {
		logger.Error("Failed to load config file.")
		logger.Error(err)
		return msgErr
	}
	cfg := config.Get()

	// Scaffold into the project cwd (where Forge.toml lives), not Project.Name.
	for _, nameFolder := range listFolders {
		if err := os.MkdirAll(nameFolder, 0755); err != nil {
			logger.Error(err)
			return msgErr
		}
	}

	if err := generateFiles(&cfg); err != nil {
		logger.Error(err)
		return msgErr
	}

	dep, ok := registry["cmsis"]
	if !ok {
		logger.Error("Failed to load registery")
		return msgErr
	}

	logger.Info("Start downloading dependency.")

	cacheDir, err := download(dep)
	if err != nil {
		logger.Error(err)
		return msgErr
	}

	logger.Success(dep.Name + " downaloed in to " + cacheDir + " successfully.")

	cfg = config.Get()
	cfg.Dependencies = append(cfg.Dependencies, dep.Name+"@"+dep.Version)
	config.Set(cfg)
	return config.Write()

}
