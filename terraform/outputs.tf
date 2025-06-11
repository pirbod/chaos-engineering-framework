output "rg_name" {
  value       = azurerm_resource_group.rg.name
  description = "Name of the resource group"
}

output "lab_id" {
  value       = azurerm_dev_test_lab.lab.id
  description = "ID of the Azure DevTest Lab"
  sensitive   = true
}

output "vnet_id" {
  value       = azurerm_dev_test_virtual_network.vnet.id
  description = "ID of the DevTest Lab virtual network"
}
