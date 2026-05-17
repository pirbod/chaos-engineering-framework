resource "azurerm_log_analytics_workspace" "this" {
  name                = "${var.prefix}-${var.environment}-law"
  location            = var.location
  resource_group_name = var.resource_group_name
  sku                 = "PerGB2018"
  retention_in_days   = var.retention_in_days
  tags                = var.tags
}

resource "azurerm_monitor_action_group" "platform" {
  name                = "${var.prefix}-${var.environment}-platform-ag"
  resource_group_name = var.resource_group_name
  short_name          = "pch"
  tags                = var.tags

  dynamic "email_receiver" {
    for_each = var.alert_email == "" ? [] : [var.alert_email]
    content {
      name          = "platform-oncall"
      email_address = email_receiver.value
    }
  }
}
