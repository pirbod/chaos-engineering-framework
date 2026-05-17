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

module "azure_lab" {
  source      = "../../modules/azure-lab"
  prefix      = var.prefix
  environment = "demo"
  location    = var.location
  tags        = var.common_tags
}
