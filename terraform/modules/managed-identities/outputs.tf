output "principal_id" {
  description = "Managed identity principal ID."
  value       = azurerm_user_assigned_identity.platform.principal_id
}

output "client_id" {
  description = "Managed identity client ID."
  value       = azurerm_user_assigned_identity.platform.client_id
}

output "id" {
  description = "Managed identity resource ID."
  value       = azurerm_user_assigned_identity.platform.id
}
