package cmd

import (
	"fmt"
	"strings"

	"forge/internal/toolchain"
)

// Toolchain installs catalog binaries. It does not read Forge.toml.
func Toolchain(args ...string) error {
	dryRun := false
	var pos []string
	for _, arg := range args {
		switch {
		case arg == "--dry-run":
			dryRun = true
		case strings.HasPrefix(arg, "-"):
			return fmt.Errorf("unknown flag %q", arg)
		default:
			pos = append(pos, arg)
		}
	}
	if len(pos) < 1 {
		return fmt.Errorf("usage: forge toolchain install <name|all>")
	}
	if pos[0] != "install" {
		return fmt.Errorf("unknown toolchain command %q", pos[0])
	}
	if len(pos) != 2 {
		return fmt.Errorf("usage: forge toolchain install <name|all>")
	}

	opts := toolchain.InstallOptions{DryRun: dryRun}
	if pos[1] == "all" {
		tools, err := toolchain.List()
		if err != nil {
			return err
		}
		for _, tool := range tools {
			if err := toolchain.Install(tool, opts); err != nil {
				return err
			}
		}
		return nil
	}

	tool, err := toolchain.Lookup(pos[1])
	if err != nil {
		return err
	}
	return toolchain.Install(tool, opts)
}
