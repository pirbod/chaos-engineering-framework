# AI Operations

The AI feature is intentionally practical: it helps summarize operational evidence after or during an experiment.

## Local Provider

The default provider is deterministic and requires no API key. It uses experiment name, affected service, metrics, logs/events and risk score to produce:

- executive summary
- likely impact
- probable root cause
- next actions
- recommended runbook
- confidence level

## Optional Providers

Set `AI_PROVIDER=openai` or `AI_PROVIDER=azure-openai` plus the documented environment variables to prepare for a real provider. The current implementation keeps the local deterministic behavior as the safe demo fallback.

## Prompting Pattern

Real deployments should keep prompts short, structured and evidence-first. Avoid sending secrets, raw tokens or sensitive customer payloads to external providers.
