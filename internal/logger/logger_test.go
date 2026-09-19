package logger

import (
	"bytes"
	"io"
	"log"
	"os"
	"strings"
	"testing"
)

func TestFormatMessage_NoColorForBuffer(t *testing.T) {
	var buf bytes.Buffer
	got := formatMessage(&buf, red, "ERROR:", "boom\n")
	if got != "ERROR: boom" {
		t.Errorf("formatMessage = %q", got)
	}
}

func TestColorEnabled_NOCOLOR(t *testing.T) {
	t.Setenv("NO_COLOR", "1")
	t.Setenv("TERM", "xterm")
	if colorEnabled(os.Stdout) {
		t.Fatal("expected color disabled when NO_COLOR is set")
	}
}

func TestColorEnabled_DumbTERM(t *testing.T) {
	t.Setenv("NO_COLOR", "")
	t.Setenv("TERM", "dumb")
	if colorEnabled(os.Stdout) {
		t.Fatal("expected color disabled when TERM=dumb")
	}
}

func TestColorEnabled_NonFile(t *testing.T) {
	t.Setenv("NO_COLOR", "")
	t.Setenv("TERM", "xterm")
	if colorEnabled(&bytes.Buffer{}) {
		t.Fatal("expected color disabled for non-file writer")
	}
}

func TestInfoAndErrorPrefixes(t *testing.T) {
	t.Setenv("NO_COLOR", "1")
	out, errOut := captureStd(t, func() {
		Info("hello")
		Infof("n=%d", 3)
		Success("ok")
		Successf("done %s", "x")
		Warning("careful")
		Warningf("w=%s", "y")
		Error("fail")
		Errorf("e=%s", "z")
	})
	if !strings.Contains(out, "INFO: hello") || !strings.Contains(out, "INFO: n=3") {
		t.Errorf("stdout = %q", out)
	}
	if !strings.Contains(out, "SUCCESS: ok") || !strings.Contains(out, "SUCCESS: done x") {
		t.Errorf("stdout success = %q", out)
	}
	if !strings.Contains(errOut, "WARNING: careful") || !strings.Contains(errOut, "ERROR: fail") {
		t.Errorf("stderr = %q", errOut)
	}
	if strings.Contains(out, "\033[") || strings.Contains(errOut, "\033[") {
		t.Fatal("expected no ANSI color when NO_COLOR is set")
	}
}

func TestPrintln(t *testing.T) {
	var buf bytes.Buffer
	old := log.Writer()
	log.SetOutput(&buf)
	log.SetFlags(0)
	t.Cleanup(func() {
		log.SetOutput(old)
	})
	Println("plain")
	if !strings.Contains(buf.String(), "plain") {
		t.Errorf("Println wrote %q", buf.String())
	}
}

func TestFatal_UsesExitHook(t *testing.T) {
	t.Setenv("NO_COLOR", "1")
	old := exit
	code := 0
	exit = func(c int) { code = c }
	t.Cleanup(func() { exit = old })

	_, errOut := captureStd(t, func() {
		Fatal("stop")
		Fatalf("code=%d", 7)
	})
	if code != 1 {
		t.Errorf("exit code = %d, want 1", code)
	}
	if !strings.Contains(errOut, "FATAL: stop") || !strings.Contains(errOut, "FATAL: code=7") {
		t.Errorf("stderr = %q", errOut)
	}
}

func captureStd(t *testing.T, fn func()) (stdout, stderr string) {
	t.Helper()
	rOut, wOut, err := os.Pipe()
	if err != nil {
		t.Fatalf("stdout pipe: %v", err)
	}
	rErr, wErr, err := os.Pipe()
	if err != nil {
		t.Fatalf("stderr pipe: %v", err)
	}
	oldOut, oldErr := os.Stdout, os.Stderr
	os.Stdout, os.Stderr = wOut, wErr
	defer func() {
		os.Stdout, os.Stderr = oldOut, oldErr
	}()

	fn()

	_ = wOut.Close()
	_ = wErr.Close()
	outBytes, err := io.ReadAll(rOut)
	if err != nil {
		t.Fatalf("read stdout: %v", err)
	}
	errBytes, err := io.ReadAll(rErr)
	if err != nil {
		t.Fatalf("read stderr: %v", err)
	}
	return string(outBytes), string(errBytes)
}
