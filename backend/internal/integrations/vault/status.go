package vault

import (
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/pirbod/chaos-engineering-framework/backend/internal/catalog"
)

func Status() catalog.IntegrationStatus {
	status := "unconfigured"
	details := "Set VAULT_ADDR and VAULT_TOKEN to enable connectivity checks in non-demo environments."
	if os.Getenv("VAULT_ADDR") != "" && os.Getenv("VAULT_TOKEN") != "" {
		status = "configured"
		details = "Vault environment variables are present. Live health probe succeeded."
		if err := probeVault(os.Getenv("VAULT_ADDR"), os.Getenv("VAULT_TOKEN")); err != nil {
			status = "degraded"
			details = "Vault is configured, but health probe failed: " + err.Error()
		}
	}
	return catalog.IntegrationStatus{
		ID:          "vault",
		Name:        "Vault",
		Status:      status,
		Category:    "secrets",
		Details:     details,
		Signals:     []string{"secret lease health", "policy scope", "external secret sync"},
		DocsURL:     "/docs/vault-integration.md",
		LastChecked: time.Now().UTC().Format(time.RFC3339),
	}
}

func probeVault(addr, token string) error {
	req, err := http.NewRequest(http.MethodGet, strings.TrimRight(addr, "/")+"/v1/sys/health", nil)
	if err != nil {
		return err
	}
	req.Header.Set("X-Vault-Token", token)
	client := &http.Client{Timeout: 1500 * time.Millisecond}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	switch resp.StatusCode {
	case http.StatusOK, http.StatusTooManyRequests, 472, 473:
		return nil
	default:
		return &statusError{code: resp.StatusCode}
	}
}

type statusError struct {
	code int
}

func (e *statusError) Error() string {
	return http.StatusText(e.code)
}
