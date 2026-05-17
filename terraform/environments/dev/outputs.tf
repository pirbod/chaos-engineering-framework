output "resource_group_name" {
  description = "Dev resource group."
  value       = module.azure_lab.resource_group_name
}

output "managed_identity_client_id" {
  description = "Managed identity client ID for workload identity demos."
  value       = module.managed_identities.client_id
}

output "log_analytics_workspace_id" {
  description = "Observability workspace ID."
  value       = module.aks_observability.log_analytics_workspace_id
}
