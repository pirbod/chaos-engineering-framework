# Security Model

## Principles

- Local demo runs without cloud credentials.
- Integration credentials are optional and environment-driven.
- Secrets must live in Vault, CI protected variables or local ignored files.
- Managed identity is preferred for Azure automation.
- Critical experiments require identity, secret and data sensitivity review.

## Controls

- `.env.example` documents configuration without values.
- `.gitignore` excludes `.env`, Terraform state and generated artifacts.
- Helm chart includes service account and network policy templates.
- Security workflow includes secret scan and filesystem vulnerability scanning.
- Vault docs define least-privilege token and External Secrets patterns.

## Threats to Watch

- Over-broad Vault policy used during an experiment.
- CI token with write access exposed through logs.
- Datadog event payloads containing sensitive customer data.
- Chaos experiment scoped wider than the approved namespace or label selector.
