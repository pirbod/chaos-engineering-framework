package integrations

import (
	"os"
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
		details = "Prometheus URL is configured for future query integration."
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
		details = "Grafana URL is configured for deep links from the developer portal."
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
