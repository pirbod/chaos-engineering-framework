variable "prefix" {
  description = "Prefix for Azure resource names."
  type        = string
}

variable "environment" {
  description = "Environment name."
  type        = string
}

variable "location" {
  description = "Azure region."
  type        = string
}

variable "resource_group_name" {
  description = "Resource group for observability resources."
  type        = string
}

variable "retention_in_days" {
  description = "Log Analytics retention period."
  type        = number
  default     = 30
}

variable "alert_email" {
  description = "Optional platform on-call email receiver."
  type        = string
  default     = ""
}

variable "tags" {
  description = "Tags applied to observability resources."
  type        = map(string)
  default     = {}
}
