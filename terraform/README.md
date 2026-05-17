# Terraform Azure Reference Infrastructure

The Terraform in this repository is intentionally small and demo-friendly. It keeps the original Azure DevTest Lab idea, then adds managed identity and observability modules to show how the platform would be wired in a real Azure environment.

## Modules

- `modules/azure-lab`: resource group, Azure DevTest Lab and lab virtual network.
- `modules/managed-identities`: user-assigned managed identity for workload identity or GitOps controllers.
- `modules/aks-observability`: Log Analytics workspace and optional alert action group.

## Usage

```bash
cd terraform/environments/dev
cp terraform.tfvars.example terraform.tfvars
terraform init
terraform plan
```

No secrets should be committed. Use environment variables, Azure workload identity or a secure remote backend for real environments.

## Security Notes

- Prefer managed identity over client secrets for automation.
- Scope identities to the minimum resource group or namespace required.
- Store Terraform state in a locked remote backend for shared environments.
- Do not put Vault, Datadog, GitLab or Azure credentials in `.tfvars` files.
