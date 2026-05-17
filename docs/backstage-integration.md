# Backstage Integration

Backstage is the intended single developer entry point.

## Included Assets

- `backstage/catalog-info.yaml`
- `backstage/templates/onboard-service/template.yaml`
- `backstage/templates/add-experiment/template.yaml`
- `backstage/plugin-platform-chaos-hub/`

## Developer Experience

Backstage should expose:

- service owner and escalation metadata
- approved experiments
- runbook links
- SLO and dashboard links
- maturity score and risk factors
- approval requests for high and critical blast radius tests

The plugin skeleton is lightweight by design. Teams can embed it into an internal Backstage app and call the Go API.
