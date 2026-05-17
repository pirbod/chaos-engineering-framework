package tests

import (
	"bytes"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/pirbod/chaos-engineering-framework/backend/internal/ai"
	"github.com/pirbod/chaos-engineering-framework/backend/internal/api"
	"github.com/pirbod/chaos-engineering-framework/backend/internal/observability"
	"github.com/pirbod/chaos-engineering-framework/backend/internal/store"
)

func testServer() http.Handler {
	return api.NewRouter(
		store.NewInMemoryStore(),
		ai.LocalProvider{},
		observability.NewRecorder(),
		"test",
		slog.New(slog.NewTextHandler(io.Discard, nil)),
	)
}

func TestHealthz(t *testing.T) {
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	testServer().ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
	if !strings.Contains(rr.Body.String(), "healthy") {
		t.Fatalf("expected healthy response, got %s", rr.Body.String())
	}
}

func TestExperimentAndServiceAPIs(t *testing.T) {
	handler := testServer()
	for _, path := range []string{"/api/experiments", "/api/services", "/api/runbooks"} {
		rr := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, path, nil)
		handler.ServeHTTP(rr, req)
		if rr.Code != http.StatusOK {
			t.Fatalf("%s expected 200, got %d", path, rr.Code)
		}
		var response []map[string]any
		if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
			t.Fatalf("%s invalid json: %v", path, err)
		}
		if len(response) == 0 {
			t.Fatalf("%s expected non-empty response", path)
		}
	}
}

func TestServiceRiskAPI(t *testing.T) {
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/services/payments-api/risk", nil)
	testServer().ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
	var response struct {
		Score int    `json:"score"`
		Level string `json:"level"`
	}
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Fatalf("invalid json: %v", err)
	}
	if response.Score == 0 || response.Level == "" {
		t.Fatalf("expected populated risk response, got %+v", response)
	}
}

func TestAISummaryAPI(t *testing.T) {
	body := bytes.NewBufferString(`{
		"experimentName":"Network Latency",
		"affectedService":"payments-api",
		"metricsSnapshot":{"5xx_rate":"4.1%"},
		"logsOrEvents":["timeout waiting for provider response"],
		"riskScore":76
	}`)
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/ai/summarize", body)
	testServer().ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rr.Code, rr.Body.String())
	}
	if !strings.Contains(rr.Body.String(), "executiveSummary") {
		t.Fatalf("expected summary response, got %s", rr.Body.String())
	}
}

func TestApprovalRequestAPI(t *testing.T) {
	body := bytes.NewBufferString(`{
		"experimentId":"node-drain",
		"serviceId":"payments-api",
		"requestedBy":"dev@example.com",
		"riskScore":91,
		"riskLevel":"critical",
		"blastRadius":"critical",
		"runbookId":"failed-chaos-experiment",
		"rollbackThreshold":"abort if 5xx exceeds 2% for two minutes"
	}`)
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/approvals/requests", body)
	testServer().ServeHTTP(rr, req)

	if rr.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", rr.Code, rr.Body.String())
	}
	if !strings.Contains(rr.Body.String(), "requiredApprovers") {
		t.Fatalf("expected approval guardrails, got %s", rr.Body.String())
	}
}
