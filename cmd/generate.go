package cmd

import (
	"bytes"
	"fmt"
	"forge/internal/config"
	"forge/internal/devices"
	"forge/internal/logger"
	thirdparty "forge/internal/package"
	"forge/internal/templates"
	"html/template"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func cmsisCoreHeader(cpu string) string {
	switch strings.ToLower(strings.TrimSpace(cpu)) {
	case "cortex-m0":
		return "core_cm0.h"
	case "cortex-m0plus", "cortex-m0+":
		return "core_cm0plus.h"
	case "cortex-m1":
		return "core_cm1.h"
	case "cortex-m3":
		return "core_cm3.h"
	case "cortex-m4":
		return "core_cm4.h"
	case "cortex-m7":
		return "core_cm7.h"
	default:
		return "core_cm3.h"
	}
}

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
	major, minor := parseCMakeVersion(minVer)

	debugBuildType := "Debug"
	releaseBuildType := "Release"
	if p, ok := catalog.Presets["debug"]; ok && p.BuildType != "" {
		debugBuildType = p.BuildType
	}
	if p, ok := catalog.Presets["release"]; ok && p.BuildType != "" {
		releaseBuildType = p.BuildType
	}

	var tmplData = templates.TemplateData{
		ProjectName:                 cfg.Project.Name,
		ProjectVersion:              cfg.Project.Version,
		ProjectDescription:          desc,
		CStandard:                   cStd,
		CMakeMinimumRequiredVersion: minVer,
		CMakeMinimumRequiredMajor:   major,
		CMakeMinimumRequiredMinor:   minor,
		ToolchainCompiler:           cfg.Toolchain.Compiler,
		DebugBuildType:              debugBuildType,
		ReleaseBuildType:            releaseBuildType,
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
	tmplData.STM32Device = catalog.STM32Device
	tmplData.CMSISCoreHeader = cmsisCoreHeader(catalog.CPU)

	cacheDir, err := thirdparty.PackagesDir()
	if err != nil {
		return templates.TemplateData{}, err
	}
	tmplData.CacheDir = filepath.ToSlash(cacheDir)

	for _, spec := range cfg.Dependencies {
		name, version, err := thirdparty.ParseSpec(spec)
		if err != nil {
			return templates.TemplateData{}, err
		}
		pkg, err := thirdparty.Register(name, version)
		if err != nil {
			return templates.TemplateData{}, err
		}
		sources, err := pkg.SourceFiles()
		if err != nil {
			return templates.TemplateData{}, err
		}
		if len(pkg.Sources) > 0 && len(sources) == 0 {
			return templates.TemplateData{}, fmt.Errorf("no C sources found for package %s in cache", pkg.Name)
		}
		kind := "INTERFACE"
		if len(sources) > 0 {
			kind = "STATIC"
		}
		var defines []string
		if catalog.STM32Device != "" && (kind == "STATIC" || len(pkg.ProjectSources) > 0 || pkg.StartupDir != "") {
			defines = []string{"USE_HAL_DRIVER", catalog.STM32Device}
		}
		tmplData.Packages = append(tmplData.Packages, templates.PackageData{
			Name:        pkg.Name,
			Kind:        kind,
			IncludeDirs: relCachePaths(cacheDir, pkg.IncludeDirs()),
			Sources:     relCachePaths(cacheDir, sources),
			Defines:     defines,
			Depends:     pkg.Depends,
		})
		tmplData.AppSources = append(tmplData.AppSources, pkg.ProjectFileNames(catalog.STM32Device)...)
	}

	return tmplData, nil
}

func parseCMakeVersion(v string) (major, minor int) {
	major, minor = 3, 20
	parts := strings.SplitN(strings.TrimSpace(v), ".", 3)
	if len(parts) >= 1 {
		if n, err := strconv.Atoi(parts[0]); err == nil && n > 0 {
			major = n
		}
	}
	if len(parts) >= 2 {
		if n, err := strconv.Atoi(parts[1]); err == nil && n >= 0 {
			minor = n
		}
	}
	return major, minor
}

func relCachePaths(cacheDir string, paths []string) []string {
	out := make([]string, 0, len(paths))
	for _, p := range paths {
		rel, err := filepath.Rel(cacheDir, p)
		if err != nil {
			out = append(out, filepath.ToSlash(p))
			continue
		}
		out = append(out, filepath.ToSlash(rel))
	}
	return out
}

func downloadDependencies(cfg *config.Config) error {

	logger.Info("Start downloading dependenceis ...")
	for _, spec := range cfg.Dependencies {

		name, version, err := thirdparty.ParseSpec(spec)
		if err != nil {
			return err
		}
		pkg, err := thirdparty.Register(name, version)
		if err != nil {
			return err
		}

		if err := pkg.Download(); err != nil {
			return err
		}
		if err := pkg.Extract(); err != nil {
			return err
		}
	}
	return nil
}

func copyHalConf(cfg *config.Config) error {
	for _, spec := range cfg.Dependencies {
		name, version, err := thirdparty.ParseSpec(spec)
		if err != nil {
			return err
		}
		pkg, err := thirdparty.Register(name, version)
		if err != nil {
			return err
		}
		confs, err := pkg.ConfTemplates()
		if err != nil {
			return err
		}
		for _, src := range confs {
			destName := strings.Replace(filepath.Base(src), "_template", "", 1)
			dest := filepath.Join("include", destName)
			if _, err := os.Stat(dest); err == nil {
				continue
			}
			data, err := os.ReadFile(src)
			if err != nil {
				return err
			}
			if err := os.MkdirAll("include", 0755); err != nil {
				return err
			}
			if err := os.WriteFile(dest, data, 0644); err != nil {
				return err
			}
			logger.Infof("Copied %s", dest)
		}
	}
	return nil
}

func copyProjectSources(cfg *config.Config) error {
	catalog, err := devices.Lookup(cfg.Target.Device)
	if err != nil {
		return err
	}
	for _, spec := range cfg.Dependencies {
		name, version, err := thirdparty.ParseSpec(spec)
		if err != nil {
			return err
		}
		pkg, err := thirdparty.Register(name, version)
		if err != nil {
			return err
		}
		if err := pkg.CopyProjectFiles(".", catalog.STM32Device); err != nil {
			return err
		}
	}
	return nil
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

	if err := downloadDependencies(&cfg); err != nil {
		logger.Error(err)
		return msgErr
	}

	if err := copyHalConf(&cfg); err != nil {
		logger.Error(err)
		return msgErr
	}

	if err := copyProjectSources(&cfg); err != nil {
		logger.Error(err)
		return msgErr
	}

	if err := generateFiles(&cfg); err != nil {
		logger.Error(err)
		return msgErr
	}

	return nil

}
