# Create a log analytics workspace
resource "azurerm_log_analytics_workspace" "this" {
  name                = var.name
  resource_group_name = var.resource_group_name
  location            = var.location
  sku                 = "Free"
  retention_in_days   = 7

  tags = {
    environment = var.environment
  }
}
