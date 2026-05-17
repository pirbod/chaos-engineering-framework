package backstage

import (
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/pirbod/chaos-engineering-framework/backend/internal/catalog"
)

func Status() catalog.IntegrationStatus {
	status := "catalog-ready"
	details := "Backstage catalog descriptors and software templates are included in the repository."
	if os.Getenv("BACKSTAGE_BASE_URL") != "" {
		status = "configured"
		details = "Backstage base URL is configured and reachable."
		if err := probeBackstage(os.Getenv("BACKSTAGE_BASE_URL")); err != nil {
			status = "degraded"
			details = "Backstage is configured, but health probe failed: " + err.Error()
		}
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

func probeBackstage(baseURL string) error {
	req, err := http.NewRequest(http.MethodGet, strings.TrimRight(baseURL, "/")+"/api/catalog/entities?filter=kind=component", nil)
	if err != nil {
		return err
	}
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
