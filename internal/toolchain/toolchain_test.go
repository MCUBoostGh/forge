package toolchain

import (
	"archive/tar"
	"archive/zip"
	"bytes"
	"compress/gzip"
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestList_FourTools(t *testing.T) {
	tools, err := list("linux", "amd64")
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	got := map[string]bool{}
	for _, tool := range tools {
		got[tool.Name] = true
		if strings.Contains(tool.URL, "{version}") || strings.Contains(tool.ArchiveRoot, "{version}") {
			t.Errorf("%s still has {version}: %s %s", tool.Name, tool.URL, tool.ArchiveRoot)
		}
		if !strings.HasSuffix(tool.URL, ".tar.gz") {
			t.Errorf("%s linux URL = %s, want .tar.gz", tool.Name, tool.URL)
		}
		if len(tool.SHA256) != 64 {
			t.Errorf("%s sha256 length = %d", tool.Name, len(tool.SHA256))
		}
	}
	for _, name := range []string{"cmake", "ninja", "gcc-arm-none-eabi", "openocd"} {
		if !got[name] {
			t.Errorf("missing %s", name)
		}
	}
	if len(tools) != 4 {
		t.Errorf("len = %d, want 4", len(tools))
	}
	if got["gcc"] {
		t.Error("host gcc must not be a catalog entry")
	}
}

func TestLookup_CMakeHosts(t *testing.T) {
	linux, err := lookup("cmake", "linux", "amd64")
	if err != nil {
		t.Fatalf("linux: %v", err)
	}
	if linux.URL != "https://github.com/Kitware/CMake/releases/download/v4.5.1/cmake-4.5.1-linux-x86_64.tar.gz" {
		t.Errorf("linux URL = %s", linux.URL)
	}
	if linux.ArchiveRoot != "cmake-4.5.1-linux-x86_64" {
		t.Errorf("linux archive root = %s", linux.ArchiveRoot)
	}
	if linux.SHA256 != "83b71181fb62b4fe8f26067a9321448572e7a43bd04ff6c7f54a24c33698f16d" {
		t.Errorf("linux sha256 = %s", linux.SHA256)
	}
	if linux.Binary != "bin/cmake" {
		t.Errorf("binary = %s", linux.Binary)
	}

	windows, err := lookup("cmake", "windows", "amd64")
	if err != nil {
		t.Fatalf("windows: %v", err)
	}
	if windows.URL != "https://github.com/Kitware/CMake/releases/download/v4.5.1/cmake-4.5.1-windows-x86_64.zip" {
		t.Errorf("windows URL = %s", windows.URL)
	}
	if windows.ArchiveRoot != "cmake-4.5.1-windows-x86_64" {
		t.Errorf("windows archive root = %s", windows.ArchiveRoot)
	}
	if !strings.HasSuffix(windows.URL, ".zip") {
		t.Errorf("windows URL = %s, want .zip", windows.URL)
	}
}

func TestLookup_Errors(t *testing.T) {
	if _, err := lookup("gcc", "linux", "amd64"); err == nil {
		t.Fatal("expected error for unknown tool")
	}
	if _, err := lookup("cmake", "darwin", "amd64"); err == nil {
		t.Fatal("expected error for darwin")
	}
	if _, err := lookup("cmake", "linux", "arm64"); err == nil {
		t.Fatal("expected error for arm64")
	}
}

func TestList_MatchesHost(t *testing.T) {
	tools, err := List()
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	suffix := ".tar.gz"
	if runtime.GOOS == "windows" {
		suffix = ".zip"
	}
	if len(tools) != 4 {
		t.Fatalf("len = %d", len(tools))
	}
	for _, tool := range tools {
		if !strings.HasSuffix(tool.URL, suffix) {
			t.Errorf("%s URL = %s, want suffix %s", tool.Name, tool.URL, suffix)
		}
	}
}

func TestCacheDirAndBinaryPath(t *testing.T) {
	tool, err := lookup("cmake", "linux", "amd64")
	if err != nil {
		t.Fatalf("lookup: %v", err)
	}
	dir, err := CacheDir(tool)
	if err != nil {
		t.Fatalf("CacheDir: %v", err)
	}
	want := filepath.Join("forge", "toolchains", "cmake", "4.5.1")
	if !strings.HasSuffix(dir, want) {
		t.Errorf("CacheDir = %s, want suffix %s", dir, want)
	}
	bin, err := BinaryPath(tool)
	if err != nil {
		t.Fatalf("BinaryPath: %v", err)
	}
	binSuffix := filepath.Join("bin", "cmake")
	if runtime.GOOS == "windows" {
		binSuffix += ".exe"
	}
	if !strings.HasSuffix(bin, binSuffix) {
		t.Errorf("BinaryPath = %s, want suffix %s", bin, binSuffix)
	}
}

func TestExtract_TarGzAndZip(t *testing.T) {
	root := t.TempDir()

	tarPath := filepath.Join(root, "tool.tar.gz")
	if err := writeTarGz(tarPath, []tarItem{
		{name: "bin/cmake", body: "cmake-bin", mode: 0755},
		{name: "bin/cmake-link", link: "cmake"},
	}); err != nil {
		t.Fatalf("write tar: %v", err)
	}
	tarDest := filepath.Join(root, "tar")
	if err := Extract(tarPath, tarDest); err != nil {
		t.Fatalf("extract tar: %v", err)
	}
	body, err := os.ReadFile(filepath.Join(tarDest, "bin", "cmake"))
	if err != nil {
		t.Fatalf("read cmake: %v", err)
	}
	if string(body) != "cmake-bin" {
		t.Errorf("cmake body = %q", body)
	}
	linkBody, err := os.ReadFile(filepath.Join(tarDest, "bin", "cmake-link"))
	if err != nil {
		t.Fatalf("read link: %v", err)
	}
	if string(linkBody) != "cmake-bin" {
		t.Errorf("link body = %q", linkBody)
	}
	if canSymlink(t) {
		info, err := os.Lstat(filepath.Join(tarDest, "bin", "cmake-link"))
		if err != nil {
			t.Fatal(err)
		}
		if info.Mode()&os.ModeSymlink == 0 {
			t.Errorf("cmake-link mode = %v, want symlink", info.Mode())
		}
	}

	zipPath := filepath.Join(root, "tool.zip")
	if err := writeZip(zipPath, "bin/ninja", "ninja-bin"); err != nil {
		t.Fatalf("write zip: %v", err)
	}
	zipDest := filepath.Join(root, "zip")
	if err := Extract(zipPath, zipDest); err != nil {
		t.Fatalf("extract zip: %v", err)
	}
	got, err := os.ReadFile(filepath.Join(zipDest, "bin", "ninja"))
	if err != nil {
		t.Fatalf("read ninja: %v", err)
	}
	if string(got) != "ninja-bin" {
		t.Errorf("ninja body = %q", got)
	}
}

func TestExtract_RejectsPathTraversal(t *testing.T) {
	root := t.TempDir()

	tarPath := filepath.Join(root, "slip.tar.gz")
	if err := writeTarGz(tarPath, []tarItem{{name: "../evil.txt", body: "no"}}); err != nil {
		t.Fatalf("write tar: %v", err)
	}
	if err := Extract(tarPath, filepath.Join(root, "out-tar")); err == nil {
		t.Fatal("expected tar slip error")
	} else if !errors.Is(err, os.ErrInvalid) {
		t.Fatalf("tar slip err = %v", err)
	}

	zipPath := filepath.Join(root, "slip.zip")
	if err := writeZip(zipPath, "../evil.txt", "no"); err != nil {
		t.Fatalf("write zip: %v", err)
	}
	if err := Extract(zipPath, filepath.Join(root, "out-zip")); err == nil {
		t.Fatal("expected zip slip error")
	} else if !errors.Is(err, os.ErrInvalid) {
		t.Fatalf("zip slip err = %v", err)
	}
	if _, err := os.Stat(filepath.Join(root, "evil.txt")); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("evil.txt stat = %v", err)
	}
}

func TestExtract_RejectsEscapingSymlink(t *testing.T) {
	root := t.TempDir()
	tarPath := filepath.Join(root, "link.tar.gz")
	if err := writeTarGz(tarPath, []tarItem{{name: "bin/escape", link: "../../evil.txt"}}); err != nil {
		t.Fatalf("write tar: %v", err)
	}
	err := Extract(tarPath, filepath.Join(root, "out"))
	if err == nil {
		t.Fatal("expected escaping symlink error")
	}
	if !errors.Is(err, os.ErrInvalid) {
		t.Fatalf("err = %v", err)
	}
}

type tarItem struct {
	name string
	body string
	link string
	mode int64
}

func writeTarGz(path string, items []tarItem) error {
	var buf bytes.Buffer
	tw := tar.NewWriter(&buf)
	for _, item := range items {
		hdr := &tar.Header{Name: item.name}
		if item.link != "" {
			hdr.Typeflag = tar.TypeSymlink
			hdr.Linkname = item.link
			hdr.Mode = 0777
		} else {
			hdr.Typeflag = tar.TypeReg
			hdr.Mode = item.mode
			if hdr.Mode == 0 {
				hdr.Mode = 0644
			}
			hdr.Size = int64(len(item.body))
		}
		if err := tw.WriteHeader(hdr); err != nil {
			return err
		}
		if item.link == "" {
			if _, err := tw.Write([]byte(item.body)); err != nil {
				return err
			}
		}
	}
	if err := tw.Close(); err != nil {
		return err
	}
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	gz := gzip.NewWriter(f)
	if _, err := gz.Write(buf.Bytes()); err != nil {
		gz.Close()
		return err
	}
	return gz.Close()
}

func writeZip(path, name, body string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	zw := zip.NewWriter(f)
	w, err := zw.Create(name)
	if err != nil {
		zw.Close()
		return err
	}
	if _, err := w.Write([]byte(body)); err != nil {
		zw.Close()
		return err
	}
	return zw.Close()
}

func canSymlink(t *testing.T) bool {
	t.Helper()
	dir := t.TempDir()
	src := filepath.Join(dir, "a")
	if err := os.WriteFile(src, []byte("x"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink("a", filepath.Join(dir, "b")); err != nil {
		return false
	}
	return true
}
