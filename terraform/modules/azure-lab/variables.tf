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

variable "tags" {
  description = "Tags applied to all resources."
  type        = map(string)
  default     = {}
}
