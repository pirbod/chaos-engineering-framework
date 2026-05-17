terraform {
  required_version = ">= 1.4.0, < 2.0.0"

  required_providers {
    azurerm = {
      source  = "hashicorp/azurerm"
      version = "~> 3.117"
    }
  }
}

provider "azurerm" {
  features {}
  subscription_id = var.subscription_id
}

locals {
  tags = merge(var.common_tags, {
    environment = var.environment
    managed_by  = "terraform"
  })
}

module "azure_lab" {
  source      = "../../modules/azure-lab"
  prefix      = var.prefix
  environment = var.environment
  location    = var.location
  tags        = local.tags
}

module "managed_identities" {
  source              = "../../modules/managed-identities"
  prefix              = var.prefix
  environment         = var.environment
  location            = var.location
  resource_group_name = module.azure_lab.resource_group_name
  tags                = local.tags
}

module "aks_observability" {
  source              = "../../modules/aks-observability"
  prefix              = var.prefix
  environment         = var.environment
  location            = var.location
  resource_group_name = module.azure_lab.resource_group_name
  alert_email         = var.alert_email
  tags                = local.tags
}
