package log

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestInitLogger_WritesInfoToFile(t *testing.T) {
	dir := t.TempDir()
	logPath := filepath.Join(dir, "app.log")

	InitLogger(logPath, "info")
	t.Cleanup(func() {
		if Logger != nil {
			_ = Logger.Sync()
		}
	})

	Logger.Info("hello-from-test", String("k", "v"))
	_ = Logger.Sync()

	b, err := os.ReadFile(logPath)
	if err != nil {
		t.Fatalf("read log file: %v", err)
	}
	got := string(b)

	if !strings.Contains(got, "hello-from-test") {
		t.Fatalf("expected log file to contain message; got:\n%s", got)
	}
	if !strings.Contains(got, `"k": "v"`) {
		t.Fatalf("expected log file to contain field k; got:\n%s", got)
	}
	if !strings.Contains(got, `"application": "chat-room"`) {
		t.Fatalf("expected log file to contain default field application; got:\n%s", got)
	}
}

func TestInitLogger_LevelFiltersDebug(t *testing.T) {
	dir := t.TempDir()
	logPath := filepath.Join(dir, "level.log")

	InitLogger(logPath, "info")
	t.Cleanup(func() {
		if Logger != nil {
			_ = Logger.Sync()
		}
	})

	Logger.Debug("debug-should-not-appear")
	Logger.Info("info-should-appear")
	_ = Logger.Sync()

	b, err := os.ReadFile(logPath)
	if err != nil {
		t.Fatalf("read log file: %v", err)
	}
	got := string(b)

	if strings.Contains(got, "debug-should-not-appear") {
		t.Fatalf("did not expect debug message in info level; got:\n%s", got)
	}
	if !strings.Contains(got, "info-should-appear") {
		t.Fatalf("expected info message; got:\n%s", got)
	}
}

