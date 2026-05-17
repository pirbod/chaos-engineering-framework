variable "subscription_id" {
  description = "Azure subscription ID."
  type        = string
}

variable "prefix" {
  description = "Resource name prefix."
  type        = string
  default     = "pch"
}

variable "location" {
  description = "Azure region."
  type        = string
  default     = "westeurope"
}

variable "common_tags" {
  description = "Common tags."
  type        = map(string)
  default = {
    project     = "platform-chaos-hub"
    environment = "demo"
  }
}
