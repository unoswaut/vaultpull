package sync_test

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/example/vaultpull/internal/config"
	"github.com/example/vaultpull/internal/sync"
	"github.com/example/vaultpull/internal/vault"
)

func makeConfig(t *testing.T, mappings []config.Mapping, rotate bool) *config.Config {
	t.Helper()
	return &config.Config{
		Mappings: mappings,
		Rotate:   rotate,
	}
}

func TestSyncer_Run_Success(t *testing.T) {
	dir := t.TempDir()
	outFile := filepath.Join(dir, ".env")

	client := &vault.MockClient{
		Secrets: map[string]map[string]string{
			"secret/app": {"KEY": "value"},
		},
	}
	cfg := makeConfig(t, []config.Mapping{{VaultPath: "secret/app", EnvFile: outFile}}, false)

	s := sync.New(client, cfg)
	results := s.Run()

	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Err != nil {
		t.Fatalf("unexpected error: %v", results[0].Err)
	}
	if _, err := os.Stat(outFile); err != nil {
		t.Fatalf("output file not created: %v", err)
	}
}

func TestSyncer_Run_VaultError(t *testing.T) {
	dir := t.TempDir()
	outFile := filepath.Join(dir, ".env")

	client := &vault.MockClient{
		Err: errors.New("vault unavailable"),
	}
	cfg := makeConfig(t, []config.Mapping{{VaultPath: "secret/app", EnvFile: outFile}}, false)

	s := sync.New(client, cfg)
	results := s.Run()

	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestSyncer_Run_MultipleMappings(t *testing.T) {
	dir := t.TempDir()

	client := &vault.MockClient{
		Secrets: map[string]map[string]string{
			"secret/a": {"A": "1"},
			"secret/b": {"B": "2"},
		},
	}
	cfg := makeConfig(t, []config.Mapping{
		{VaultPath: "secret/a", EnvFile: filepath.Join(dir, "a.env")},
		{VaultPath: "secret/b", EnvFile: filepath.Join(dir, "b.env")},
	}, false)

	s := sync.New(client, cfg)
	results := s.Run()

	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
	if sync.HasErrors(results) {
		t.Fatal("unexpected errors in results")
	}
}

func TestHasErrors_False(t *testing.T) {
	results := []sync.Result{{Path: "p", OutFile: "f", Err: nil}}
	if sync.HasErrors(results) {
		t.Fatal("expected no errors")
	}
}

func TestHasErrors_True(t *testing.T) {
	results := []sync.Result{{Err: errors.New("oops")}}
	if !sync.HasErrors(results) {
		t.Fatal("expected errors")
	}
}
