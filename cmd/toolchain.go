package cmd

import (
	"fmt"
	"strings"

	"forge/internal/toolchain"
)

// Toolchain installs and checks catalog binaries. It does not read Forge.toml.
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
		return toolchainUsage()
	}
	switch pos[0] {
	case "install":
		return installToolchain(pos[1:], dryRun)
	case "check":
		return checkToolchain(pos[1:])
	default:
		return fmt.Errorf("unknown toolchain command %q", pos[0])
	}
}

func installToolchain(args []string, dryRun bool) error {
	if len(args) != 1 {
		return toolchainUsage()
	}
	opts := toolchain.InstallOptions{DryRun: dryRun}
	if args[0] == "all" {
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
	tool, err := toolchain.Lookup(args[0])
	if err != nil {
		return err
	}
	return toolchain.Install(tool, opts)
}

func checkToolchain(args []string) error {
	if len(args) > 1 {
		return toolchainUsage()
	}
	var tools []toolchain.Tool
	if len(args) == 1 {
		tool, err := toolchain.Lookup(args[0])
		if err != nil {
			return err
		}
		tools = []toolchain.Tool{tool}
	} else {
		var err error
		tools, err = toolchain.List()
		if err != nil {
			return err
		}
	}

	var failed []string
	results := make([]toolchain.CheckResult, 0, len(tools))
	for _, tool := range tools {
		result, err := toolchain.Check(tool, "")
		if err != nil {
			return err
		}
		results = append(results, result)
		if result.Status != toolchain.StatusOK {
			failed = append(failed, tool.Name+" "+result.Status)
		}
	}
	fmt.Print(toolchain.FormatReport(results))
	if len(results) > 0 {
		fmt.Println()
	}
	if len(failed) > 0 {
		return fmt.Errorf("toolchain check: %s", strings.Join(failed, ", "))
	}
	return nil
}

func toolchainUsage() error {
	return fmt.Errorf("usage: forge toolchain install <name|all> | forge toolchain check [name]")
}
