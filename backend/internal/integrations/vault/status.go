package vault

import (
	"os"
	"time"

	"github.com/pirbod/chaos-engineering-framework/backend/internal/catalog"
)

func Status() catalog.IntegrationStatus {
	status := "unconfigured"
	details := "Set VAULT_ADDR and VAULT_TOKEN to enable connectivity checks in non-demo environments."
	if os.Getenv("VAULT_ADDR") != "" && os.Getenv("VAULT_TOKEN") != "" {
		status = "configured"
		details = "Vault environment variables are present. The demo does not make live secret reads by default."
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
