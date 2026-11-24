package toolchain

import (
	"embed"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"
)

//go:embed tools.yaml
var catalogYAML embed.FS

// Tool is one relocatable toolchain binary resolved for a host OS.
type Tool struct {
	Name        string
	Version     string
	Binary      string
	URL         string
	SHA256      string
	ArchiveRoot string
}

type toolFile struct {
	Name    string             `yaml:"name"`
	Version string             `yaml:"version"`
	Binary  string             `yaml:"binary"`
	OS      map[string]osAsset `yaml:"os"`
}

type osAsset struct {
	URL         string `yaml:"url"`
	SHA256      string `yaml:"sha256"`
	ArchiveRoot string `yaml:"archive_root"`
}

// List returns every catalog tool for this host, with {version} substituted.
// install all means this list.
func List() ([]Tool, error) {
	return list(runtime.GOOS, runtime.GOARCH)
}

// Lookup returns one catalog tool for this host.
func Lookup(name string) (Tool, error) {
	return lookup(name, runtime.GOOS, runtime.GOARCH)
}

// CacheDir is the versioned install directory under the user cache.
func CacheDir(tool Tool) (string, error) {
	base, err := os.UserCacheDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(base, "forge", "toolchains", tool.Name, tool.Version), nil
}

// BinaryPath is the executable inside CacheDir. Windows paths gain a .exe suffix.
func BinaryPath(tool Tool) (string, error) {
	dir, err := CacheDir(tool)
	if err != nil {
		return "", err
	}
	rel := tool.Binary
	if runtime.GOOS == "windows" && !strings.HasSuffix(strings.ToLower(rel), ".exe") {
		rel += ".exe"
	}
	return filepath.Join(dir, filepath.FromSlash(rel)), nil
}

func list(goos, goarch string) ([]Tool, error) {
	catalog, err := loadCatalog()
	if err != nil {
		return nil, err
	}
	names := make([]string, 0, len(catalog))
	for name := range catalog {
		names = append(names, name)
	}
	sort.Strings(names)

	tools := make([]Tool, 0, len(names))
	for _, name := range names {
		tool, err := resolve(catalog[name], name, goos, goarch)
		if err != nil {
			return nil, err
		}
		tools = append(tools, tool)
	}
	return tools, nil
}

func lookup(name, goos, goarch string) (Tool, error) {
	catalog, err := loadCatalog()
	if err != nil {
		return Tool{}, err
	}
	raw, ok := catalog[name]
	if !ok {
		return Tool{}, fmt.Errorf("unknown toolchain %q", name)
	}
	return resolve(raw, name, goos, goarch)
}

func loadCatalog() (map[string]toolFile, error) {
	data, err := catalogYAML.ReadFile("tools.yaml")
	if err != nil {
		return nil, err
	}
	var file map[string]toolFile
	if err := yaml.Unmarshal(data, &file); err != nil {
		return nil, err
	}
	return file, nil
}

func resolve(raw toolFile, id, goos, goarch string) (Tool, error) {
	if goarch != "amd64" || (goos != "linux" && goos != "windows") {
		return Tool{}, fmt.Errorf("unsupported toolchain host %s/%s", goos, goarch)
	}
	asset, ok := raw.OS[goos]
	if !ok {
		return Tool{}, fmt.Errorf("toolchain %q has no %s archive", id, goos)
	}
	name := raw.Name
	if name == "" {
		name = id
	}
	if raw.Version == "" || raw.Binary == "" || asset.URL == "" || asset.SHA256 == "" {
		return Tool{}, fmt.Errorf("toolchain %q is missing version, binary, url, or sha256", id)
	}
	version := raw.Version
	return Tool{
		Name:        name,
		Version:     version,
		Binary:      raw.Binary,
		URL:         strings.ReplaceAll(asset.URL, "{version}", version),
		SHA256:      asset.SHA256,
		ArchiveRoot: strings.ReplaceAll(asset.ArchiveRoot, "{version}", version),
	}, nil
}
