package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

const (
	reset  = "\033[0m"
	red    = "\033[31m"
	green  = "\033[32m"
	yellow = "\033[33m"
	cyan   = "\033[36m"
	bold   = "\033[1m"
)

func Info(v ...any) {
	printLine(os.Stdout, cyan, "INFO:", v...)
}

func Infof(format string, v ...any) {
	printFormat(os.Stdout, cyan, "INFO:", format, v...)
}

func Success(v ...any) {
	printLine(os.Stdout, green, "SUCCESS:", v...)
}

func Successf(format string, v ...any) {
	printFormat(os.Stdout, green, "SUCCESS:", format, v...)
}

func Warning(v ...any) {
	printLine(os.Stderr, yellow, "WARNING:", v...)
}

func Warningf(format string, v ...any) {
	printFormat(os.Stderr, yellow, "WARNING:", format, v...)
}

func Error(v ...any) {
	printLine(os.Stderr, red, "ERROR:", v...)
}

func Errorf(format string, v ...any) {
	printFormat(os.Stderr, red, "ERROR:", format, v...)
}

var exit = os.Exit

func Fatal(v ...any) {
	printLine(os.Stderr, bold+red, "FATAL:", v...)
	exit(1)
}

func Println(v ...any) {
	log.Println(v...)
}

func Fatalf(format string, v ...any) {
	printFormat(os.Stderr, bold+red, "FATAL:", format, v...)
	exit(1)
}

func printLine(w io.Writer, color, prefix string, v ...any) {
	fmt.Fprintln(w, formatMessage(w, color, prefix, fmt.Sprint(v...)))
}

func printFormat(w io.Writer, color, prefix, format string, v ...any) {
	fmt.Fprintln(w, formatMessage(w, color, prefix, fmt.Sprintf(format, v...)))
}

func formatMessage(w io.Writer, color, prefix, msg string) string {
	msg = strings.TrimRight(msg, "\n")
	if colorEnabled(w) {
		return color + prefix + reset + " " + msg
	}
	return prefix + " " + msg
}

func colorEnabled(w io.Writer) bool {
	if os.Getenv("NO_COLOR") != "" || os.Getenv("TERM") == "dumb" {
		return false
	}
	f, ok := w.(*os.File)
	if !ok {
		return false
	}
	info, err := f.Stat()
	if err != nil {
		return false
	}
	return info.Mode()&os.ModeCharDevice != 0
}
