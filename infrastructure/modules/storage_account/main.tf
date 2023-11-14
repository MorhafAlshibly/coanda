# Create Azure Storage Account
resource "azurerm_storage_account" "this" {
  name                     = var.name
  resource_group_name      = var.resource_group_name
  location                 = var.location
  account_tier             = "Standard"
  account_replication_type = "LRS"

  network_rules {
    default_action = "Deny"
    bypass         = ["None"]
    virtual_network_subnet_ids = [
      var.vnet_subnet_id,
    ]
  }

  tags = {
    environment = var.environment
  }
}

# Grant access to the storage account by managed identity
resource "azurerm_role_assignment" "this" {
  scope                = azurerm_storage_account.this.id
  role_definition_name = "Storage Table Data Contributor"
  principal_id         = var.managed_identity_principal_id
}


# Add storage account connection string to app configuration
resource "azurerm_app_configuration_key" "storage_account_connection_string" {
  configuration_store_id = var.app_configuration_id
  key                    = var.
  value                  = azurerm_storage_account.this.primary_connection_string
}