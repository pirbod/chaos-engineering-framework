# Runbook: GitLab Runner Degradation

## Trigger

Pipeline queue time, executor failures or runner latency exceed platform SLO.

## Immediate Checks

1. Check queued jobs and runner online count.
2. Review executor logs for capacity or network failures.
3. Compare queue time against the CI platform SLO.
4. Identify projects blocked by release windows.

## Remediation

1. Route critical projects to a healthy runner group.
2. Scale runner workers if capacity-driven.
3. Pause non-critical scheduled pipelines.
4. Retry failed jobs after queue time recovers.

## Rollback

Restore autoscaling settings and remove temporary runner overrides.

## Follow-Up

Record queue time trend and add a chaos test for the dominant failure mode.
