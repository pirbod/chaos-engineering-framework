variable "prefix" {
  description = "Prefix for Azure resource names."
  type        = string
}

variable "environment" {
  description = "Environment name used for resource naming and tags."
  type        = string
}

variable "location" {
  description = "Azure region."
  type        = string
}

variable "address_space" {
  description = "DevTest Lab virtual network address space."
  type        = list(string)
  default     = ["10.1.0.0/16"]
}

variable "tags" {
  description = "Tags applied to all resources."
  type        = map(string)
  default     = {}
}
