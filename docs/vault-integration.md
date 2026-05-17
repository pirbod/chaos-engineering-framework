# Vault Integration

Vault is optional for the local demo. In real environments it should provide scoped secrets for platform integrations and workloads.

## Environment Variables

- `VAULT_ADDR`
- `VAULT_TOKEN`

The backend reports Vault as configured when both are present, but it does not read live secrets in the local demo.

## Kubernetes Pattern

Use External Secrets or Vault Agent with namespace-scoped policies:

```yaml
apiVersion: external-secrets.io/v1beta1
kind: ExternalSecret
metadata:
  name: platform-chaos-hub
spec:
  refreshInterval: 1h
  secretStoreRef:
    name: vault-platform
    kind: ClusterSecretStore
  target:
    name: platform-chaos-hub-secrets
  data:
    - secretKey: datadog-api-key
      remoteRef:
        key: platform/chaos-hub/datadog
        property: api_key
```

## Security Guidance

- Use least privilege policies.
- Prefer short-lived tokens or Kubernetes auth.
- Rotate tokens through GitOps and verify consumers before revocation.
- Never print secrets in CI logs or AI summaries.
