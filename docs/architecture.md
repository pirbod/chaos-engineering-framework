# Architecture

Platform Chaos Reliability Hub is a compact platform service around GitOps-managed chaos experiments.

```mermaid
flowchart TB
  subgraph Developer Experience
    Backstage["Backstage catalog and templates"]
    Frontend["React dashboard"]
  end
  subgraph API
    Go["Go HTTP API"]
    Store["Pluggable store interface"]
    Scoring["Risk scoring"]
    AI["AI summary provider"]
    Integrations["Integration status"]
  end
  subgraph Delivery
    Git["Git repository"]
    CI["GitHub Actions / GitLab CI"]
    Argo["Argo CD"]
    Helm["Helm chart"]
  end
  subgraph Runtime
    K8s["Kubernetes"]
    Chaos["Litmus/Chaos Mesh experiments"]
    Prom["Prometheus"]
    Grafana["Grafana"]
    Datadog["Datadog pattern"]
    Vault["Vault pattern"]
  end

  Backstage --> Frontend --> Go
  Go --> Store
  Go --> Scoring
  Go --> AI
  Go --> Integrations
  Go --> Prom
  Git --> CI --> Argo --> Helm --> K8s --> Chaos
  Prom --> Grafana
  Go --> Datadog
  Go --> Vault
```

## Design Choices

- The backend uses local in-memory sample data by default so the demo needs no database.
- The store is an interface so SQLite, Postgres or a catalog service can be added later.
- Risk scoring is deterministic and explainable rather than opaque.
- AI summaries are local and deterministic by default, with OpenAI-compatible configuration placeholders.
- Backstage is the intended single entry point, but the React UI can run standalone.

## Request Flow

```mermaid
sequenceDiagram
  participant User
  participant UI as React UI
  participant API as Go API
  participant Store as Catalog Store
  participant Metrics as Prometheus Metrics

  User->>UI: Open service risk page
  UI->>API: GET /api/services
  API->>Store: List services
  Store-->>API: Services
  UI->>API: GET /api/services/payments-api/risk
  API->>Store: Load service and experiments
  API->>API: Calculate deterministic risk
  API->>Metrics: Record risk score
  API-->>UI: Score, factors, actions
```
