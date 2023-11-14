# Create a user-assigned identity
resource "azurerm_user_assigned_identity" "this" {
  name                = var.name
  resource_group_name = var.resource_group_name
  location            = var.location
  tags = {
    environment = var.environment
  }
}

resource "azurerm_role_assignment" "acr_pull" {
  scope                = var.container_registry_id
  role_definition_name = "AcrPull"
  principal_id         = azurerm_user_assigned_identity.this.principal_id
}

data "azurerm_user_assigned_identity" "after_acr_pull" {
  depends_on          = [azurerm_role_assignment.acr_pull]
  name                = var.name
  resource_group_name = var.resource_group_name
}
