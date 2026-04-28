package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoad_Defaults(t *testing.T) {
	t.Setenv("VAULT_TOKEN", "test-token")
	t.Setenv("VAULTPULL_VAULT_PATH", "secret/data/myapp")

	cfg, err := Load("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.VaultAddr != "http://127.0.0.1:8200" {
		t.Errorf("expected default vault_addr, got %q", cfg.VaultAddr)
	}
	if cfg.OutputFile != ".env" {
		t.Errorf("expected default output_file '.env', got %q", cfg.OutputFile)
	}
	if cfg.Rotate {
		t.Errorf("expected rotate to default to false")
	}
}

func TestLoad_FromFile(t *testing.T) {
	dir := t.TempDir()
	cfgPath := filepath.Join(dir, "config.yaml")

	content := []byte(`
vault_addr: "http://vault.example.com:8200"
vault_token: "s.abc123"
vault_path: "secret/data/prod"
output_file: "prod.env"
rotate: true
`)
	if err := os.WriteFile(cfgPath, content, 0o600); err != nil {
		t.Fatalf("writing temp config: %v", err)
	}

	cfg, err := Load(cfgPath)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.VaultAddr != "http://vault.example.com:8200" {
		t.Errorf("unexpected vault_addr: %q", cfg.VaultAddr)
	}
	if cfg.VaultToken != "s.abc123" {
		t.Errorf("unexpected vault_token: %q", cfg.VaultToken)
	}
	if cfg.VaultPath != "secret/data/prod" {
		t.Errorf("unexpected vault_path: %q", cfg.VaultPath)
	}
	if cfg.OutputFile != "prod.env" {
		t.Errorf("unexpected output_file: %q", cfg.OutputFile)
	}
	if !cfg.Rotate {
		t.Errorf("expected rotate to be true")
	}
}

func TestLoad_ValidationErrors(t *testing.T) {
	tests := []struct {
		name    string
		envVars map[string]string
	}{
		{
			name: "missing vault_token",
			envVars: map[string]string{
				"VAULTPULL_VAULT_PATH": "secret/data/app",
			},
		},
		{
			name: "missing vault_path",
			envVars: map[string]string{
				"VAULT_TOKEN": "s.token",
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			for k, v := range tc.envVars {
				t.Setenv(k, v)
			}
			_, err := Load("")
			if err == nil {
				t.Error("expected validation error, got nil")
			}
		})
	}
}
