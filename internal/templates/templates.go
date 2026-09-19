package templates

import (
	"embed"
)

//go:embed *.tmpl cmake/*.tmpl ld/*.tmpl
var FS embed.FS

type TemplateData struct {
	ProjectName                 string
	ProjectVersion              string
	ProjectDescription          string
	CStandard                   string
	CMakeMinimumRequiredVersion string
	CMakeMinimumRequiredMajor   int
	CMakeMinimumRequiredMinor   int
	ToolchainCompiler           string
	BuildType                   string
	TargetDevice                string
	TargetCPU                   string
	TargetFPU                   string
	TargetFloatABI              string
	TargetVendor                string
	TargetFamily                string
	TargetSeries                string
	TargetFlashKB               int
	TargetRAMKB                 int
	CMSISCoreHeader             string
	Packages                    []PackageData
	// keep Target etc. if MCU tmpls need them later
}

type PackageData struct {
	Name        string
	IncludeDirs []string
}
