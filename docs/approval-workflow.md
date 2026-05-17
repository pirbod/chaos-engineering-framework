# Approval Workflow

Critical and high-risk experiments need lightweight approval before execution. The goal is not bureaucracy; it is making rollback, ownership and evidence explicit before blast radius increases.

## API

- `GET /api/approvals/requests`
- `POST /api/approvals/requests`
- `POST /api/approvals/requests/{id}/decision`

## Required Guardrails

Critical requests require:

- experiment ID
- service ID
- requester
- risk score
- rollback threshold
- runbook ID

The API derives required approvers:

- low/medium: service owner
- high: service owner and platform on-call
- critical: service owner, platform on-call and security platform

## Example

```bash
curl -s http://localhost:8080/api/approvals/requests \
  -H 'Content-Type: application/json' \
  -d '{
    "experimentId":"node-drain",
    "serviceId":"payments-api",
    "requestedBy":"dev@example.com",
    "riskScore":91,
    "riskLevel":"critical",
    "blastRadius":"critical",
    "runbookId":"failed-chaos-experiment",
    "rollbackThreshold":"abort if 5xx exceeds 2% for two minutes"
  }'
```

## Production Extension

Replace the in-memory approval store with SQLite/Postgres, wire identity from Backstage or an ingress auth proxy, and emit approval events to Datadog or an audit sink.
