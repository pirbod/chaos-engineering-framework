# GitLab Integration

GitLab support is included for teams using GitLab CI or GitLab-hosted services.

## Included Assets

- `gitlab/.gitlab-ci.example.yml`
- backend integration status placeholder
- GitLab runner degradation runbook

## Optional Environment Variables

- `GITLAB_BASE_URL`
- `GITLAB_TOKEN`

When configured, the backend can be extended to call GitLab APIs for project pipeline health, runner queue time and failed job trends.

## Chaos Validation in GitLab

Use the example pipeline to run:

- Go tests
- frontend type checks and tests
- Terraform format
- Helm lint/template
- Trivy filesystem scan
- container build
