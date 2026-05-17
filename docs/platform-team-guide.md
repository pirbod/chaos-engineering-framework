# Platform Team Guide

The platform team owns guardrails, reusable experiments, integration patterns and operational feedback loops.

## Weekly Operating Loop

1. Review high and critical service risk scores.
2. Pair with owners to close metadata, runbook and observability gaps.
3. Approve only experiments with clear blast radius and rollback thresholds.
4. Review MTTD trend and noisy alerts after each experiment.
5. Convert repeated findings into templates, docs or automated checks.
6. Use approval requests to capture decisions for high and critical blast radius tests.

## Approval Checklist

- Owner and escalation route are present.
- SLO and dashboard links are present.
- Runbook includes checks and rollback.
- Experiment is scoped to the smallest useful blast radius.
- Security-sensitive services have Vault and identity scope reviewed.
- Alerts route to the right team and do not page unrelated services.

## Platform Partner Behavior

The platform should help teams make decisions, not just block risky changes. Prefer specific next actions, reusable templates and pairing sessions over broad policy statements.
