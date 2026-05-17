# Observability

The observability model exists to reduce Mean Time To Detect.

## Backend Metrics

The Go API exposes:

- `platform_chaos_hub_http_requests_total`
- `platform_chaos_hub_http_request_duration_seconds_sum`
- `platform_chaos_hub_http_request_duration_seconds_count`
- `platform_chaos_hub_service_risk_score`
- `platform_chaos_hub_ai_summary_confidence`

## Dashboards

`monitoring/grafana/platform-chaos-dashboard.json` includes:

- API request rate
- API latency
- experiment success rate
- experiment failure rate
- critical risk services
- AI summary confidence
- MTTD placeholder trend

## Alerts

`monitoring/alerts/platform-chaos-alerts.yaml` includes alert rules for high chaos failure rate, critical risk, missing metadata/runbook review and high MTTD.

## MTTD Workflow

1. Tag experiment events by service and risk level.
2. Route alerts to the owning team.
3. Use runbook links in the alert.
4. Generate an AI summary for handoff.
5. Review whether detection happened before customer impact.
