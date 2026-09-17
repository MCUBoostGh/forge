package cmd

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"forge/internal/logger"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type Dependency struct {
	Name    string
	Version string
	URL     string
	Size    int64
}

var registry = map[string]Dependency{
	"cmsis": {
		Name:    "cmsis",
		Version: "5.9.0",
		URL:     "https://github.com/ARM-software/CMSIS_5/archive/refs/tags/5.9.0.tar.gz",
		Size:    40658690,
	},
}

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

const progressBarWidth = 30

type progressReader struct {
	name  string
	r     io.Reader
	total int64
	read  int64
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

func (p *progressReader) Read(b []byte) (int, error) {
	n, err := p.r.Read(b)
	p.read += int64(n)
	p.render()
	return n, err
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

func download(dep Dependency) (string, error) {

	msgErr := fmt.Errorf("Failed to download dependency.")

	base, err := os.UserCacheDir()
	if err != nil {
		logger.Error(err)
		return "", msgErr
	}

	cacheDir := filepath.Join(base, "forge", "packages", dep.Name, dep.Version)

	if _, err := os.Stat(cacheDir); err == nil {
		return cacheDir, nil
	}

	req, err := http.NewRequest(http.MethodGet, dep.URL, nil)
	if err != nil {
		logger.Error(err)
		return "", msgErr
	}
	req.Header.Set("Accept-Encoding", "identity")

	transport := http.DefaultTransport.(*http.Transport).Clone()
	transport.DisableCompression = true
	client := &http.Client{Transport: transport}

	response, err := client.Do(req)
	if err != nil {
		logger.Error(err)
		return "", msgErr
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return "", fmt.Errorf("download failed: %s", response.Status)
	}

	progress := &progressReader{
		name:  dep.Name,
		r:     response.Body,
		total: contentLength(response, dep.Size),
	}

	tmp, err := os.CreateTemp("", "forge-"+dep.Name+"-*.tar.gz")
	if err != nil {
		logger.Error(err)
		return "", msgErr
	}

	defer os.Remove(tmp.Name())

	if _, err := io.Copy(tmp, progress); err != nil {
		fmt.Println()
		logger.Error(err)
		return "", msgErr
	}
	progress.done()

	tmp.Close()

	if err := extractTarGz(tmp.Name(), cacheDir); err != nil {
		logger.Error(err)
		return "", msgErr
	}

	return cacheDir, nil
}
