package audit_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"strings"
	"testing"

	"github.com/example/vaultpull/internal/audit"
)

func TestLogger_Log_Success(t *testing.T) {
	var buf bytes.Buffer
	l := audit.NewLogger(&buf)

	if err := l.Log("sync", "secret/app", ".env", nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var entry audit.Entry
	if err := json.Unmarshal(buf.Bytes(), &entry); err != nil {
		t.Fatalf("could not unmarshal output: %v", err)
	}

	if entry.Event != "sync" {
		t.Errorf("expected event=sync, got %q", entry.Event)
	}
	if entry.Path != "secret/app" {
		t.Errorf("expected path=secret/app, got %q", entry.Path)
	}
	if entry.Target != ".env" {
		t.Errorf("expected target=.env, got %q", entry.Target)
	}
	if entry.Error != "" {
		t.Errorf("expected no error field, got %q", entry.Error)
	}
}

func TestLogger_Log_WithError(t *testing.T) {
	var buf bytes.Buffer
	l := audit.NewLogger(&buf)

	if err := l.Log("sync", "secret/app", ".env", errors.New("permission denied")); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var entry audit.Entry
	if err := json.Unmarshal(buf.Bytes(), &entry); err != nil {
		t.Fatalf("could not unmarshal output: %v", err)
	}

	if entry.Error != "permission denied" {
		t.Errorf("expected error field, got %q", entry.Error)
	}
}

func TestLogger_Log_MultipleEntries(t *testing.T) {
	var buf bytes.Buffer
	l := audit.NewLogger(&buf)

	_ = l.Log("sync", "secret/a", "a.env", nil)
	_ = l.Log("rotate", "secret/b", "b.env", nil)

	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) != 2 {
		t.Fatalf("expected 2 lines, got %d", len(lines))
	}
}

func TestNewLogger_NilUsesStdout(t *testing.T) {
	// Should not panic
	l := audit.NewLogger(nil)
	if l == nil {
		t.Fatal("expected non-nil logger")
	}
}
