package store

import (
	"context"
	"testing"
)

func TestSampleCatalogLoads(t *testing.T) {
	s := NewInMemoryStore()
	experiments, err := s.ListExperiments(context.Background())
	if err != nil {
		t.Fatalf("list experiments: %v", err)
	}
	services, err := s.ListServices(context.Background())
	if err != nil {
		t.Fatalf("list services: %v", err)
	}
	runbooks, err := s.ListRunbooks(context.Background())
	if err != nil {
		t.Fatalf("list runbooks: %v", err)
	}

	if len(experiments) < 7 {
		t.Fatalf("expected at least 7 experiments, got %d", len(experiments))
	}
	if len(services) < 5 {
		t.Fatalf("expected at least 5 services, got %d", len(services))
	}
	if len(runbooks) < 5 {
		t.Fatalf("expected at least 5 runbooks, got %d", len(runbooks))
	}
}

func TestMissingExperimentReturnsNotFound(t *testing.T) {
	_, err := NewInMemoryStore().GetExperiment(context.Background(), "missing")
	if err != ErrNotFound {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}
