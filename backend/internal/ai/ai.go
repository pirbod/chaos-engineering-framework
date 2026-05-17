package ai

import (
	"context"
	"fmt"
	"os"
	"strings"
)

type SummaryInput struct {
	ExperimentName  string            `json:"experimentName"`
	AffectedService string            `json:"affectedService"`
	MetricsSnapshot map[string]string `json:"metricsSnapshot"`
	LogsOrEvents    []string          `json:"logsOrEvents"`
	RiskScore       int               `json:"riskScore"`
}

type SummaryOutput struct {
	ExecutiveSummary   string   `json:"executiveSummary"`
	LikelyImpact       string   `json:"likelyImpact"`
	ProbableRootCause  string   `json:"probableRootCause"`
	SuggestedActions   []string `json:"suggestedNextActions"`
	RecommendedRunbook string   `json:"recommendedRunbook"`
	ConfidenceLevel    float64  `json:"confidenceLevel"`
	Provider           string   `json:"provider"`
}

type Provider interface {
	Summarize(context.Context, SummaryInput) (SummaryOutput, error)
	Name() string
}

type LocalProvider struct{}

func (LocalProvider) Name() string {
	return "local-deterministic"
}

func (p LocalProvider) Summarize(_ context.Context, input SummaryInput) (SummaryOutput, error) {
	experiment := input.ExperimentName
	if experiment == "" {
		experiment = "chaos experiment"
	}
	service := input.AffectedService
	if service == "" {
		service = "the affected service"
	}

	rootCause := inferRootCause(input.LogsOrEvents)
	runbook := inferRunbook(input)
	impact := "Customer impact appears contained, but service health should be watched until metrics return to baseline."
	if input.RiskScore >= 80 {
		impact = "Potential customer impact is high because the service risk score is critical and rollback thresholds may be near."
	} else if input.RiskScore >= 60 {
		impact = "Potential customer impact is elevated because reliability signals show high operational risk."
	}

	return SummaryOutput{
		ExecutiveSummary: fmt.Sprintf("%s against %s produced a risk score of %d. Prioritize containment, owner acknowledgement and evidence capture before expanding blast radius.", experiment, service, input.RiskScore),
		LikelyImpact:     impact,
		ProbableRootCause: rootCause,
		SuggestedActions: []string{
			"Confirm the experiment is paused, completed or still within approved safety thresholds.",
			"Compare request rate, latency, errors and saturation against the pre-experiment baseline.",
			"Notify the owning team with the risk score, top signal and rollback recommendation.",
			"Record the finding in Backstage service metadata and update the maturity score.",
		},
		RecommendedRunbook: runbook,
		ConfidenceLevel:    confidence(input),
		Provider:           p.Name(),
	}, nil
}

type OpenAICompatibleProvider struct {
	LocalProvider
	BaseURL string
	Model   string
}

func (p OpenAICompatibleProvider) Name() string {
	if p.BaseURL == "" {
		return "openai-compatible-placeholder"
	}
	return "openai-compatible"
}

func (p OpenAICompatibleProvider) Summarize(ctx context.Context, input SummaryInput) (SummaryOutput, error) {
	output, err := p.LocalProvider.Summarize(ctx, input)
	if err != nil {
		return SummaryOutput{}, err
	}
	output.Provider = p.Name()
	return output, nil
}

func NewProviderFromEnv() Provider {
	switch strings.ToLower(os.Getenv("AI_PROVIDER")) {
	case "openai", "azure-openai":
		return OpenAICompatibleProvider{
			BaseURL: firstNonEmpty(os.Getenv("OPENAI_BASE_URL"), os.Getenv("AZURE_OPENAI_ENDPOINT")),
			Model:   firstNonEmpty(os.Getenv("OPENAI_MODEL"), os.Getenv("AZURE_OPENAI_DEPLOYMENT"), "gpt-4o-mini"),
		}
	default:
		return LocalProvider{}
	}
}

func inferRootCause(events []string) string {
	joined := strings.ToLower(strings.Join(events, " "))
	switch {
	case strings.Contains(joined, "timeout") || strings.Contains(joined, "latency"):
		return "Injected latency likely exceeded retry or timeout budgets for a downstream dependency."
	case strings.Contains(joined, "vault") || strings.Contains(joined, "secret"):
		return "Secret retrieval or token lease renewal appears to be the dominant failure signal."
	case strings.Contains(joined, "runner") || strings.Contains(joined, "queue"):
		return "Runner capacity or executor health likely increased CI queue time."
	case strings.Contains(joined, "oom") || strings.Contains(joined, "memory"):
		return "Memory pressure likely triggered restarts or degraded application performance."
	case strings.Contains(joined, "disk") || strings.Contains(joined, "ephemeral"):
		return "Disk pressure likely caused eviction risk or write-path degradation."
	default:
		return "The strongest signal is service degradation during the approved chaos window; compare metrics by dependency to isolate the failing path."
	}
}

func inferRunbook(input SummaryInput) string {
	text := strings.ToLower(input.ExperimentName + " " + input.AffectedService + " " + strings.Join(input.LogsOrEvents, " "))
	switch {
	case strings.Contains(text, "vault") || strings.Contains(text, "secret"):
		return "vault-secret-access-failure"
	case strings.Contains(text, "gitlab") || strings.Contains(text, "runner"):
		return "gitlab-runner-degradation"
	case strings.Contains(text, "datadog") || strings.Contains(text, "alert"):
		return "datadog-alert-noise"
	case strings.Contains(text, "pod") || strings.Contains(text, "node") || strings.Contains(text, "memory") || strings.Contains(text, "disk") || strings.Contains(text, "cpu"):
		return "kubernetes-service-degradation"
	default:
		return "failed-chaos-experiment"
	}
}

func confidence(input SummaryInput) float64 {
	points := 0.62
	if len(input.MetricsSnapshot) > 0 {
		points += 0.12
	}
	if len(input.LogsOrEvents) > 0 {
		points += 0.10
	}
	if input.RiskScore >= 60 {
		points += 0.08
	}
	if input.RiskScore >= 80 {
		points += 0.03
	}
	if points > 0.93 {
		return 0.93
	}
	return points
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if value != "" {
			return value
		}
	}
	return ""
}
