package gitlab

import (
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/pirbod/chaos-engineering-framework/backend/internal/catalog"
)

func Status() catalog.IntegrationStatus {
	status := "unconfigured"
	details := "Set GITLAB_BASE_URL and GITLAB_TOKEN to check project pipeline health from GitLab."
	if os.Getenv("GITLAB_BASE_URL") != "" && os.Getenv("GITLAB_TOKEN") != "" {
		status = "configured"
		details = "GitLab configuration is present. Live version probe succeeded."
		if err := probeGitLab(os.Getenv("GITLAB_BASE_URL"), os.Getenv("GITLAB_TOKEN")); err != nil {
			status = "degraded"
			details = "GitLab is configured, but version probe failed: " + err.Error()
		}
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

func probeGitLab(baseURL, token string) error {
	req, err := http.NewRequest(http.MethodGet, strings.TrimRight(baseURL, "/")+"/api/v4/version", nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+token)
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
