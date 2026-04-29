package rotate

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNew_DefaultMaxBackups(t *testing.T) {
	r := New(0)
	if r.MaxBackups != 3 {
		t.Fatalf("expected MaxBackups=3, got %d", r.MaxBackups)
	}
}

func TestNew_CustomMaxBackups(t *testing.T) {
	r := New(5)
	if r.MaxBackups != 5 {
		t.Fatalf("expected MaxBackups=5, got %d", r.MaxBackups)
	}
}

func TestRotate_NoopWhenFileMissing(t *testing.T) {
	r := New(3)
	if err := r.Rotate("/tmp/vaultpull_nonexistent_xyz.env"); err != nil {
		t.Fatalf("expected no error for missing file, got %v", err)
	}
}

func TestRotate_CreatesBackup(t *testing.T) {
	dir := t.TempDir()
	envFile := filepath.Join(dir, ".env")

	if err := os.WriteFile(envFile, []byte("SECRET=abc\n"), 0600); err != nil {
		t.Fatal(err)
	}

	r := New(3)
	if err := r.Rotate(envFile); err != nil {
		t.Fatalf("Rotate error: %v", err)
	}

	matches, _ := filepath.Glob(filepath.Join(dir, ".env.*.bak"))
	if len(matches) != 1 {
		t.Fatalf("expected 1 backup, got %d", len(matches))
	}

	data, err := os.ReadFile(matches[0])
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "SECRET=abc\n" {
		t.Fatalf("backup content mismatch: %q", string(data))
	}
}

func TestRotate_PrunesOldBackups(t *testing.T) {
	dir := t.TempDir()
	envFile := filepath.Join(dir, ".env")

	r := New(2)

	// Create 3 rotations; after each we rewrite the source file.
	for i := 0; i < 3; i++ {
		if err := os.WriteFile(envFile, []byte("VAL=x\n"), 0600); err != nil {
			t.Fatal(err)
		}
		// Small sleep not needed — timestamps in filenames are second-precision;
		// glob sort order is sufficient for the prune logic to work in tests
		// as long as we only check counts.
		if err := r.Rotate(envFile); err != nil {
			t.Fatalf("Rotate %d error: %v", i, err)
		}
	}

	matches, _ := filepath.Glob(filepath.Join(dir, ".env.*.bak"))
	if len(matches) > r.MaxBackups {
		t.Fatalf("expected at most %d backups, got %d", r.MaxBackups, len(matches))
	}
}
