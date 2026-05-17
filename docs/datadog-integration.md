# Datadog Integration

Datadog is optional for the local demo. The repository includes monitor examples and event tag mapping.

## Environment Variables

- `DD_API_KEY`
- `DD_SITE`

## Event Tags

Map experiment events to:

- `service:<service-id>`
- `experiment_id:<id>`
- `risk_score:<score>`
- `risk_level:<low|medium|high|critical>`
- `experiment_outcome:<passed|failed|partial>`

## Monitor Examples

See `monitoring/datadog/monitors.json` for high failure rate and critical service risk monitor examples.

## Noise Reduction

Use experiment tags to group alerts, suppress approved windows and identify monitors that fired without action. Do not mute critical SLO burn alerts for customer-facing services.
