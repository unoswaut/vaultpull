package audit_test

import (
	"bufio"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/example/vaultpull/internal/audit"
)

func TestNewFileLogger_CreatesFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "audit", "audit.log")

	fl, err := audit.NewFileLogger(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer fl.Close()

	if _, err := os.Stat(path); err != nil {
		t.Fatalf("expected file to exist: %v", err)
	}
}

func TestNewFileLogger_WritesEntry(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "audit.log")

	fl, err := audit.NewFileLogger(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if err := fl.Log("sync", "secret/db", "db.env", nil); err != nil {
		t.Fatalf("log error: %v", err)
	}
	fl.Close()

	f, err := os.Open(path)
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	if !scanner.Scan() {
		t.Fatal("expected at least one line")
	}

	var entry audit.Entry
	if err := json.Unmarshal(scanner.Bytes(), &entry); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if entry.Event != "sync" {
		t.Errorf("expected event=sync, got %q", entry.Event)
	}
}

func TestNewFileLogger_InvalidPath(t *testing.T) {
	// /dev/null/audit.log is not a valid directory path on Linux
	_, err := audit.NewFileLogger("/dev/null/subdir/audit.log")
	if err == nil {
		t.Fatal("expected error for invalid path")
	}
}
