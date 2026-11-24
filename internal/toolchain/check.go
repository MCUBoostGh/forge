package toolchain

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

const (
	StatusOK       = "ok"
	StatusMissing  = "missing"
	StatusMismatch = "version mismatch"
)

// CheckResult is one catalog tool after probing the cache and PATH.
type CheckResult struct {
	Tool     Tool
	Path     string
	Detected string
	Status   string
}

var upstreamVersion = regexp.MustCompile(`^\d+\.\d+\.\d+`)

// Check runs tool --version from the cache directory, then from PATH.
// An empty dir uses CacheDir. It does not read Forge.toml or change PATH.
func Check(tool Tool, dir string) (CheckResult, error) {
	result := CheckResult{Tool: tool, Status: StatusMissing}
	if dir == "" {
		var err error
		dir, err = CacheDir(tool)
		if err != nil {
			return result, err
		}
	}

	path := executablePath(dir, tool.Binary)
	if _, err := os.Stat(path); err != nil {
		found, lookErr := exec.LookPath(filepath.Base(tool.Binary))
		if lookErr != nil {
			return result, nil
		}
		path = found
	}

	out, err := commandOutput(path, "--version")
	result.Path = path
	result.Detected = firstLine(out)
	if result.Detected == "" && err != nil {
		result.Detected = strings.Join(strings.Fields(err.Error()), " ")
	}
	if versionMatches(tool.Version, out) {
		result.Status = StatusOK
	} else {
		result.Status = StatusMismatch
	}
	return result, nil
}

func versionMatches(catalog, output string) bool {
	if catalog != "" && strings.Contains(output, catalog) {
		return true
	}
	base := upstreamVersion.FindString(catalog)
	return base != "" && strings.Contains(output, base)
}

func commandOutput(name string, args ...string) (string, error) {
	cmd := exec.Command(name, args...)
	var buf bytes.Buffer
	cmd.Stdout = &buf
	cmd.Stderr = &buf
	err := cmd.Run()
	return buf.String(), err
}

func firstLine(s string) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return ""
	}
	if i := strings.IndexAny(s, "\r\n"); i >= 0 {
		s = s[:i]
	}
	return strings.Join(strings.Fields(s), " ")
}

// FormatReport renders check results as an aligned block per tool.
func FormatReport(results []CheckResult) string {
	nameW := 0
	verW := 0
	for _, r := range results {
		if n := len(r.Tool.Name); n > nameW {
			nameW = n
		}
		if n := len(r.Tool.Version); n > verW {
			verW = n
		}
	}
	var b strings.Builder
	for i, r := range results {
		if i > 0 {
			b.WriteByte('\n')
		}
		detected := r.Detected
		if detected == "" {
			detected = "-"
		}
		fmt.Fprintf(&b, "%-*s  %-*s  %s\n", nameW, r.Tool.Name, verW, r.Tool.Version, r.Status)
		fmt.Fprintf(&b, "  %s\n", detected)
		if r.Path != "" {
			fmt.Fprintf(&b, "  %s\n", r.Path)
		}
	}
	return b.String()
}

func (r CheckResult) String() string {
	return strings.TrimRight(FormatReport([]CheckResult{r}), "\n")
}
