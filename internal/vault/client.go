package vault

import (
	"context"
	"fmt"
	"time"

	vaultapi "github.com/hashicorp/vault/api"
)

// Client wraps the Vault API client with helper methods.
type Client struct {
	api     *vaultapi.Client
	mount   string
	timeout time.Duration
}

// Config holds the parameters needed to create a Vault client.
type Config struct {
	Address string
	Token   string
	Mount   string
	Timeout time.Duration
}

// New creates a new authenticated Vault client.
func New(cfg Config) (*Client, error) {
	if cfg.Address == "" {
		return nil, fmt.Errorf("vault address must not be empty")
	}
	if cfg.Token == "" {
		return nil, fmt.Errorf("vault token must not be empty")
	}

	apiCfg := vaultapi.DefaultConfig()
	apiCfg.Address = cfg.Address

	if cfg.Timeout > 0 {
		apiCfg.Timeout = cfg.Timeout
	}

	c, err := vaultapi.NewClient(apiCfg)
	if err != nil {
		return nil, fmt.Errorf("creating vault api client: %w", err)
	}
	c.SetToken(cfg.Token)

	mount := cfg.Mount
	if mount == "" {
		mount = "secret"
	}

	return &Client{api: c, mount: mount, timeout: cfg.Timeout}, nil
}

// ReadSecrets fetches key/value pairs from the given secret path.
func (c *Client) ReadSecrets(ctx context.Context, path string) (map[string]string, error) {
	ctx, cancel := context.WithTimeout(ctx, max(c.timeout, 10*time.Second))
	defer cancel()

	secret, err := c.api.KVv2(c.mount).Get(ctx, path)
	if err != nil {
		return nil, fmt.Errorf("reading secret at %q: %w", path, err)
	}
	if secret == nil || secret.Data == nil {
		return nil, fmt.Errorf("no data found at path %q", path)
	}

	result := make(map[string]string, len(secret.Data))
	for k, v := range secret.Data {
		str, ok := v.(string)
		if !ok {
			return nil, fmt.Errorf("value for key %q is not a string", k)
		}
		result[k] = str
	}
	return result, nil
}

func max(a, b time.Duration) time.Duration {
	if a > b {
		return a
	}
	return b
}
