package templates

import (
	"embed"
)

//go:embed *.tmpl host/*.tmpl
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
	// keep Target etc. if MCU tmpls need them later
}
