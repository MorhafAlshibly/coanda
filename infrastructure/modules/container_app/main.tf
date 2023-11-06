# Create an azure container apps service
data "azurerm_container_app_environment" "this" {
  name                = var.container_app_environment_name
  resource_group_name = var.resource_group_name
}

# Create an azure container apps service
data "azurerm_container_app" "this" {
  name                         = var.name
  resource_group_name          = var.resource_group_name
  container_app_environment_id = data.azurerm_container_app_environment.this.id

  identity {
    type = "UserAssigned"
    identity_ids = [
      var.managed_identity_id
    ]
  }
}
