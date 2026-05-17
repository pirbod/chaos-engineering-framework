# Runbook: Kubernetes Service Degradation

## Trigger

Service latency, errors, restarts or resource pressure increase during an experiment.

## Immediate Checks

1. Check deployment rollout status.
2. Inspect pod restarts, events and readiness.
3. Review CPU, memory, disk and network signals.
4. Compare p95 latency and 5xx rate to baseline.

## Remediation

1. Pause or abort the experiment if thresholds are breached.
2. Scale replicas if capacity is constrained.
3. Remove unhealthy pods only after confirming replacement capacity.
4. Revert recent resource or HPA changes if they caused instability.

## Rollback

Revert the experiment manifest and restore previous workload settings.

## Follow-Up

Update pod disruption budgets, resource requests and experiment guardrails.
