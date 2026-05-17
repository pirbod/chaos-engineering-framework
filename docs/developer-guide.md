# Developer Guide

This platform helps product teams request and review chaos experiments without learning every platform implementation detail.

## Request an Experiment

1. Open the experiment catalog.
2. Filter by your service or failure mode.
3. Review blast radius and safety notes.
4. Confirm owner, SLO, dashboard and runbook links.
5. Open a pull request with the experiment manifest.

## Read the Risk Score

Risk is scored from 0 to 100:

- `low`: safe for normal change controls
- `medium`: run canary-first and review runbook coverage
- `high`: platform review recommended
- `critical`: platform approval required before execution

The score includes blast radius, criticality, error rate, ownership, runbooks, observability, last experiment result and security sensitivity.

## Reduce Your Score

- Add missing owner and escalation metadata.
- Attach a clear runbook with rollback.
- Improve dashboards and SLO alerts.
- Start with smaller blast radius.
- Close findings from failed or partial experiments.
