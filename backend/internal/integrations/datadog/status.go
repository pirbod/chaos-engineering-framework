package datadog

import (
	"os"
	"time"

	"github.com/pirbod/chaos-engineering-framework/backend/internal/catalog"
)

func Status() catalog.IntegrationStatus {
	status := "unconfigured"
	details := "Set DD_API_KEY and DD_SITE to publish chaos events and risk tags to Datadog."
	if os.Getenv("DD_API_KEY") != "" {
		status = "configured"
		details = "Datadog API key is present. The demo exposes event tag mapping without sending data."
	}
	return catalog.IntegrationStatus{
		ID:          "datadog",
		Name:        "Datadog",
		Status:      status,
		Category:    "observability",
		Details:     details,
		Signals:     []string{"chaos events", "monitor noise", "SLO burn", "risk tags"},
		DocsURL:     "/docs/datadog-integration.md",
		LastChecked: time.Now().UTC().Format(time.RFC3339),
	}
}
