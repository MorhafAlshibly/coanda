# Create a azure app configuration instance
resource "azurerm_app_configuration" "this" {
  name                = var.name
  resource_group_name = var.resource_group_name
  location            = var.location
  sku                 = "Free"

  tags = {
    environment = var.environment
  }
}
