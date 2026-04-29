package sync

import (
	"fmt"
	"log"

	"github.com/example/vaultpull/internal/config"
	"github.com/example/vaultpull/internal/env"
	"github.com/example/vaultpull/internal/vault"
)

// VaultClient is the interface used by Syncer to read secrets.
type VaultClient interface {
	ReadSecrets(path string) (map[string]string, error)
}

// Syncer orchestrates reading secrets from Vault and writing them to .env files.
type Syncer struct {
	client VaultClient
	cfg    *config.Config
}

// New creates a new Syncer with the provided client and config.
func New(client VaultClient, cfg *config.Config) *Syncer {
	return &Syncer{client: client, cfg: cfg}
}

// Result holds the outcome of syncing a single mapping.
type Result struct {
	Path    string
	OutFile string
	Err     error
}

// Run iterates over all configured mappings, fetches secrets, and writes env files.
func (s *Syncer) Run() []Result {
	results := make([]Result, 0, len(s.cfg.Mappings))

	for _, m := range s.cfg.Mappings {
		r := Result{Path: m.VaultPath, OutFile: m.EnvFile}

		secrets, err := s.client.ReadSecrets(m.VaultPath)
		if err != nil {
			r.Err = fmt.Errorf("reading %q: %w", m.VaultPath, err)
			results = append(results, r)
			log.Printf("[error] %v", r.Err)
			continue
		}

		w := env.NewWriter(m.EnvFile, s.cfg.Rotate)
		if err := w.Write(secrets); err != nil {
			r.Err = fmt.Errorf("writing %q: %w", m.EnvFile, err)
			results = append(results, r)
			log.Printf("[error] %v", r.Err)
			continue
		}

		log.Printf("[ok] synced %q -> %q (%d keys)", m.VaultPath, m.EnvFile, len(secrets))
		results = append(results, r)
	}

	return results
}

// HasErrors returns true if any result contains an error.
func HasErrors(results []Result) bool {
	for _, r := range results {
		if r.Err != nil {
			return true
		}
	}
	return false
}

// ensure vault.Client satisfies VaultClient at compile time.
var _ VaultClient = (*vault.Client)(nil)
