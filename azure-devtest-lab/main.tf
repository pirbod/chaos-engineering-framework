terraform {
  required_version = ">= 1.4.0, < 2.0.0"
  required_providers {
    azurerm = {
      source  = "hashicorp/azurerm"
      version = "~> 3.75"
    }
  }
}

provider "azurerm" {
  features        = {}
  subscription_id = var.subscription_id
}

resource "azurerm_resource_group" "rg" {
  name     = "${var.prefix}-${var.environment}-rg"
  location = var.location
  tags     = var.common_tags
}

resource "azurerm_dev_test_lab" "lab" {
  name                = "${var.prefix}-${var.environment}-lab"
  location            = azurerm_resource_group.rg.location
  resource_group_name = azurerm_resource_group.rg.name
  lab_storage_type    = "Standard"
  tags                = merge(var.common_tags, { environment = var.environment })
}

resource "azurerm_dev_test_virtual_network" "vnet" {
  name                = "${var.prefix}-${var.environment}-vnet"
  resource_group_name = azurerm_resource_group.rg.name
  lab_name            = azurerm_dev_test_lab.lab.name
  location            = azurerm_resource_group.rg.location
  address_space       = ["10.1.0.0/16"]
}
