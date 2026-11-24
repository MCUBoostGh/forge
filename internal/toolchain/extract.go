package toolchain

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

// Extract unpacks a .zip or .tar.gz archive into dest.
// Paths that escape dest are rejected. Relative tar symlinks are recreated.
func Extract(archivePath, dest string) error {
	switch {
	case strings.HasSuffix(archivePath, ".tar.gz"), strings.HasSuffix(archivePath, ".tgz"):
		return extractTarGz(archivePath, dest)
	case strings.HasSuffix(archivePath, ".zip"):
		return extractZip(archivePath, dest)
	default:
		return fmt.Errorf("unsupported archive %q", archivePath)
	}
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

	if err := os.MkdirAll(dest, 0755); err != nil {
		return err
	}

	tr := tar.NewReader(gz)
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
			return fmt.Errorf("extract %s: %w", hdr.Name, err)
		}

		switch hdr.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(target, 0755); err != nil {
				return err
			}
		case tar.TypeReg, tar.TypeRegA:
			if err := writeFile(target, tr, hdr.FileInfo().Mode().Perm()); err != nil {
				return err
			}
		case tar.TypeSymlink:
			if err := writeSymlink(dest, target, hdr.Linkname); err != nil {
				return fmt.Errorf("extract %s: %w", hdr.Name, err)
			}
		default:
			return fmt.Errorf("unsupported tar entry %q type %c", hdr.Name, hdr.Typeflag)
		}
	}
}

func extractZip(src, dest string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer r.Close()

	if err := os.MkdirAll(dest, 0755); err != nil {
		return err
	}

	for _, file := range r.File {
		target, err := safeJoin(dest, file.Name)
		if err != nil {
			return fmt.Errorf("extract %s: %w", file.Name, err)
		}
		mode := file.FileInfo().Mode()
		if file.FileInfo().IsDir() || strings.HasSuffix(file.Name, "/") {
			if err := os.MkdirAll(target, 0755); err != nil {
				return err
			}
			continue
		}
		if mode&os.ModeSymlink != 0 {
			rc, err := file.Open()
			if err != nil {
				return err
			}
			link, err := io.ReadAll(rc)
			rc.Close()
			if err != nil {
				return err
			}
			if err := writeSymlink(dest, target, string(link)); err != nil {
				return fmt.Errorf("extract %s: %w", file.Name, err)
			}
			continue
		}
		rc, err := file.Open()
		if err != nil {
			return err
		}
		err = writeFile(target, rc, mode.Perm())
		rc.Close()
		if err != nil {
			return err
		}
	}
	return nil
}

func writeFile(target string, r io.Reader, perm os.FileMode) error {
	if perm == 0 {
		perm = 0644
	}
	if err := os.MkdirAll(filepath.Dir(target), 0755); err != nil {
		return err
	}
	out, err := os.OpenFile(target, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, perm)
	if err != nil {
		return err
	}
	_, copyErr := io.Copy(out, r)
	closeErr := out.Close()
	if copyErr != nil {
		return copyErr
	}
	return closeErr
}

func writeSymlink(dest, target, linkname string) error {
	if filepath.IsAbs(linkname) {
		return fmt.Errorf("absolute symlink %q", linkname)
	}
	rel := filepath.FromSlash(linkname)
	resolved := filepath.Clean(filepath.Join(filepath.Dir(target), rel))
	if !inside(dest, resolved) {
		return os.ErrInvalid
	}
	if err := os.MkdirAll(filepath.Dir(target), 0755); err != nil {
		return err
	}
	if err := os.Symlink(rel, target); err != nil {
		if runtime.GOOS != "windows" {
			return err
		}
		return copyRegular(resolved, target)
	}
	return nil
}

func copyRegular(src, dst string) error {
	info, err := os.Lstat(src)
	if err != nil {
		return err
	}
	if !info.Mode().IsRegular() {
		return fmt.Errorf("symlink target %s is not a regular file", src)
	}
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	return writeFile(dst, in, info.Mode().Perm())
}

func safeJoin(dest, name string) (string, error) {
	clean := filepath.Clean(name)
	if clean == ".." || strings.HasPrefix(clean, ".."+string(os.PathSeparator)) || filepath.IsAbs(clean) {
		return "", os.ErrInvalid
	}
	target := filepath.Join(dest, clean)
	if !inside(dest, target) {
		return "", os.ErrInvalid
	}
	return target, nil
}

func inside(dest, path string) bool {
	dest = filepath.Clean(dest)
	path = filepath.Clean(path)
	if path == dest {
		return true
	}
	return strings.HasPrefix(path, dest+string(os.PathSeparator))
}
