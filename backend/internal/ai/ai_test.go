package ai

import (
	"context"
	"reflect"
	"testing"
)

func TestLocalProviderSummarizeIsDeterministic(t *testing.T) {
	input := SummaryInput{
		ExperimentName:  "Network Latency",
		AffectedService: "payments-api",
		MetricsSnapshot: map[string]string{"5xx_rate": "4.1%", "p95": "920ms"},
		LogsOrEvents:    []string{"timeout waiting for provider response"},
		RiskScore:       76,
	}

	provider := LocalProvider{}
	first, err := provider.Summarize(context.Background(), input)
	if err != nil {
		t.Fatalf("summarize failed: %v", err)
	}
	second, err := provider.Summarize(context.Background(), input)
	if err != nil {
		t.Fatalf("summarize failed: %v", err)
	}

	if !reflect.DeepEqual(first, second) {
		t.Fatalf("expected deterministic summary\nfirst=%+v\nsecond=%+v", first, second)
	}
	if first.RecommendedRunbook != "failed-chaos-experiment" {
		t.Fatalf("unexpected runbook %q", first.RecommendedRunbook)
	}
	if first.ConfidenceLevel <= 0.7 {
		t.Fatalf("expected useful confidence, got %.2f", first.ConfidenceLevel)
	}
}
