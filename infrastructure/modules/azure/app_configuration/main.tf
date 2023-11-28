# Create a azure app configuration instance
resource "azurerm_app_configuration" "this" {
  name                = var.name
  resource_group_name = var.resource_group_name
  location            = var.location
  sku                 = "free"

  identity {
    type = "UserAssigned"
    identity_ids = [
      var.managed_identity_id,
    ]
  }

  tags = {
    environment = var.environment
  }
}

# Give the managed identity access to the app configuration instance
resource "azurerm_role_assignment" "app_configuration" {
  scope                = azurerm_app_configuration.this.id
  role_definition_name = "App Configuration Data Reader"
  principal_id         = var.managed_identity_principal_id
}
