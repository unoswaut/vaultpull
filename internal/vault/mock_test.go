package vault_test

import (
	"context"
	"errors"
	"testing"

	"github.com/user/vaultpull/internal/vault"
)

func TestMockClient_ReturnsConfiguredSecrets(t *testing.T) {
	mock := &vault.MockClient{
		Secrets: map[string]map[string]string{
			"app/dev": {"API_KEY": "abc123", "LOG_LEVEL": "debug"},
		},
	}

	got, err := mock.ReadSecrets(context.Background(), "app/dev")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got["API_KEY"] != "abc123" {
		t.Errorf("expected API_KEY=abc123, got %q", got["API_KEY"])
	}
	if len(mock.Calls) != 1 || mock.Calls[0] != "app/dev" {
		t.Errorf("unexpected calls: %v", mock.Calls)
	}
}

func TestMockClient_ReturnsCopiedMap(t *testing.T) {
	mock := &vault.MockClient{
		Secrets: map[string]map[string]string{
			"app/dev": {"KEY": "original"},
		},
	}

	result, _ := mock.ReadSecrets(context.Background(), "app/dev")
	result["KEY"] = "mutated"

	again, _ := mock.ReadSecrets(context.Background(), "app/dev")
	if again["KEY"] != "original" {
		t.Errorf("mock data was mutated; expected %q got %q", "original", again["KEY"])
	}
}

func TestMockClient_ReturnsError(t *testing.T) {
	expected := errors.New("vault unavailable")
	mock := &vault.MockClient{Err: expected}

	_, err := mock.ReadSecrets(context.Background(), "any/path")
	if !errors.Is(err, expected) {
		t.Errorf("expected %v, got %v", expected, err)
	}
}

func TestMockClient_MissingPath(t *testing.T) {
	mock := &vault.MockClient{
		Secrets: map[string]map[string]string{},
	}

	_, err := mock.ReadSecrets(context.Background(), "missing/path")
	if err == nil {
		t.Fatal("expected error for missing path")
	}
}
