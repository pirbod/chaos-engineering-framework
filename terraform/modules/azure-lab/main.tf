resource "azurerm_resource_group" "this" {
  name     = "${var.prefix}-${var.environment}-rg"
  location = var.location
  tags     = var.tags
}

resource "azurerm_dev_test_lab" "this" {
  name                = "${var.prefix}-${var.environment}-lab"
  location            = azurerm_resource_group.this.location
  resource_group_name = azurerm_resource_group.this.name
  tags                = merge(var.tags, { environment = var.environment })
}

resource "azurerm_dev_test_virtual_network" "this" {
  name                = "${var.prefix}-${var.environment}-vnet"
  resource_group_name = azurerm_resource_group.this.name
  lab_name            = azurerm_dev_test_lab.this.name
  description         = "DevTest Lab virtual network for isolated chaos engineering demos."
  tags                = merge(var.tags, { environment = var.environment })

  subnet {
    use_in_virtual_machine_creation = "Allow"
    use_public_ip_address           = "Deny"
  }
}
