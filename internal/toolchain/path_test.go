package toolchain

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestPrependPath_CacheDirsFirst(t *testing.T) {
	sep := string(osPathSep())
	dirs := []string{`C:\cache\cmake\bin`, `C:\cache\ninja\bin`}
	got := PrependPath(`C:\Windows`, dirs)
	want := strings.Join(append(dirs, `C:\Windows`), sep)
	if got != want {
		t.Fatalf("PrependPath = %q, want %q", got, want)
	}
	if !strings.HasPrefix(got, dirs[0]+sep+dirs[1]+sep) {
		t.Fatalf("cache dirs are not first: %q", got)
	}
}

func TestFindExecutable_PrefersFirstDir(t *testing.T) {
	root := t.TempDir()
	cache := filepath.Join(root, "cache")
	other := filepath.Join(root, "other")
	if err := os.MkdirAll(cache, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(other, 0o755); err != nil {
		t.Fatal(err)
	}
	name := "cmake"
	if runtime.GOOS == "windows" {
		name = "cmake.exe"
	}
	want := filepath.Join(cache, name)
	if err := os.WriteFile(want, []byte("cache"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(other, name), []byte("other"), 0o755); err != nil {
		t.Fatal(err)
	}
	sep := string(os.PathListSeparator)
	got, err := FindExecutable("cmake", cache+sep+other)
	if err != nil {
		t.Fatal(err)
	}
	if got != want {
		t.Fatalf("FindExecutable = %q, want %q", got, want)
	}
}

func TestFindExecutable_Missing(t *testing.T) {
	_, err := FindExecutable("cmake", t.TempDir())
	if !os.IsNotExist(err) && err != exec.ErrNotFound {
		t.Fatalf("err = %v, want not found", err)
	}
}

func osPathSep() byte {
	if runtime.GOOS == "windows" {
		return ';'
	}
	return ':'
}
