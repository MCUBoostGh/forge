package toolchain

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// InstallOptions controls where and whether Install writes.
// An empty Dir uses CacheDir.
type InstallOptions struct {
	DryRun bool
	Dir    string
}

// Install downloads one catalog tool into the user cache.
// It does not read Forge.toml or change PATH.
func Install(tool Tool, opts InstallOptions) error {
	dest, err := installDir(tool, opts)
	if err != nil {
		return err
	}
	if opts.DryRun {
		fmt.Printf("%s %s\n%s\n%s\n", tool.Name, tool.Version, tool.URL, dest)
		return nil
	}
	bin := executablePath(dest, tool.Binary)
	if _, err := os.Stat(bin); err == nil {
		fmt.Printf("%s %s already installed\n", tool.Name, tool.Version)
		return nil
	}

	parent := filepath.Dir(dest)
	if err := os.MkdirAll(parent, 0755); err != nil {
		return err
	}
	suffix, err := archiveSuffix(tool.URL)
	if err != nil {
		return err
	}
	partial := filepath.Join(parent, tool.Name+".partial"+suffix)
	if err := download(tool, partial); err != nil {
		os.Remove(partial)
		return err
	}
	defer os.Remove(partial)

	if err := verifySHA256(partial, tool.SHA256); err != nil {
		return fmt.Errorf("%s: %w", tool.Name, err)
	}

	tmp, err := os.MkdirTemp(parent, tool.Name+".extract-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tmp)

	if err := Extract(partial, tmp); err != nil {
		return fmt.Errorf("extract %s: %w", tool.Name, err)
	}
	root := filepath.Join(tmp, filepath.FromSlash(tool.ArchiveRoot))
	info, err := os.Stat(root)
	if err != nil || !info.IsDir() {
		return fmt.Errorf("archive root %q not found in %s", tool.ArchiveRoot, tool.Name)
	}
	if _, err := os.Stat(dest); err == nil {
		if err := os.RemoveAll(dest); err != nil {
			return err
		}
	}
	if err := os.Rename(root, dest); err != nil {
		return fmt.Errorf("install %s: %w", tool.Name, err)
	}
	return nil
}

func installDir(tool Tool, opts InstallOptions) (string, error) {
	if opts.Dir != "" {
		return opts.Dir, nil
	}
	return CacheDir(tool)
}

func archiveSuffix(rawURL string) (string, error) {
	u := rawURL
	if i := strings.IndexByte(u, '?'); i >= 0 {
		u = u[:i]
	}
	switch {
	case strings.HasSuffix(u, ".tar.gz"), strings.HasSuffix(u, ".tgz"):
		return ".tar.gz", nil
	case strings.HasSuffix(u, ".zip"):
		return ".zip", nil
	default:
		return "", fmt.Errorf("unsupported archive %q", rawURL)
	}
}

func download(tool Tool, dest string) error {
	req, err := http.NewRequest(http.MethodGet, tool.URL, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Accept-Encoding", "identity")

	transport := http.DefaultTransport.(*http.Transport).Clone()
	transport.DisableCompression = true
	client := &http.Client{Transport: transport}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("download %s: %s", tool.Name, resp.Status)
	}

	out, err := os.Create(dest)
	if err != nil {
		return err
	}
	progress := &progressReader{name: tool.Name, r: resp.Body, total: resp.ContentLength}
	_, copyErr := io.Copy(out, progress)
	closeErr := out.Close()
	progress.done()
	if copyErr != nil {
		return copyErr
	}
	return closeErr
}

func verifySHA256(path, want string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()
	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return err
	}
	got := hex.EncodeToString(h.Sum(nil))
	if !strings.EqualFold(got, want) {
		return fmt.Errorf("sha256 mismatch: got %s want %s", got, want)
	}
	return nil
}

type progressReader struct {
	name  string
	r     io.Reader
	total int64
	read  int64
}

func (p *progressReader) Read(b []byte) (int, error) {
	n, err := p.r.Read(b)
	p.read += int64(n)
	if p.total > 0 {
		pct := float64(p.read) / float64(p.total) * 100
		if pct > 100 {
			pct = 100
		}
		fmt.Printf("\r%s %5.1f%%", p.name, pct)
	} else if n > 0 {
		fmt.Printf("\r%s %d bytes", p.name, p.read)
	}
	return n, err
}

func (p *progressReader) done() {
	fmt.Println()
}
