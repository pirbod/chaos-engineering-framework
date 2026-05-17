output "log_analytics_workspace_id" {
  description = "Log Analytics workspace ID for AKS/container insights."
  value       = azurerm_log_analytics_workspace.this.id
}

output "action_group_id" {
  description = "Action group ID for platform alerts."
  value       = azurerm_monitor_action_group.platform.id
}
