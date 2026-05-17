# Runbook: Failed Chaos Experiment

## Trigger

An experiment fails its hypothesis, exceeds rollback threshold or expands beyond approved blast radius.

## Immediate Checks

1. Confirm experiment status in Argo CD or the chaos controller.
2. Check service SLO burn, error rate, latency and saturation.
3. Confirm the owning team is aware.
4. Capture affected namespace, deployment and experiment ID.

## Remediation

1. Pause or abort the experiment.
2. Revert the experiment manifest if GitOps continues to apply it.
3. Scale affected workloads back to the last healthy replica count.
4. Open an incident if customer impact is possible.

## Rollback

Revert the experiment pull request or suspend the Argo CD application until the finding is understood.

## Follow-Up

Update the risk score, runbook and safety notes before re-running.
