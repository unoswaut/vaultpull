package config

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

// Config holds all configuration for vaultpull.
type Config struct {
	VaultAddr  string `mapstructure:"vault_addr"`
	VaultToken string `mapstructure:"vault_token"`
	VaultPath  string `mapstructure:"vault_path"`
	OutputFile string `mapstructure:"output_file"`
	Rotate     bool   `mapstructure:"rotate"`
}

// Load reads configuration from a file and environment variables.
// Environment variables take precedence over the config file.
func Load(cfgFile string) (*Config, error) {
	v := viper.New()

	v.SetDefault("vault_addr", "http://127.0.0.1:8200")
	v.SetDefault("output_file", ".env")
	v.SetDefault("rotate", false)

	v.SetEnvPrefix("VAULTPULL")
	v.AutomaticEnv()

	// Allow VAULT_ADDR and VAULT_TOKEN from standard Vault env vars.
	_ = v.BindEnv("vault_addr", "VAULT_ADDR")
	_ = v.BindEnv("vault_token", "VAULT_TOKEN")

	if cfgFile != "" {
		v.SetConfigFile(cfgFile)
	} else {
		v.SetConfigName(".vaultpull")
		v.SetConfigType("yaml")
		v.AddConfigPath(".")
		v.AddConfigPath(os.Getenv("HOME"))
	}

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("reading config file: %w", err)
		}
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("unmarshalling config: %w", err)
	}

	if err := cfg.validate(); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func (c *Config) validate() error {
	if c.VaultAddr == "" {
		return fmt.Errorf("vault_addr must not be empty")
	}
	if c.VaultToken == "" {
		return fmt.Errorf("vault_token must not be empty (set VAULT_TOKEN or vault_token in config)")
	}
	if c.VaultPath == "" {
		return fmt.Errorf("vault_path must not be empty")
	}
	return nil
}
