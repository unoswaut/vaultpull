package vault

import (
	"context"
	"fmt"
)

// Reader is the interface satisfied by Client for reading secrets.
type Reader interface {
	ReadSecrets(ctx context.Context, path string) (map[string]string, error)
}

// MockClient is a test double for Reader.
type MockClient struct {
	// Secrets maps path -> key/value pairs to return.
	Secrets map[string]map[string]string
	// Err is returned for every call when set.
	Err error
	// Calls records paths that were requested.
	Calls []string
}

// ReadSecrets implements Reader.
func (m *MockClient) ReadSecrets(_ context.Context, path string) (map[string]string, error) {
	m.Calls = append(m.Calls, path)
	if m.Err != nil {
		return nil, m.Err
	}
	data, ok := m.Secrets[path]
	if !ok {
		return nil, fmt.Errorf("mock: no data for path %q", path)
	}
	// Return a copy to avoid mutation across calls.
	copy := make(map[string]string, len(data))
	for k, v := range data {
		copy[k] = v
	}
	return copy, nil
}
