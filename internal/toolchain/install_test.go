package toolchain

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync/atomic"
	"testing"
)

func TestInstall_DownloadsAndSkips(t *testing.T) {
	root := t.TempDir()
	archive := filepath.Join(root, "tool.tar.gz")
	entry := "pkg/bin/tool"
	if runtime.GOOS == "windows" {
		entry = "pkg/bin/tool.exe"
	}
	if err := writeTarGz(archive, []tarItem{{name: entry, body: "hello", mode: 0755}}); err != nil {
		t.Fatalf("tar: %v", err)
	}
	sum := fileSHA256(t, archive)

	var hits atomic.Int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hits.Add(1)
		http.ServeFile(w, r, archive)
	}))
	t.Cleanup(srv.Close)

	dest := filepath.Join(root, "installed")
	tool := Tool{
		Name:        "demo",
		Version:     "1",
		Binary:      "bin/tool",
		URL:         srv.URL + "/tool.tar.gz",
		SHA256:      sum,
		ArchiveRoot: "pkg",
	}
	if err := Install(tool, InstallOptions{Dir: dest}); err != nil {
		t.Fatalf("Install: %v", err)
	}
	body, err := os.ReadFile(executablePath(dest, tool.Binary))
	if err != nil {
		t.Fatalf("read binary: %v", err)
	}
	if string(body) != "hello" {
		t.Errorf("body = %q", body)
	}
	if err := Install(tool, InstallOptions{Dir: dest}); err != nil {
		t.Fatalf("second Install: %v", err)
	}
	if hits.Load() != 1 {
		t.Fatalf("download hits = %d, want 1", hits.Load())
	}
}

func TestInstall_SHA256Mismatch(t *testing.T) {
	root := t.TempDir()
	archive := filepath.Join(root, "tool.tar.gz")
	if err := writeTarGz(archive, []tarItem{{name: "pkg/bin/tool", body: "hello"}}); err != nil {
		t.Fatalf("tar: %v", err)
	}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, archive)
	}))
	t.Cleanup(srv.Close)

	dest := filepath.Join(root, "installed")
	err := Install(Tool{
		Name:        "demo",
		Version:     "1",
		Binary:      "bin/tool",
		URL:         srv.URL + "/tool.tar.gz",
		SHA256:      strings.Repeat("ab", 32),
		ArchiveRoot: "pkg",
	}, InstallOptions{Dir: dest})
	if err == nil || !strings.Contains(err.Error(), "sha256 mismatch") {
		t.Fatalf("err = %v", err)
	}
	if _, statErr := os.Stat(dest); !os.IsNotExist(statErr) {
		t.Fatalf("dest stat = %v, want not exist", statErr)
	}
}

func TestInstall_DryRun(t *testing.T) {
	dest := filepath.Join(t.TempDir(), "installed")
	tool := Tool{
		Name:    "demo",
		Version: "1",
		Binary:  "bin/tool",
		URL:     "https://example.test/tool.tar.gz",
	}
	out := captureInstall(t, func() {
		if err := Install(tool, InstallOptions{DryRun: true, Dir: dest}); err != nil {
			t.Fatalf("Install: %v", err)
		}
	})
	if !strings.Contains(out, tool.URL) || !strings.Contains(out, dest) {
		t.Fatalf("dry-run output = %q", out)
	}
	if _, err := os.Stat(dest); !os.IsNotExist(err) {
		t.Fatalf("dest stat = %v, want not exist", err)
	}
}

func fileSHA256(t *testing.T, path string) string {
	t.Helper()
	f, err := os.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		t.Fatal(err)
	}
	return hex.EncodeToString(h.Sum(nil))
}

func captureInstall(t *testing.T, fn func()) string {
	t.Helper()
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	old := os.Stdout
	os.Stdout = w
	defer func() { os.Stdout = old }()
	fn()
	_ = w.Close()
	b, err := io.ReadAll(r)
	if err != nil {
		t.Fatal(err)
	}
	return string(b)
}
