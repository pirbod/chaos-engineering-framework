package experiments

import "github.com/pirbod/chaos-engineering-framework/backend/internal/catalog"

func ForService(service catalog.Service, all []catalog.Experiment) []catalog.Experiment {
	allowed := map[string]bool{}
	for _, id := range service.ExperimentIDs {
		allowed[id] = true
	}
	matched := make([]catalog.Experiment, 0, len(service.ExperimentIDs))
	for _, experiment := range all {
		if allowed[experiment.ID] || experiment.Target == service.ID {
			matched = append(matched, experiment)
		}
	}
	return matched
}
