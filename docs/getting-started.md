# Getting Started

## Prerequisites

- Go 1.22+
- Node.js 20+
- npm 10+
- Docker and Docker Compose for container demo
- Optional: Terraform, Helm, Ansible, kubectl

## Local Demo

```bash
make setup
make dev
```

Open:

- Frontend: http://localhost:5173
- Backend: http://localhost:8080
- Metrics: http://localhost:8080/metrics

## Docker Demo

```bash
docker-compose up --build
```

## Validation

```bash
make test
make lint
make validate
```

If optional tools such as Terraform or Helm are missing, the local validation script reports the skip. CI runs those checks in tool-specific jobs.

## Useful API Calls

```bash
curl http://localhost:8080/healthz
curl http://localhost:8080/api/experiments
curl http://localhost:8080/api/services
curl http://localhost:8080/api/services/payments-api/risk
curl http://localhost:8080/api/approvals/requests
```
