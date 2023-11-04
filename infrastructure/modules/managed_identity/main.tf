resource "azurerm_user_assigned_identity" "this" {
  name                = var.name
  resource_group_name = var.resource_group_name
  location            = var.location
  tags = {
    environment = var.environment
  }
}

resource "azurerm_role_assignment" "this" {
  scope                = azurerm_user_assigned_identity.this.id
  role_definition_name = "AcrPush"
  principal_id         = azurerm_user_assigned_identity.this.principal_id
}
