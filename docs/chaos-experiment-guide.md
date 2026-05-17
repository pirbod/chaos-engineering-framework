# Chaos Experiment Guide

## Experiment Lifecycle

1. Define the hypothesis.
2. Select the smallest blast radius that can prove it.
3. Confirm service owner, runbook, dashboard and rollback threshold.
4. Open a pull request with the manifest.
5. Let CI validate structure and GitOps deploy after approval.
6. Observe SLO, RED and saturation signals.
7. Capture outcome and maturity improvements.

## Guardrails

- Do not target production databases in the demo pattern.
- Do not run critical blast radius experiments without platform approval.
- Abort when rollback thresholds are breached.
- Keep experiment duration short and explicit.
- Prefer canary pods before namespace or node-level disruption.

## Example

```yaml
apiVersion: chaos.litmus.io/v1alpha1
kind: ChaosEngine
metadata:
  name: cpu-hog-checkout-api
spec:
  appinfo:
    appns: default
    applabel: app=checkout-api
    appkind: deployment
  experiments:
    - name: pod-cpu-hog
```
