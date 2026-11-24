package toolchain

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

// BinDirs returns the bin directory of each catalog tool that is already installed.
func BinDirs() ([]string, error) {
	tools, err := List()
	if err != nil {
		return nil, err
	}
	var dirs []string
	seen := map[string]bool{}
	for _, tool := range tools {
		bin, err := BinaryPath(tool)
		if err != nil {
			return nil, err
		}
		if _, err := os.Stat(bin); err != nil {
			continue
		}
		dir := filepath.Dir(bin)
		if seen[dir] {
			continue
		}
		seen[dir] = true
		dirs = append(dirs, dir)
	}
	return dirs, nil
}

// PrependPath puts dirs at the front of a PATH string.
func PrependPath(path string, dirs []string) string {
	if len(dirs) == 0 {
		return path
	}
	prefix := strings.Join(dirs, string(os.PathListSeparator))
	if path == "" {
		return prefix
	}
	return prefix + string(os.PathListSeparator) + path
}

// FindExecutable searches pathEnv for name. On Windows it also tries name.exe.
// An absolute path skips the search.
func FindExecutable(name, pathEnv string) (string, error) {
	if strings.TrimSpace(name) == "" {
		return "", exec.ErrNotFound
	}
	if filepath.IsAbs(name) || strings.ContainsAny(name, `/\`) {
		if info, err := os.Stat(name); err == nil && !info.IsDir() {
			return name, nil
		}
		return "", exec.ErrNotFound
	}

	candidates := []string{name}
	if runtime.GOOS == "windows" && !strings.EqualFold(filepath.Ext(name), ".exe") {
		candidates = append(candidates, name+".exe")
	}
	for _, dir := range filepath.SplitList(pathEnv) {
		if dir == "" {
			continue
		}
		for _, cand := range candidates {
			path := filepath.Join(dir, cand)
			info, err := os.Stat(path)
			if err != nil || info.IsDir() {
				continue
			}
			return path, nil
		}
	}
	return "", exec.ErrNotFound
}
