package approvals

import (
	"context"
	"testing"
)

func TestCriticalApprovalRequiresRunbookAndRollback(t *testing.T) {
	store := NewStore()
	_, err := store.Create(context.Background(), RequestInput{
		ExperimentID: "node-drain",
		ServiceID:    "payments-api",
		RequestedBy:  "dev@example.com",
		RiskScore:    91,
		BlastRadius:  "critical",
	})
	if err != ErrInvalid {
		t.Fatalf("expected ErrInvalid, got %v", err)
	}
}

func TestApprovalDecision(t *testing.T) {
	store := NewStore()
	request, err := store.Create(context.Background(), RequestInput{
		ExperimentID:      "node-drain",
		ServiceID:         "payments-api",
		RequestedBy:       "dev@example.com",
		RiskScore:         91,
		BlastRadius:       "critical",
		RunbookID:         "failed-chaos-experiment",
		RollbackThreshold: "abort if 5xx exceeds 2% for two minutes",
	})
	if err != nil {
		t.Fatalf("create approval: %v", err)
	}
	if len(request.RequiredApprovers) != 3 {
		t.Fatalf("expected 3 approvers, got %d", len(request.RequiredApprovers))
	}

	decided, err := store.Decide(context.Background(), request.ID, DecisionInput{
		Decision: "approved",
		Approver: "platform-oncall@example.com",
		Comment:  "approved for demo window",
	})
	if err != nil {
		t.Fatalf("decide approval: %v", err)
	}
	if decided.Status != "approved" {
		t.Fatalf("expected approved, got %q", decided.Status)
	}
}
