package sync

import (
	"fmt"

	"github.com/example/vaultpull/internal/audit"
	"github.com/example/vaultpull/internal/config"
	"github.com/example/vaultpull/internal/env"
	"github.com/example/vaultpull/internal/rotate"
	"github.com/example/vaultpull/internal/vault"
)

// Syncer orchestrates reading secrets from Vault and writing .env files.
type Syncer struct {
	client  vault.Client
	cfg     *config.Config
	rotator *rotate.Rotator
	auditor *audit.Logger
	errors  []error
}

// New creates a Syncer from the given config and vault client.
func New(cfg *config.Config, client vault.Client, auditor *audit.Logger) *Syncer {
	r := rotate.New(cfg.MaxBackups)
	return &Syncer{
		client:  client,
		cfg:     cfg,
		rotator: r,
		auditor: auditor,
	}
}

// Run executes the sync for every mapping in the config.
func (s *Syncer) Run() {
	for _, m := range s.cfg.Mappings {
		secrets, err := s.client.ReadSecrets(m.VaultPath)
		if err != nil {
			s.errors = append(s.errors, fmt.Errorf("vault read %s: %w", m.VaultPath, err))
			if s.auditor != nil {
				_ = s.auditor.Log("error", m.VaultPath, m.EnvFile, err)
			}
			continue
		}

		if err := s.rotator.Rotate(m.EnvFile); err != nil {
			s.errors = append(s.errors, fmt.Errorf("rotate %s: %w", m.EnvFile, err))
			if s.auditor != nil {
				_ = s.auditor.Log("error", m.VaultPath, m.EnvFile, err)
			}
			continue
		}

		w := env.NewWriter(m.EnvFile)
		if err := w.Write(secrets); err != nil {
			s.errors = append(s.errors, fmt.Errorf("write %s: %w", m.EnvFile, err))
			if s.auditor != nil {
				_ = s.auditor.Log("error", m.VaultPath, m.EnvFile, err)
			}
			continue
		}

		if s.auditor != nil {
			_ = s.auditor.Log("sync", m.VaultPath, m.EnvFile, nil)
		}
	}
}

// HasErrors reports whether any errors occurred during Run.
func (s *Syncer) HasErrors() bool {
	return len(s.errors) > 0
}

// Errors returns all errors collected during Run.
func (s *Syncer) Errors() []error {
	return s.errors
}
