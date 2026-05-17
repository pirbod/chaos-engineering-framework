package backstage

import (
	"os"
	"time"

	"github.com/pirbod/chaos-engineering-framework/backend/internal/catalog"
)

func Status() catalog.IntegrationStatus {
	status := "catalog-ready"
	details := "Backstage catalog descriptors and software templates are included in the repository."
	if os.Getenv("BACKSTAGE_BASE_URL") != "" {
		status = "configured"
		details = "Backstage base URL is configured for linking service ownership, runbooks and dashboards."
	}
	return catalog.IntegrationStatus{
		ID:          "backstage",
		Name:        "Backstage",
		Status:      status,
		Category:    "developer-portal",
		Details:     details,
		Signals:     []string{"service owner", "runbook link", "SLO metadata", "maturity score"},
		DocsURL:     "/docs/backstage-integration.md",
		LastChecked: time.Now().UTC().Format(time.RFC3339),
	}
}
