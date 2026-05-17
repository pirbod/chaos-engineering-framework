package observability

import (
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Recorder struct {
	mu              sync.RWMutex
	requests        map[string]int
	durationSeconds map[string]float64
	riskScores      map[string]riskMetric
	aiConfidence    float64
}

type riskMetric struct {
	Score int
	Level string
}

func NewRecorder() *Recorder {
	return &Recorder{
		requests:        map[string]int{},
		durationSeconds: map[string]float64{},
		riskScores:      map[string]riskMetric{},
	}
}

func (r *Recorder) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		start := time.Now()
		sw := &statusWriter{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(sw, req)
		r.RecordRequest(req.Method, normalizePath(req.URL.Path), sw.status, time.Since(start))
	})
}

func (r *Recorder) RecordRequest(method, path string, status int, duration time.Duration) {
	r.mu.Lock()
	defer r.mu.Unlock()
	key := strings.Join([]string{method, path, strconv.Itoa(status)}, "|")
	r.requests[key]++
	r.durationSeconds[key] += duration.Seconds()
}

func (r *Recorder) RecordRisk(serviceID, level string, score int) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.riskScores[serviceID] = riskMetric{Score: score, Level: level}
}

func (r *Recorder) RecordAIConfidence(confidence float64) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.aiConfidence = confidence
}

func (r *Recorder) Handler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/plain; version=0.0.4; charset=utf-8")
		r.mu.RLock()
		defer r.mu.RUnlock()

		fmt.Fprintln(w, "# HELP platform_chaos_hub_http_requests_total Total HTTP requests by method, route and status.")
		fmt.Fprintln(w, "# TYPE platform_chaos_hub_http_requests_total counter")
		keys := sortedKeys(r.requests)
		for _, key := range keys {
			parts := strings.Split(key, "|")
			fmt.Fprintf(w, "platform_chaos_hub_http_requests_total{method=%s,path=%s,status=%s} %d\n", quote(parts[0]), quote(parts[1]), quote(parts[2]), r.requests[key])
		}

		fmt.Fprintln(w, "# HELP platform_chaos_hub_http_request_duration_seconds_sum Total HTTP request duration in seconds.")
		fmt.Fprintln(w, "# TYPE platform_chaos_hub_http_request_duration_seconds_sum counter")
		for _, key := range sortedKeysFloat(r.durationSeconds) {
			parts := strings.Split(key, "|")
			fmt.Fprintf(w, "platform_chaos_hub_http_request_duration_seconds_sum{method=%s,path=%s,status=%s} %.6f\n", quote(parts[0]), quote(parts[1]), quote(parts[2]), r.durationSeconds[key])
		}
		fmt.Fprintln(w, "# HELP platform_chaos_hub_http_request_duration_seconds_count Total HTTP request count for duration aggregation.")
		fmt.Fprintln(w, "# TYPE platform_chaos_hub_http_request_duration_seconds_count counter")
		for _, key := range keys {
			parts := strings.Split(key, "|")
			fmt.Fprintf(w, "platform_chaos_hub_http_request_duration_seconds_count{method=%s,path=%s,status=%s} %d\n", quote(parts[0]), quote(parts[1]), quote(parts[2]), r.requests[key])
		}

		fmt.Fprintln(w, "# HELP platform_chaos_hub_service_risk_score Latest calculated service risk score.")
		fmt.Fprintln(w, "# TYPE platform_chaos_hub_service_risk_score gauge")
		serviceIDs := make([]string, 0, len(r.riskScores))
		for serviceID := range r.riskScores {
			serviceIDs = append(serviceIDs, serviceID)
		}
		sort.Strings(serviceIDs)
		for _, serviceID := range serviceIDs {
			risk := r.riskScores[serviceID]
			fmt.Fprintf(w, "platform_chaos_hub_service_risk_score{service=%s,level=%s} %d\n", quote(serviceID), quote(risk.Level), risk.Score)
		}

		fmt.Fprintln(w, "# HELP platform_chaos_hub_ai_summary_confidence Latest AI incident summary confidence.")
		fmt.Fprintln(w, "# TYPE platform_chaos_hub_ai_summary_confidence gauge")
		fmt.Fprintf(w, "platform_chaos_hub_ai_summary_confidence %.2f\n", r.aiConfidence)
	})
}

type statusWriter struct {
	http.ResponseWriter
	status int
}

func (w *statusWriter) WriteHeader(status int) {
	w.status = status
	w.ResponseWriter.WriteHeader(status)
}

func sortedKeys(values map[string]int) []string {
	keys := make([]string, 0, len(values))
	for key := range values {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}

func sortedKeysFloat(values map[string]float64) []string {
	keys := make([]string, 0, len(values))
	for key := range values {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}

func normalizePath(path string) string {
	switch {
	case strings.HasPrefix(path, "/api/experiments/"):
		return "/api/experiments/{id}"
	case strings.HasPrefix(path, "/api/services/") && strings.HasSuffix(path, "/risk"):
		return "/api/services/{id}/risk"
	case strings.HasPrefix(path, "/api/services/"):
		return "/api/services/{id}"
	case strings.HasPrefix(path, "/api/runbooks/"):
		return "/api/runbooks/{id}"
	default:
		return path
	}
}

func quote(value string) string {
	return strconv.Quote(value)
}
