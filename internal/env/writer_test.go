package env

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestWriter_Write_CreatesFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".env")

	w := NewWriter(path, false)
	secrets := map[string]string{
		"DB_HOST": "localhost",
		"DB_PORT": "5432",
	}

	if err := w.Write(secrets); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("could not read written file: %v", err)
	}

	content := string(data)
	if !strings.Contains(content, "DB_HOST=localhost") {
		t.Errorf("expected DB_HOST=localhost in output, got:\n%s", content)
	}
	if !strings.Contains(content, "DB_PORT=5432") {
		t.Errorf("expected DB_PORT=5432 in output, got:\n%s", content)
	}
}

func TestWriter_Write_SortedKeys(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".env")

	w := NewWriter(path, false)
	secrets := map[string]string{"Z_KEY": "z", "A_KEY": "a", "M_KEY": "m"}

	if err := w.Write(secrets); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data, _ := os.ReadFile(path)
	lines := strings.Split(strings.TrimSpace(string(data)), "\n")
	if lines[0] != "A_KEY=a" || lines[1] != "M_KEY=m" || lines[2] != "Z_KEY=z" {
		t.Errorf("keys not sorted, got: %v", lines)
	}
}

func TestWriter_Write_BackupOnRotate(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".env")

	// Write initial file.
	if err := os.WriteFile(path, []byte("OLD=value\n"), 0600); err != nil {
		t.Fatal(err)
	}

	w := NewWriter(path, true)
	if err := w.Write(map[string]string{"NEW": "value"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	bak, err := os.ReadFile(path + ".bak")
	if err != nil {
		t.Fatalf("backup file not created: %v", err)
	}
	if !strings.Contains(string(bak), "OLD=value") {
		t.Errorf("backup does not contain original content: %s", bak)
	}
}

func TestEscape_WrapsSpaces(t *testing.T) {
	result := escape("hello world")
	if result != `"hello world"` {
		t.Errorf("expected quoted value, got: %s", result)
	}
}

func TestEscape_PlainValue(t *testing.T) {
	result := escape("simplevalue")
	if result != "simplevalue" {
		t.Errorf("expected unquoted value, got: %s", result)
	}
}
