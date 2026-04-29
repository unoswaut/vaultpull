package config

import (
	"errors"
	"os"

	"gopkg.in/yaml.v3"
)

// Mapping pairs a Vault secret path with a local env file destination.
type Mapping struct {
	VaultPath string `yaml:"vault_path"`
	EnvFile   string `yaml:"env_file"`
}

// Config holds all vaultpull runtime configuration.
type Config struct {
	VaultAddr string    `yaml:"vault_addr"`
	VaultToken string   `yaml:"vault_token"`
	Rotate    bool      `yaml:"rotate"`
	Mappings  []Mapping `yaml:"mappings"`
}

// Load reads configuration from the given YAML file path and applies
// environment variable overrides for sensitive fields.
func Load(path string) (*Config, error) {
	cfg := &Config{
		VaultAddr:  "http://127.0.0.1:8200",
		VaultToken: "",
		Rotate:     false,
	}

	if path != "" {
		data, err := os.ReadFile(path)
		if err != nil {
			return nil, err
		}
		if err := yaml.Unmarshal(data, cfg); err != nil {
			return nil, err
		}
	}

	if v := os.Getenv("VAULT_ADDR"); v != "" {
		cfg.VaultAddr = v
	}
	if v := os.Getenv("VAULT_TOKEN"); v != "" {
		cfg.VaultToken = v
	}

	if err := cfg.validate(); err != nil {
		return nil, err
	}
	return cfg, nil
}

func (c *Config) validate() error {
	if c.VaultAddr == "" {
		return errors.New("vault_addr is required")
	}
	if c.VaultToken == "" {
		return errors.New("vault_token is required")
	}
	for i, m := range c.Mappings {
		if m.VaultPath == "" {
			return errors.New("mapping missing vault_path")
		}
		if m.EnvFile == "" {
			return errors.New("mapping missing env_file")
		}
		_ = i
	}
	return nil
}
