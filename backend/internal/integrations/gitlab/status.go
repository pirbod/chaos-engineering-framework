package gitlab

import (
	"os"
	"time"

	"github.com/pirbod/chaos-engineering-framework/backend/internal/catalog"
)

func Status() catalog.IntegrationStatus {
	status := "unconfigured"
	details := "Set GITLAB_BASE_URL and GITLAB_TOKEN to check project pipeline health from GitLab."
	if os.Getenv("GITLAB_BASE_URL") != "" && os.Getenv("GITLAB_TOKEN") != "" {
		status = "configured"
		details = "GitLab configuration is present. Live calls are intentionally disabled in the local demo."
	}
	return catalog.IntegrationStatus{
		ID:          "gitlab",
		Name:        "GitLab",
		Status:      status,
		Category:    "ci-cd",
		Details:     details,
		Signals:     []string{"pipeline success rate", "runner queue time", "failed job trend"},
		DocsURL:     "/docs/gitlab-integration.md",
		LastChecked: time.Now().UTC().Format(time.RFC3339),
	}
}
