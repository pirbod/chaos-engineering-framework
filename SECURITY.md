# Security Policy

## Reporting a Vulnerability

Please do not open public issues for suspected vulnerabilities. Email the repository owner with:

- affected component
- steps to reproduce
- expected impact
- suggested remediation if known

The project does not require real cloud, Vault, GitLab or Datadog credentials for the local demo. Report any accidental secret exposure immediately so keys can be revoked and history can be reviewed.

## Security Model

- Local demo data is embedded and non-sensitive.
- External integrations are optional and disabled unless environment variables are configured.
- Secrets belong in Vault, CI protected variables or local `.env` files that are never committed.
- Kubernetes deployments use a dedicated service account and network policy by default.
