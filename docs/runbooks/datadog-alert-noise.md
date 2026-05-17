# Runbook: Datadog Alert Noise

## Trigger

Datadog monitors page repeatedly during approved chaos windows without actionable signal.

## Immediate Checks

1. Group alerts by `service`, `experiment_id` and `risk_level`.
2. Identify monitors that fired without requiring action.
3. Confirm SLO burn and critical customer alerts still page.

## Remediation

1. Add chaos event tags to dashboards and monitors.
2. Tune thresholds only with SLO evidence.
3. Use short mute windows for approved experiments.
4. Keep critical risk alerts active.

## Rollback

Revert monitor threshold changes and remove temporary mute windows.

## Follow-Up

Update Datadog monitor examples and platform alert guidance.
