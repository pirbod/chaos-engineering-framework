package api

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/pirbod/chaos-engineering-framework/backend/internal/ai"
	"github.com/pirbod/chaos-engineering-framework/backend/internal/experiments"
	"github.com/pirbod/chaos-engineering-framework/backend/internal/integrations"
	"github.com/pirbod/chaos-engineering-framework/backend/internal/observability"
	"github.com/pirbod/chaos-engineering-framework/backend/internal/scoring"
	"github.com/pirbod/chaos-engineering-framework/backend/internal/store"
)

type Server struct {
	store    store.Store
	ai       ai.Provider
	metrics  *observability.Recorder
	version  string
	logger   *slog.Logger
}

func NewRouter(data store.Store, provider ai.Provider, metrics *observability.Recorder, version string, logger *slog.Logger) http.Handler {
	if version == "" {
		version = "dev"
	}
	if logger == nil {
		logger = slog.Default()
	}
	s := &Server{
		store:   data,
		ai:      provider,
		metrics: metrics,
		version: version,
		logger:  logger,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", s.handleIndex)
	mux.HandleFunc("/healthz", s.handleHealth)
	mux.HandleFunc("/readyz", s.handleReady)
	mux.HandleFunc("/version", s.handleVersion)
	mux.Handle("/metrics", metrics.Handler())
	mux.HandleFunc("/api/experiments", s.handleExperiments)
	mux.HandleFunc("/api/experiments/", s.handleExperimentDetail)
	mux.HandleFunc("/api/services", s.handleServices)
	mux.HandleFunc("/api/services/", s.handleServiceRoutes)
	mux.HandleFunc("/api/runbooks", s.handleRunbooks)
	mux.HandleFunc("/api/runbooks/", s.handleRunbookDetail)
	mux.HandleFunc("/api/ai/summarize", s.handleAISummary)
	mux.HandleFunc("/api/integrations", s.handleIntegrations)

	return withCORS(withLogging(logger, metrics.Middleware(mux)))
}

func (s *Server) handleIndex(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		writeError(w, http.StatusNotFound, "route not found")
		return
	}
	if !requireMethod(w, r, http.MethodGet) {
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{
		"name":    "Platform Chaos Reliability Hub",
		"version": s.version,
		"status":  "ready",
	})
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	if !requireMethod(w, r, http.MethodGet) {
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "healthy"})
}

func (s *Server) handleReady(w http.ResponseWriter, r *http.Request) {
	if !requireMethod(w, r, http.MethodGet) {
		return
	}
	if _, err := s.store.ListExperiments(r.Context()); err != nil {
		writeError(w, http.StatusServiceUnavailable, "catalog is not ready")
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "ready"})
}

func (s *Server) handleVersion(w http.ResponseWriter, r *http.Request) {
	if !requireMethod(w, r, http.MethodGet) {
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{
		"version":    s.version,
		"aiProvider": s.ai.Name(),
	})
}

func (s *Server) handleExperiments(w http.ResponseWriter, r *http.Request) {
	if !requireMethod(w, r, http.MethodGet) {
		return
	}
	items, err := s.store.ListExperiments(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to list experiments")
		return
	}
	writeJSON(w, http.StatusOK, items)
}

func (s *Server) handleExperimentDetail(w http.ResponseWriter, r *http.Request) {
	if !requireMethod(w, r, http.MethodGet) {
		return
	}
	id := strings.TrimPrefix(r.URL.Path, "/api/experiments/")
	item, err := s.store.GetExperiment(r.Context(), id)
	if err != nil {
		writeStoreError(w, err, "experiment not found")
		return
	}
	writeJSON(w, http.StatusOK, item)
}

func (s *Server) handleServices(w http.ResponseWriter, r *http.Request) {
	if !requireMethod(w, r, http.MethodGet) {
		return
	}
	items, err := s.store.ListServices(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to list services")
		return
	}
	writeJSON(w, http.StatusOK, items)
}

func (s *Server) handleServiceRoutes(w http.ResponseWriter, r *http.Request) {
	if !requireMethod(w, r, http.MethodGet) {
		return
	}
	path := strings.TrimPrefix(r.URL.Path, "/api/services/")
	if strings.HasSuffix(path, "/risk") {
		id := strings.TrimSuffix(path, "/risk")
		id = strings.TrimSuffix(id, "/")
		s.handleServiceRisk(w, r, id)
		return
	}
	service, err := s.store.GetService(r.Context(), path)
	if err != nil {
		writeStoreError(w, err, "service not found")
		return
	}
	writeJSON(w, http.StatusOK, service)
}

func (s *Server) handleServiceRisk(w http.ResponseWriter, r *http.Request, id string) {
	service, err := s.store.GetService(r.Context(), id)
	if err != nil {
		writeStoreError(w, err, "service not found")
		return
	}
	allExperiments, err := s.store.ListExperiments(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to list experiments")
		return
	}
	result := scoring.Calculate(service, experiments.ForService(service, allExperiments))
	s.metrics.RecordRisk(result.ServiceID, result.Level, result.Score)
	writeJSON(w, http.StatusOK, result)
}

func (s *Server) handleRunbooks(w http.ResponseWriter, r *http.Request) {
	if !requireMethod(w, r, http.MethodGet) {
		return
	}
	items, err := s.store.ListRunbooks(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to list runbooks")
		return
	}
	writeJSON(w, http.StatusOK, items)
}

func (s *Server) handleRunbookDetail(w http.ResponseWriter, r *http.Request) {
	if !requireMethod(w, r, http.MethodGet) {
		return
	}
	id := strings.TrimPrefix(r.URL.Path, "/api/runbooks/")
	item, err := s.store.GetRunbook(r.Context(), id)
	if err != nil {
		writeStoreError(w, err, "runbook not found")
		return
	}
	writeJSON(w, http.StatusOK, item)
}

func (s *Server) handleAISummary(w http.ResponseWriter, r *http.Request) {
	if !requireMethod(w, r, http.MethodPost) {
		return
	}
	defer r.Body.Close()
	var input ai.SummaryInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "invalid summary request body")
		return
	}
	output, err := s.ai.Summarize(r.Context(), input)
	if err != nil {
		writeError(w, http.StatusBadGateway, "failed to generate incident summary")
		return
	}
	s.metrics.RecordAIConfidence(output.ConfidenceLevel)
	writeJSON(w, http.StatusOK, output)
}

func (s *Server) handleIntegrations(w http.ResponseWriter, r *http.Request) {
	if !requireMethod(w, r, http.MethodGet) {
		return
	}
	writeJSON(w, http.StatusOK, integrations.Statuses())
}

func requireMethod(w http.ResponseWriter, r *http.Request, allowed string) bool {
	if r.Method == allowed {
		return true
	}
	w.Header().Set("Allow", allowed)
	writeError(w, http.StatusMethodNotAllowed, "method not allowed")
	return false
}

func writeStoreError(w http.ResponseWriter, err error, fallback string) {
	if errors.Is(err, store.ErrNotFound) {
		writeError(w, http.StatusNotFound, fallback)
		return
	}
	writeError(w, http.StatusInternalServerError, "catalog lookup failed")
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(value); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
	}
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}

func withCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET,POST,OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type,Authorization,X-Request-ID")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func withLogging(logger *slog.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		traceID := r.Header.Get("X-Request-ID")
		if traceID == "" {
			traceID = time.Now().UTC().Format("20060102150405.000000000")
		}
		w.Header().Set("X-Request-ID", traceID)
		sw := &loggingResponseWriter{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(sw, r)
		logger.Info("http_request",
			"method", r.Method,
			"path", r.URL.Path,
			"status", sw.status,
			"duration_ms", time.Since(start).Milliseconds(),
			"trace_id", traceID,
		)
	})
}

type loggingResponseWriter struct {
	http.ResponseWriter
	status int
}

func (w *loggingResponseWriter) WriteHeader(status int) {
	w.status = status
	w.ResponseWriter.WriteHeader(status)
}
