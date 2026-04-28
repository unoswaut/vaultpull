package vault_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/user/vaultpull/internal/vault"
)

func newFakeVault(t *testing.T, data map[string]interface{}) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		payload := map[string]interface{}{
			"data": map[string]interface{}{
				"data": data,
			},
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(payload)
	}))
}

func TestNew_MissingAddress(t *testing.T) {
	_, err := vault.New(vault.Config{Token: "tok"})
	if err == nil {
		t.Fatal("expected error for missing address")
	}
}

func TestNew_MissingToken(t *testing.T) {
	_, err := vault.New(vault.Config{Address: "http://localhost:8200"})
	if err == nil {
		t.Fatal("expected error for missing token")
	}
}

func TestNew_Success(t *testing.T) {
	srv := newFakeVault(t, nil)
	defer srv.Close()

	c, err := vault.New(vault.Config{
		Address: srv.URL,
		Token:   "test-token",
		Mount:   "secret",
		Timeout: 5 * time.Second,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c == nil {
		t.Fatal("expected non-nil client")
	}
}

func TestReadSecrets_ReturnsValues(t *testing.T) {
	srv := newFakeVault(t, map[string]interface{}{
		"DB_HOST": "localhost",
		"DB_PORT": "5432",
	})
	defer srv.Close()

	c, err := vault.New(vault.Config{
		Address: srv.URL,
		Token:   "test-token",
		Timeout: 5 * time.Second,
	})
	if err != nil {
		t.Fatalf("creating client: %v", err)
	}

	secrets, err := c.ReadSecrets(context.Background(), "myapp/prod")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if secrets["DB_HOST"] != "localhost" {
		t.Errorf("expected DB_HOST=localhost, got %q", secrets["DB_HOST"])
	}
	if secrets["DB_PORT"] != "5432" {
		t.Errorf("expected DB_PORT=5432, got %q", secrets["DB_PORT"])
	}
}
