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

# Add storage table connection string to app configuration
resource "azurerm_app_configuration_key" "storage_table_connection_string" {
  configuration_store_id = var.app_configuration_id
  key                    = var.configuration_key
  value                  = format("https://%s.table.core.windows.net", azurerm_storage_account.this.name)
}
