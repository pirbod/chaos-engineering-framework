# Runbook: Vault Secret Access Failure

## Trigger

Applications cannot read or refresh secrets from Vault or External Secrets.

## Immediate Checks

1. Check Vault health, seal status and active leader.
2. Validate Kubernetes auth role and service account binding.
3. Inspect ExternalSecret reconciliation errors.
4. Confirm no broad token was introduced during remediation.

## Remediation

1. Restore the previous Vault policy if a recent change caused denial.
2. Renew or rotate only the affected token or role.
3. Restart External Secrets reconciliation if needed.
4. Verify application pods read the expected secret version.

## Rollback

Revert the policy or ExternalSecret manifest through GitOps.

## Follow-Up

Reduce policy scope and add lease/error alerts if missing.
