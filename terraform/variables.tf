
variable "subscription_id" {
  type        = string
  description = "Azure Subscription ID for provider"
}

variable "prefix" {
  type        = string
  description = "Prefix for all resource names"
}

variable "environment" {
  type        = string
  description = "Deployment environment (dev, test, prod)"
  validation {
    condition     = contains(["dev","test","prod"], var.environment)
    error_message = "environment must be one of dev, test, prod."
  }
}

variable "location" {
  type        = string
  description = "Azure region for resource deployment"
  default     = "westeurope"
}

variable "common_tags" {
  type        = map(string)
  description = "Common tags to apply to all resources"
  default     = {
    project = "chaos-engineering-framework"
  }
}
