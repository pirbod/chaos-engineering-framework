output "resource_group_name" {
  description = "Name of the DevTest Lab resource group."
  value       = azurerm_resource_group.this.name
}

output "lab_id" {
  description = "Azure DevTest Lab ID."
  value       = azurerm_dev_test_lab.this.id
}

output "virtual_network_id" {
  description = "DevTest Lab virtual network ID."
  value       = azurerm_dev_test_virtual_network.this.id
}
