package toolchain

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestMain(m *testing.M) {
	base := strings.TrimSuffix(strings.ToLower(filepath.Base(os.Args[0])), ".exe")
	switch base {
	case "cmake":
		fmt.Println("cmake version 4.5.1")
		os.Exit(0)
	case "ninja":
		fmt.Println("1.12.0")
		os.Exit(0)
	default:
		os.Exit(m.Run())
	}
}

func TestVersionMatches(t *testing.T) {
	cases := []struct {
		catalog string
		output  string
		ok      bool
	}{
		{"4.5.1", "cmake version 4.5.1\n", true},
		{"1.13.1-1", "1.13.1\n", true},
		{"15.2.1-1.1", "arm-none-eabi-gcc (xPack) 15.2.1\n", true},
		{"0.12.0-7", "Open On-Chip Debugger 0.12.0\n", true},
		{"1.13.1-1", "1.12.0\n", false},
		{"4.5.1", "cmake version 3.28.0\n", false},
	}
	for _, tc := range cases {
		if got := versionMatches(tc.catalog, tc.output); got != tc.ok {
			t.Errorf("versionMatches(%q, %q) = %v, want %v", tc.catalog, tc.output, got, tc.ok)
		}
	}
}

func TestFormatReport_AlignsAndCollapses(t *testing.T) {
	got := firstLine("\n    xPack Open On-Chip Debugger 0.12.0\t(dirty)\nLicensed\n")
	if got != "xPack Open On-Chip Debugger 0.12.0 (dirty)" {
		t.Fatalf("firstLine = %q", got)
	}
	report := FormatReport([]CheckResult{
		{
			Tool:     Tool{Name: "cmake", Version: "4.5.1"},
			Detected: "cmake version 4.5.1",
			Path:     `C:\forge\cmake\bin\cmake.exe`,
			Status:   StatusOK,
		},
		{
			Tool:     Tool{Name: "gcc-arm-none-eabi", Version: "15.2.1-1.1"},
			Detected: "arm-none-eabi-gcc (xPack) 15.2.1",
			Path:     `C:\forge\gcc\bin\arm-none-eabi-gcc.exe`,
			Status:   StatusOK,
		},
	})
	if strings.Contains(report, "\t") {
		t.Fatalf("report still uses tabs:\n%s", report)
	}
	lines := strings.Split(report, "\n")
	if len(lines) < 4 {
		t.Fatalf("report = %q", report)
	}
	if !strings.HasPrefix(lines[0], "cmake") || !strings.Contains(lines[0], "ok") {
		t.Errorf("line 0 = %q", lines[0])
	}
	gcc := lines[4]
	if !strings.HasPrefix(gcc, "gcc-arm-none-eabi") || !strings.Contains(gcc, "15.2.1-1.1") {
		t.Errorf("gcc line = %q\nreport:\n%s", gcc, report)
	}
}

func TestCheck_CacheHit(t *testing.T) {
	dir := t.TempDir()
	copyFakeTool(t, filepath.Join(dir, "bin", exeName("cmake")))
	withPATH(t, t.TempDir())

	got, err := Check(Tool{Name: "cmake", Version: "4.5.1", Binary: "bin/cmake"}, dir)
	if err != nil {
		t.Fatalf("Check: %v", err)
	}
	if got.Status != StatusOK {
		t.Fatalf("status = %q, detected %q, path %s", got.Status, got.Detected, got.Path)
	}
	if !strings.Contains(got.Detected, "4.5.1") {
		t.Errorf("detected = %q", got.Detected)
	}
	if got.Path != executablePath(dir, "bin/cmake") {
		t.Errorf("path = %s", got.Path)
	}
}

func TestCheck_VersionMismatch(t *testing.T) {
	dir := t.TempDir()
	copyFakeTool(t, filepath.Join(dir, "bin", exeName("ninja")))
	withPATH(t, t.TempDir())

	got, err := Check(Tool{Name: "ninja", Version: "1.13.1-1", Binary: "bin/ninja"}, dir)
	if err != nil {
		t.Fatalf("Check: %v", err)
	}
	if got.Status != StatusMismatch {
		t.Fatalf("status = %q, detected %q", got.Status, got.Detected)
	}
	if got.Detected != "1.12.0" {
		t.Errorf("detected = %q", got.Detected)
	}
}

func TestCheck_Missing(t *testing.T) {
	withPATH(t, t.TempDir())
	got, err := Check(Tool{Name: "cmake", Version: "4.5.1", Binary: "bin/cmake"}, t.TempDir())
	if err != nil {
		t.Fatalf("Check: %v", err)
	}
	if got.Status != StatusMissing || got.Path != "" {
		t.Fatalf("status = %q path = %q", got.Status, got.Path)
	}
}

func TestCheck_PATHFallback(t *testing.T) {
	binDir := t.TempDir()
	installed := filepath.Join(binDir, exeName("cmake"))
	copyFakeTool(t, installed)
	withPATH(t, binDir)

	got, err := Check(Tool{Name: "cmake", Version: "4.5.1", Binary: "bin/cmake"}, t.TempDir())
	if err != nil {
		t.Fatalf("Check: %v", err)
	}
	if got.Status != StatusOK {
		t.Fatalf("status = %q detected %q path %s", got.Status, got.Detected, got.Path)
	}
	if filepath.Clean(got.Path) != filepath.Clean(installed) {
		t.Errorf("path = %s, want %s", got.Path, installed)
	}
}

func exeName(name string) string {
	if runtime.GOOS == "windows" {
		return name + ".exe"
	}
	return name
}

func copyFakeTool(t *testing.T, dest string) {
	t.Helper()
	src, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}
	if resolved, err := filepath.EvalSymlinks(src); err == nil {
		src = resolved
	}
	in, err := os.Open(src)
	if err != nil {
		t.Fatal(err)
	}
	defer in.Close()
	if err := os.MkdirAll(filepath.Dir(dest), 0755); err != nil {
		t.Fatal(err)
	}
	out, err := os.OpenFile(dest, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0755)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := io.Copy(out, in); err != nil {
		out.Close()
		t.Fatal(err)
	}
	if err := out.Close(); err != nil {
		t.Fatal(err)
	}
}

func withPATH(t *testing.T, dir string) {
	t.Helper()
	old := os.Getenv("PATH")
	if err := os.Setenv("PATH", dir); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Setenv("PATH", old) })
}
