output "resource_group_name" {
  description = "Demo resource group."
  value       = module.azure_lab.resource_group_name
}

output "lab_id" {
  description = "Demo DevTest Lab ID."
  value       = module.azure_lab.lab_id
}
