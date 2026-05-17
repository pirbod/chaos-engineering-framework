package scoring

import (
	"testing"

	"github.com/pirbod/chaos-engineering-framework/backend/internal/catalog"
)

func TestCalculateCriticalServiceRisk(t *testing.T) {
	service := catalog.Service{
		ID:                    "payments-api",
		Criticality:           "critical",
		RecentErrorRate:       0.04,
		OwnershipCompleteness: 80,
		RunbookIDs:            []string{"failed-chaos-experiment"},
		ObservabilityMaturity: 70,
		LastExperimentOutcome: "failed",
		SecuritySensitivity:   "critical",
	}
	experiments := []catalog.Experiment{{ID: "node-drain", BlastRadius: "critical"}}

	result := Calculate(service, experiments)

	if result.Score < 70 {
		t.Fatalf("expected high risk score, got %d", result.Score)
	}
	if result.Level != "high" && result.Level != "critical" {
		t.Fatalf("expected high or critical level, got %q", result.Level)
	}
	if len(result.TopFactors) == 0 {
		t.Fatal("expected contributing factors")
	}
}

func TestCalculateLowRiskService(t *testing.T) {
	service := catalog.Service{
		ID:                    "recommendations",
		Criticality:           "low",
		RecentErrorRate:       0.001,
		OwnershipCompleteness: 100,
		RunbookIDs:            []string{"kubernetes-service-degradation"},
		ObservabilityMaturity: 95,
		LastExperimentOutcome: "passed",
		SecuritySensitivity:   "low",
	}
	experiments := []catalog.Experiment{{ID: "memory-hog", BlastRadius: "low"}}

	result := Calculate(service, experiments)

	if result.Score >= 35 {
		t.Fatalf("expected low score, got %d", result.Score)
	}
	if result.Level != "low" {
		t.Fatalf("expected low level, got %q", result.Level)
	}
}
