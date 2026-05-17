package store

import (
	"context"
	"errors"
	"sort"
	"sync"

	"github.com/pirbod/chaos-engineering-framework/backend/internal/catalog"
)

var ErrNotFound = errors.New("not found")

type Store interface {
	ListExperiments(context.Context) ([]catalog.Experiment, error)
	GetExperiment(context.Context, string) (catalog.Experiment, error)
	ListServices(context.Context) ([]catalog.Service, error)
	GetService(context.Context, string) (catalog.Service, error)
	ListRunbooks(context.Context) ([]catalog.Runbook, error)
	GetRunbook(context.Context, string) (catalog.Runbook, error)
}

type InMemoryStore struct {
	mu          sync.RWMutex
	experiments map[string]catalog.Experiment
	services    map[string]catalog.Service
	runbooks    map[string]catalog.Runbook
}

func NewInMemoryStore() *InMemoryStore {
	s := &InMemoryStore{
		experiments: map[string]catalog.Experiment{},
		services:    map[string]catalog.Service{},
		runbooks:    map[string]catalog.Runbook{},
	}
	for _, experiment := range sampleExperiments() {
		s.experiments[experiment.ID] = experiment
	}
	for _, service := range sampleServices() {
		s.services[service.ID] = service
	}
	for _, runbook := range sampleRunbooks() {
		s.runbooks[runbook.ID] = runbook
	}
	return s
}

func (s *InMemoryStore) ListExperiments(context.Context) ([]catalog.Experiment, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	experiments := make([]catalog.Experiment, 0, len(s.experiments))
	for _, experiment := range s.experiments {
		experiments = append(experiments, experiment)
	}
	sort.Slice(experiments, func(i, j int) bool {
		return experiments[i].ID < experiments[j].ID
	})
	return experiments, nil
}

func (s *InMemoryStore) GetExperiment(_ context.Context, id string) (catalog.Experiment, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	experiment, ok := s.experiments[id]
	if !ok {
		return catalog.Experiment{}, ErrNotFound
	}
	return experiment, nil
}

func (s *InMemoryStore) ListServices(context.Context) ([]catalog.Service, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	services := make([]catalog.Service, 0, len(s.services))
	for _, service := range s.services {
		services = append(services, service)
	}
	sort.Slice(services, func(i, j int) bool {
		return services[i].ID < services[j].ID
	})
	return services, nil
}

func (s *InMemoryStore) GetService(_ context.Context, id string) (catalog.Service, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	service, ok := s.services[id]
	if !ok {
		return catalog.Service{}, ErrNotFound
	}
	return service, nil
}

func (s *InMemoryStore) ListRunbooks(context.Context) ([]catalog.Runbook, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	runbooks := make([]catalog.Runbook, 0, len(s.runbooks))
	for _, runbook := range s.runbooks {
		runbooks = append(runbooks, runbook)
	}
	sort.Slice(runbooks, func(i, j int) bool {
		return runbooks[i].ID < runbooks[j].ID
	})
	return runbooks, nil
}

func (s *InMemoryStore) GetRunbook(_ context.Context, id string) (catalog.Runbook, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	runbook, ok := s.runbooks[id]
	if !ok {
		return catalog.Runbook{}, ErrNotFound
	}
	return runbook, nil
}
