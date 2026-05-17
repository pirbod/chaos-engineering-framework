package datadog

import (
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/pirbod/chaos-engineering-framework/backend/internal/catalog"
)

func Status() catalog.IntegrationStatus {
	status := "unconfigured"
	details := "Set DD_API_KEY and DD_SITE to publish chaos events and risk tags to Datadog."
	if os.Getenv("DD_API_KEY") != "" {
		status = "configured"
		details = "Datadog API key is present. Live API key validation succeeded."
		if err := probeDatadog(os.Getenv("DD_SITE"), os.Getenv("DD_API_KEY")); err != nil {
			status = "degraded"
			details = "Datadog is configured, but API key validation failed: " + err.Error()
		}
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

func probeDatadog(site, apiKey string) error {
	if strings.TrimSpace(site) == "" {
		site = "datadoghq.com"
	}
	req, err := http.NewRequest(http.MethodGet, "https://api."+strings.TrimPrefix(site, "api.")+"/api/v1/validate", nil)
	if err != nil {
		return err
	}
	req.Header.Set("DD-API-KEY", apiKey)
	client := &http.Client{Timeout: 1500 * time.Millisecond}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return nil
	}
	return &statusError{code: resp.StatusCode}
}

type statusError struct {
	code int
}

func (e *statusError) Error() string {
	return http.StatusText(e.code)
}
