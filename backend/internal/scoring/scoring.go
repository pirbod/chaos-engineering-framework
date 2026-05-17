package scoring

import (
	"math"
	"sort"

	"github.com/pirbod/chaos-engineering-framework/backend/internal/catalog"
)

type Factor struct {
	Name           string  `json:"name"`
	Contribution   int     `json:"contribution"`
	Weight         float64 `json:"weight"`
	Recommendation string  `json:"recommendation"`
}

type Result struct {
	ServiceID          string   `json:"serviceId"`
	Score              int      `json:"score"`
	Level              string   `json:"level"`
	TopFactors         []Factor `json:"topContributingFactors"`
	RecommendedActions []string `json:"recommendedActions"`
}

func Calculate(service catalog.Service, experiments []catalog.Experiment) Result {
	factors := []Factor{
		{
			Name:           "Experiment blast radius",
			Contribution:   maxBlastRadius(experiments),
			Weight:         0.18,
			Recommendation: "Reduce experiment scope to a canary, namespace or single dependency before expanding.",
		},
		{
			Name:           "Service criticality",
			Contribution:   criticalityScore(service.Criticality),
			Weight:         0.18,
			Recommendation: "Require platform approval and an explicit rollback threshold for critical services.",
		},
		{
			Name:           "Recent error rate",
			Contribution:   clampInt(int(math.Round(service.RecentErrorRate*450)), 0, 22),
			Weight:         0.14,
			Recommendation: "Wait for error rate to return to baseline before running a disruptive experiment.",
		},
		{
			Name:           "Ownership completeness",
			Contribution:   clampInt(int(math.Round(float64(100-service.OwnershipCompleteness)*0.18)), 0, 18),
			Weight:         0.12,
			Recommendation: "Add owner, escalation contact and Backstage metadata before approval.",
		},
		{
			Name:           "Runbook availability",
			Contribution:   runbookScore(service),
			Weight:         0.12,
			Recommendation: "Attach a tested runbook with validation and rollback steps.",
		},
		{
			Name:           "Observability maturity",
			Contribution:   clampInt(int(math.Round(float64(100-service.ObservabilityMaturity)*0.16)), 0, 16),
			Weight:         0.12,
			Recommendation: "Add RED metrics, SLO burn alerts and an experiment-specific dashboard panel.",
		},
		{
			Name:           "Last experiment result",
			Contribution:   outcomeScore(service.LastExperimentOutcome),
			Weight:         0.10,
			Recommendation: "Close the last experiment finding before running a broader scenario.",
		},
		{
			Name:           "Security sensitivity",
			Contribution:   sensitivityScore(service.SecuritySensitivity),
			Weight:         0.14,
			Recommendation: "Review secret, identity and compliance blast radius before execution.",
		},
	}

	total := 0
	for _, factor := range factors {
		total += factor.Contribution
	}
	score := clampInt(total, 0, 100)
	sort.SliceStable(factors, func(i, j int) bool {
		return factors[i].Contribution > factors[j].Contribution
	})

	top := factors
	if len(top) > 4 {
		top = top[:4]
	}
	actions := make([]string, 0, len(top))
	for _, factor := range top {
		if factor.Contribution > 0 {
			actions = append(actions, factor.Recommendation)
		}
	}
	if len(actions) == 0 {
		actions = append(actions, "Service is ready for a low-blast-radius experiment with normal change controls.")
	}

	return Result{
		ServiceID:          service.ID,
		Score:              score,
		Level:              riskLevel(score),
		TopFactors:         top,
		RecommendedActions: actions,
	}
}

func maxBlastRadius(experiments []catalog.Experiment) int {
	max := 0
	for _, experiment := range experiments {
		if score := blastRadiusScore(experiment.BlastRadius); score > max {
			max = score
		}
	}
	return max
}

func blastRadiusScore(value string) int {
	switch value {
	case "critical":
		return 24
	case "high":
		return 18
	case "medium":
		return 11
	case "low":
		return 5
	default:
		return 8
	}
}

func criticalityScore(value string) int {
	switch value {
	case "critical":
		return 20
	case "high":
		return 15
	case "medium":
		return 9
	case "low":
		return 4
	default:
		return 8
	}
}

func sensitivityScore(value string) int {
	switch value {
	case "critical":
		return 15
	case "high":
		return 11
	case "medium":
		return 6
	case "low":
		return 2
	default:
		return 5
	}
}

func outcomeScore(value string) int {
	switch value {
	case "failed":
		return 12
	case "partial":
		return 7
	case "never-run":
		return 8
	case "passed":
		return 0
	default:
		return 5
	}
}

func runbookScore(service catalog.Service) int {
	if len(service.RunbookIDs) == 0 {
		return 14
	}
	return 0
}

func riskLevel(score int) string {
	switch {
	case score >= 80:
		return "critical"
	case score >= 60:
		return "high"
	case score >= 35:
		return "medium"
	default:
		return "low"
	}
}

func clampInt(value, min, max int) int {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}
