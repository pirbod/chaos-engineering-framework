resource "azurerm_user_assigned_identity" "platform" {
  name                = "${var.prefix}-${var.environment}-platform-mi"
  location            = var.location
  resource_group_name = var.resource_group_name
  tags                = var.tags
}
