# Contributing

Thanks for improving the Platform Chaos Reliability Hub.

## Local Workflow

```bash
make setup
make test
make lint
make validate
```

## Change Guidelines

- Keep the platform compact and developer-focused.
- Add tests for scoring, API behavior and UI logic when changing behavior.
- Keep cloud integrations optional for the demo.
- Do not commit secrets, `.env` files or Terraform state.
- Document new experiments with safety notes, runbook links and expected signals.

## Dependency Updates

Review dependency updates in small batches. For Go, run `go test ./...`. For frontend updates, run `npm audit`, `npm run lint`, `npm test` and `npm run build`.
