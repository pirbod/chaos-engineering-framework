package integrations

import (
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/pirbod/chaos-engineering-framework/backend/internal/catalog"
	"github.com/pirbod/chaos-engineering-framework/backend/internal/integrations/backstage"
	"github.com/pirbod/chaos-engineering-framework/backend/internal/integrations/datadog"
	"github.com/pirbod/chaos-engineering-framework/backend/internal/integrations/gitlab"
	"github.com/pirbod/chaos-engineering-framework/backend/internal/integrations/vault"
)

func Statuses() []catalog.IntegrationStatus {
	return []catalog.IntegrationStatus{
		gitlab.Status(),
		vault.Status(),
		datadog.Status(),
		backstage.Status(),
		prometheusStatus(),
		grafanaStatus(),
	}
}

func prometheusStatus() catalog.IntegrationStatus {
	status := "local-ready"
	details := "The backend exposes /metrics in Prometheus text format."
	if os.Getenv("PROMETHEUS_URL") != "" {
		status = "configured"
		details = "Prometheus URL is configured and reachable."
		if err := probeURL(os.Getenv("PROMETHEUS_URL"), "/-/ready"); err != nil {
			status = "degraded"
			details = "Prometheus is configured, but readiness probe failed: " + err.Error()
		}
	}
	return catalog.IntegrationStatus{
		ID:          "prometheus",
		Name:        "Prometheus",
		Status:      status,
		Category:    "observability",
		Details:     details,
		Signals:     []string{"api request rate", "latency", "service risk", "AI confidence"},
		DocsURL:     "/docs/observability.md",
		LastChecked: time.Now().UTC().Format(time.RFC3339),
	}
}

func grafanaStatus() catalog.IntegrationStatus {
	status := "dashboard-ready"
	details := "A Grafana dashboard JSON is included under monitoring/grafana."
	if os.Getenv("GRAFANA_URL") != "" {
		status = "configured"
		details = "Grafana URL is configured and reachable."
		if err := probeURL(os.Getenv("GRAFANA_URL"), "/api/health"); err != nil {
			status = "degraded"
			details = "Grafana is configured, but health probe failed: " + err.Error()
		}
	}
	return catalog.IntegrationStatus{
		ID:          "grafana",
		Name:        "Grafana",
		Status:      status,
		Category:    "observability",
		Details:     details,
		Signals:     []string{"MTTD trend", "experiment success rate", "critical risk services"},
		DocsURL:     "/docs/observability.md",
		LastChecked: time.Now().UTC().Format(time.RFC3339),
	}
}

func probeURL(baseURL, path string) error {
	req, err := http.NewRequest(http.MethodGet, strings.TrimRight(baseURL, "/")+path, nil)
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
