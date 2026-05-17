package approvals

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"
)

var (
	ErrNotFound   = errors.New("approval request not found")
	ErrInvalid    = errors.New("invalid approval request")
	ErrTransition = errors.New("invalid approval transition")
)

type RequestInput struct {
	ExperimentID      string   `json:"experimentId"`
	ServiceID         string   `json:"serviceId"`
	RequestedBy       string   `json:"requestedBy"`
	RiskScore         int      `json:"riskScore"`
	RiskLevel         string   `json:"riskLevel"`
	BlastRadius       string   `json:"blastRadius"`
	RollbackThreshold string   `json:"rollbackThreshold"`
	RunbookID         string   `json:"runbookId"`
	EvidenceLinks     []string `json:"evidenceLinks"`
}

type DecisionInput struct {
	Decision string `json:"decision"`
	Approver string `json:"approver"`
	Comment  string `json:"comment"`
}

type Request struct {
	ID                string    `json:"id"`
	ExperimentID      string    `json:"experimentId"`
	ServiceID         string    `json:"serviceId"`
	RequestedBy       string    `json:"requestedBy"`
	RiskScore         int       `json:"riskScore"`
	RiskLevel         string    `json:"riskLevel"`
	BlastRadius       string    `json:"blastRadius"`
	RollbackThreshold string    `json:"rollbackThreshold"`
	RunbookID         string    `json:"runbookId"`
	EvidenceLinks     []string  `json:"evidenceLinks"`
	Status            string    `json:"status"`
	RequiredApprovers []string  `json:"requiredApprovers"`
	Guardrails        []string  `json:"guardrails"`
	DecisionBy        string    `json:"decisionBy,omitempty"`
	DecisionComment   string    `json:"decisionComment,omitempty"`
	CreatedAt         time.Time `json:"createdAt"`
	UpdatedAt         time.Time `json:"updatedAt"`
}

type Store struct {
	mu       sync.RWMutex
	requests map[string]Request
}

func NewStore() *Store {
	return &Store{requests: map[string]Request{}}
}

func (s *Store) Create(_ context.Context, input RequestInput) (Request, error) {
	if err := validateInput(input); err != nil {
		return Request{}, err
	}
	now := time.Now().UTC()
	request := Request{
		ID:                fmt.Sprintf("apr-%d", now.UnixNano()),
		ExperimentID:      strings.TrimSpace(input.ExperimentID),
		ServiceID:         strings.TrimSpace(input.ServiceID),
		RequestedBy:       strings.TrimSpace(input.RequestedBy),
		RiskScore:         input.RiskScore,
		RiskLevel:         normalizeLevel(input.RiskLevel, input.RiskScore),
		BlastRadius:       strings.TrimSpace(input.BlastRadius),
		RollbackThreshold: strings.TrimSpace(input.RollbackThreshold),
		RunbookID:         strings.TrimSpace(input.RunbookID),
		EvidenceLinks:     append([]string{}, input.EvidenceLinks...),
		Status:            "pending",
		RequiredApprovers: requiredApprovers(input),
		Guardrails:        guardrails(input),
		CreatedAt:         now,
		UpdatedAt:         now,
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.requests[request.ID] = request
	return request, nil
}

func (s *Store) List(context.Context) ([]Request, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	items := make([]Request, 0, len(s.requests))
	for _, request := range s.requests {
		items = append(items, request)
	}
	sort.Slice(items, func(i, j int) bool {
		return items[i].CreatedAt.After(items[j].CreatedAt)
	})
	return items, nil
}

func (s *Store) Decide(_ context.Context, id string, input DecisionInput) (Request, error) {
	decision := strings.ToLower(strings.TrimSpace(input.Decision))
	if decision != "approved" && decision != "rejected" {
		return Request{}, ErrInvalid
	}
	if strings.TrimSpace(input.Approver) == "" {
		return Request{}, ErrInvalid
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	request, ok := s.requests[id]
	if !ok {
		return Request{}, ErrNotFound
	}
	if request.Status != "pending" {
		return Request{}, ErrTransition
	}
	request.Status = decision
	request.DecisionBy = strings.TrimSpace(input.Approver)
	request.DecisionComment = strings.TrimSpace(input.Comment)
	request.UpdatedAt = time.Now().UTC()
	s.requests[id] = request
	return request, nil
}

func validateInput(input RequestInput) error {
	if strings.TrimSpace(input.ExperimentID) == "" || strings.TrimSpace(input.ServiceID) == "" || strings.TrimSpace(input.RequestedBy) == "" {
		return ErrInvalid
	}
	if input.RiskScore < 0 || input.RiskScore > 100 {
		return ErrInvalid
	}
	if input.RiskScore >= 80 || strings.EqualFold(input.RiskLevel, "critical") || strings.EqualFold(input.BlastRadius, "critical") {
		if strings.TrimSpace(input.RollbackThreshold) == "" || strings.TrimSpace(input.RunbookID) == "" {
			return ErrInvalid
		}
	}
	return nil
}

func requiredApprovers(input RequestInput) []string {
	if input.RiskScore >= 80 || strings.EqualFold(input.RiskLevel, "critical") || strings.EqualFold(input.BlastRadius, "critical") {
		return []string{"service-owner", "platform-oncall", "security-platform"}
	}
	if input.RiskScore >= 60 || strings.EqualFold(input.RiskLevel, "high") || strings.EqualFold(input.BlastRadius, "high") {
		return []string{"service-owner", "platform-oncall"}
	}
	return []string{"service-owner"}
}

func guardrails(input RequestInput) []string {
	items := []string{
		"Confirm owner acknowledgement before execution.",
		"Verify SLO dashboard and alert routing.",
	}
	if input.RiskScore >= 80 || strings.EqualFold(input.BlastRadius, "critical") {
		items = append(items,
			"Run only inside an approved change window.",
			"Document rollback threshold and abort command before approval.",
			"Review secret, identity and customer data blast radius.",
		)
	}
	if strings.EqualFold(input.BlastRadius, "high") {
		items = append(items, "Start with a canary or single dependency before expanding scope.")
	}
	return items
}

func normalizeLevel(level string, score int) string {
	level = strings.ToLower(strings.TrimSpace(level))
	if level != "" {
		return level
	}
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
