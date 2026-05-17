variable "subscription_id" {
  description = "Azure subscription ID."
  type        = string
}

variable "prefix" {
  description = "Resource name prefix."
  type        = string
  default     = "pch"
}

variable "environment" {
  description = "Environment name."
  type        = string
  default     = "dev"
}

variable "location" {
  description = "Azure region."
  type        = string
  default     = "westeurope"
}

variable "alert_email" {
  description = "Optional platform on-call email."
  type        = string
  default     = ""
}

variable "common_tags" {
  description = "Common tags."
  type        = map(string)
  default = {
    project = "platform-chaos-hub"
  }
}
