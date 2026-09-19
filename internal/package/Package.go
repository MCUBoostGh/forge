package Package

import (
	"archive/tar"
	"compress/gzip"
	"embed"
	"fmt"
	"forge/internal/logger"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
)

type Package struct {
	Name           string   `yaml:"name"`
	Version        string   `yaml:"version"`
	URL            string   `yaml:"url"`
	TagPrefix      string   `yaml:"tag_prefix"`
	ArchiveRoot    string   `yaml:"archive_root"`
	Include        []string `yaml:"include"`
	Sources        []string `yaml:"sources"`
	ProjectSources []string `yaml:"project_sources"`
	StartupDir     string   `yaml:"startup_dir"`
	Depends        []string `yaml:"depends"`
	Size           int64    `yaml:"size"`
}
type progressReader struct {
	name  string
	r     io.Reader
	total int64
	read  int64
}

const progressBarWidth = 30

//go:embed *.yaml
var tp embed.FS

func contentLength(resp *http.Response, fallback int64) int64 {
	if resp.ContentLength > 0 {
		return resp.ContentLength
	}
	if raw := resp.Header.Get("Content-Length"); raw != "" {
		n, err := strconv.ParseInt(raw, 10, 64)
		if err == nil && n > 0 {
			return n
		}
	}
	if fallback > 0 {
		return fallback
	}
	return 0
}

func safeJoin(dest, name string) (string, error) {
	clean := filepath.Clean(name)
	target := filepath.Join(dest, clean)
	if !strings.HasPrefix(target, filepath.Clean(dest)+string(os.PathSeparator)) &&
		target != filepath.Clean(dest) {
		return "", os.ErrInvalid // block path traversal (zip slip)
	}
	return target, nil
}

func formatBytes(n int64) string {
	const (
		kb = 1024
		mb = 1024 * kb
		gb = 1024 * mb
	)
	switch {
	case n >= gb:
		return fmt.Sprintf("%.2f GB", float64(n)/float64(gb))
	case n >= mb:
		return fmt.Sprintf("%.2f MB", float64(n)/float64(mb))
	case n >= kb:
		return fmt.Sprintf("%.2f KB", float64(n)/float64(kb))
	default:
		return fmt.Sprintf("%d B", n)
	}
}

func (p *progressReader) render() {
	var bar string
	var extra string
	if p.total > 0 {
		pct := float64(p.read) / float64(p.total)
		if pct > 1 {
			pct = 1
		}
		filled := int(pct * float64(progressBarWidth))
		if filled > progressBarWidth {
			filled = progressBarWidth
		}
		if filled > 0 && filled < progressBarWidth {
			bar = strings.Repeat("=", filled-1) + ">" + strings.Repeat(" ", progressBarWidth-filled)
		} else {
			bar = strings.Repeat("=", filled) + strings.Repeat(" ", progressBarWidth-filled)
		}
		extra = fmt.Sprintf(" %3.0f%%  %s / %s", pct*100, formatBytes(p.read), formatBytes(p.total))
	} else {
		pos := int(p.read / (256 * 1024) % int64(progressBarWidth-4))
		bar = strings.Repeat(" ", pos) + "<==>" + strings.Repeat(" ", progressBarWidth-4-pos)
		extra = fmt.Sprintf("  %s", formatBytes(p.read))
	}
	line := fmt.Sprintf("\r%s [%s]%s", p.name, bar, extra)
	fmt.Printf("%s\033[K", line)
}

func (p *progressReader) done() {
	p.render()
	fmt.Println()
}

func extractTarGz(src, dest string) error {
	f, err := os.Open(src)
	if err != nil {
		return err
	}
	defer f.Close()

	gz, err := gzip.NewReader(f)
	if err != nil {
		return err
	}
	defer gz.Close()

	tr := tar.NewReader(gz)
	if err := os.MkdirAll(dest, 0755); err != nil {
		return err
	}

	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}

		target, err := safeJoin(dest, hdr.Name)
		if err != nil {
			return err
		}

		switch hdr.Typeflag {

		case tar.TypeDir:
			if err := os.MkdirAll(target, 0755); err != nil {
				return err
			}
		case tar.TypeReg:
			if err := os.MkdirAll(filepath.Dir(target), 0755); err != nil {
				return err
			}
			out, err := os.OpenFile(target, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, os.FileMode(hdr.Mode))
			if err != nil {
				return err
			}
			if _, err := io.Copy(out, tr); err != nil {
				out.Close()
				return err
			}
			out.Close()
		}
	}
}
func (p *progressReader) Read(b []byte) (int, error) {
	n, err := p.r.Read(b)
	p.read += int64(n)
	p.render()
	return n, err
}
func loadCatalog() (map[string]Package, error) {

	catalog := make(map[string]Package)

	entries, err := tp.ReadDir(".")
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".yaml") {
			continue
		}
		data, err := tp.ReadFile(entry.Name())
		if err != nil {
			return nil, err
		}
		var file map[string]Package
		if err := yaml.Unmarshal(data, &file); err != nil {
			return nil, err
		}
		for id, pkg := range file {
			if pkg.Name == "" {
				pkg.Name = id
			}
			catalog[id] = pkg
		}
	}
	return catalog, nil
}

func tagFor(pkg Package, version string) string {
	if pkg.TagPrefix == "" {
		return version
	}
	if strings.HasPrefix(version, pkg.TagPrefix) {
		return version
	}
	return pkg.TagPrefix + version
}

func Register(name string, version string) (Package, error) {
	msgErr := fmt.Errorf("Failed to register %s third_party.", name)

	catalog, err := loadCatalog()
	if err != nil {
		logger.Error(err)
		return Package{}, msgErr
	}

	pkg, ok := catalog[name]
	if !ok {
		logger.Error("unknown third-party package:", name)
		return Package{}, msgErr
	}

	pkg.Version = version
	tag := tagFor(pkg, version)
	// GitHub tag URLs use the tag (often vX.Y.Z); archive directories drop the leading v.
	pkg.URL = strings.ReplaceAll(pkg.URL, "{version}", tag)
	pkg.ArchiveRoot = strings.ReplaceAll(pkg.ArchiveRoot, "{version}", version)
	return pkg, nil
}

func ParseSpec(spec string) (name string, version string, err error) {
	parts := strings.SplitN(strings.TrimSpace(spec), "@", 2)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return "", "", fmt.Errorf("invalid dependency spec %q (want name@version)", spec)
	}
	return parts[0], parts[1], nil
}

func PackagesDir() (string, error) {
	base, err := os.UserCacheDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(base, "forge", "packages"), nil
}

func (pkg Package) packagesDir() (string, error) {
	return PackagesDir()
}

func (pkg Package) tmpArchivePath() string {
	return filepath.Join(os.TempDir(), "forge-"+pkg.Name+"-"+pkg.Version+".tar.gz")
}

func (pkg Package) extractedRoot() (string, error) {
	dir, err := pkg.packagesDir()
	if err != nil {
		return "", err
	}
	if pkg.ArchiveRoot != "" {
		return filepath.Join(dir, pkg.ArchiveRoot), nil
	}
	return filepath.Join(dir, pkg.Name+"-"+pkg.Version), nil
}

func (pkg Package) isExtracted() bool {
	dir, err := pkg.extractedRoot()
	if err != nil {
		return false
	}
	info, err := os.Stat(dir)
	return err == nil && info.IsDir()
}

func (pkg Package) IncludeDirs() []string {
	root, err := pkg.extractedRoot()
	if err != nil {
		return nil
	}
	dirs := make([]string, 0, len(pkg.Include))
	for _, inc := range pkg.Include {
		dirs = append(dirs, filepath.ToSlash(filepath.Join(root, inc)))
	}
	return dirs
}

func skipCachedSource(path string) bool {
	base := strings.ToLower(filepath.Base(path))
	if !strings.HasSuffix(base, ".c") {
		return true
	}
	if strings.Contains(base, "_template") {
		return true
	}
	if strings.Contains(base, "_ll_") {
		return true
	}
	return false
}

// SourceFiles expands yaml sources (file, directory, or glob) under the extracted cache root.
func (pkg Package) SourceFiles() ([]string, error) {
	if len(pkg.Sources) == 0 {
		return nil, nil
	}
	root, err := pkg.extractedRoot()
	if err != nil {
		return nil, err
	}

	seen := make(map[string]struct{})
	var files []string
	for _, spec := range pkg.Sources {
		pattern := filepath.Join(root, filepath.FromSlash(spec))
		info, err := os.Stat(pattern)
		var matches []string
		if err == nil && info.IsDir() {
			matches, err = filepath.Glob(filepath.Join(pattern, "*.c"))
		} else {
			matches, err = filepath.Glob(pattern)
		}
		if err != nil {
			return nil, err
		}
		for _, match := range matches {
			if skipCachedSource(match) {
				continue
			}
			slash := filepath.ToSlash(match)
			if _, ok := seen[slash]; ok {
				continue
			}
			seen[slash] = struct{}{}
			files = append(files, slash)
		}
	}
	sort.Strings(files)
	return files, nil
}

func (pkg Package) ConfTemplates() ([]string, error) {
	var templates []string
	for _, dir := range pkg.IncludeDirs() {
		matches, err := filepath.Glob(filepath.Join(dir, "*hal_conf_template.h"))
		if err != nil {
			return nil, err
		}
		templates = append(templates, matches...)
	}
	return templates, nil
}

func (pkg Package) StartupFileName(stm32Device string) string {
	if pkg.StartupDir == "" || stm32Device == "" {
		return ""
	}
	return "startup_" + strings.ToLower(stm32Device) + ".s"
}

func (pkg Package) ProjectFileNames(stm32Device string) []string {
	names := make([]string, 0, len(pkg.ProjectSources)+1)
	for _, src := range pkg.ProjectSources {
		names = append(names, filepath.Base(src))
	}
	if name := pkg.StartupFileName(stm32Device); name != "" {
		names = append(names, name)
	}
	return names
}

func copyFileIfMissing(src, dest string) error {
	if _, err := os.Stat(dest); err == nil {
		return nil
	}
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	out, err := os.OpenFile(dest, os.O_CREATE|os.O_WRONLY|os.O_EXCL, 0644)
	if err != nil {
		if os.IsExist(err) {
			return nil
		}
		return err
	}
	_, copyErr := io.Copy(out, in)
	closeErr := out.Close()
	if copyErr != nil {
		_ = os.Remove(dest)
		return copyErr
	}
	if closeErr != nil {
		return closeErr
	}
	logger.Infof("Copied %s", dest)
	return nil
}

func (pkg Package) CopyProjectFiles(destDir, stm32Device string) error {
	if len(pkg.ProjectSources) == 0 && pkg.StartupDir == "" {
		return nil
	}
	root, err := pkg.extractedRoot()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(destDir, 0755); err != nil {
		return err
	}
	for _, rel := range pkg.ProjectSources {
		src := filepath.Join(root, filepath.FromSlash(rel))
		dest := filepath.Join(destDir, filepath.Base(rel))
		if err := copyFileIfMissing(src, dest); err != nil {
			return fmt.Errorf("copy %s: %w", rel, err)
		}
	}
	if name := pkg.StartupFileName(stm32Device); name != "" {
		src := filepath.Join(root, filepath.FromSlash(pkg.StartupDir), name)
		dest := filepath.Join(destDir, name)
		if err := copyFileIfMissing(src, dest); err != nil {
			return fmt.Errorf("copy startup %s: %w", name, err)
		}
	}
	return nil
}

func (pkg Package) Download() error {
	msgErr := fmt.Errorf("Failed to download package %s(%s).", pkg.Name, pkg.Version)

	if pkg.isExtracted() {
		return nil
	}

	dest := pkg.tmpArchivePath()
	req, err := http.NewRequest(http.MethodGet, pkg.URL, nil)
	if err != nil {
		logger.Error(err)
		return msgErr
	}

	req.Header.Set("Accept-Encoding", "identity")

	transport := http.DefaultTransport.(*http.Transport).Clone()
	transport.DisableCompression = true
	client := &http.Client{Transport: transport}

	response, err := client.Do(req)
	if err != nil {
		logger.Error(err)
		return msgErr
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("download failed: %s", response.Status)
	}

	progress := &progressReader{
		name:  pkg.Name,
		r:     response.Body,
		total: contentLength(response, pkg.Size),
	}

	out, err := os.Create(dest)
	if err != nil {
		logger.Error(err)
		return msgErr
	}

	if _, err := io.Copy(out, progress); err != nil {
		out.Close()
		fmt.Println()
		logger.Error(err)
		_ = os.Remove(dest)
		return msgErr
	}
	out.Close()
	progress.done()
	return nil
}

func (pkg Package) Extract() error {
	msgErr := fmt.Errorf("Failed to extract package %s(%s).", pkg.Name, pkg.Version)

	dest, err := pkg.packagesDir()
	if err != nil {
		logger.Error(err)
		return msgErr
	}

	if pkg.isExtracted() {
		_ = os.Remove(pkg.tmpArchivePath())
		return nil
	}

	archive := pkg.tmpArchivePath()
	if _, err := os.Stat(archive); err != nil {
		logger.Error(err)
		return msgErr
	}

	if err := os.MkdirAll(dest, 0755); err != nil {
		logger.Error(err)
		return msgErr
	}

	if err := extractTarGz(archive, dest); err != nil {
		logger.Error(err)
		return msgErr
	}
	_ = os.Remove(archive)
	return nil
}

func (pkg Package) Build() error {
	return nil
}
