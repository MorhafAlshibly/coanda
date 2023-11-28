data "azuread_client_config" "current" {}

# Create key vault to store secrets
resource "azurerm_key_vault" "this" {
  name                          = var.key_vault_name
  location                      = var.location
  resource_group_name           = var.resource_group_name
  enabled_for_disk_encryption   = true
  tenant_id                     = data.azuread_client_config.current.tenant_id
  soft_delete_retention_days    = 7
  purge_protection_enabled      = false
  public_network_access_enabled = true

  sku_name = "standard"

  access_policy {
    tenant_id = data.azuread_client_config.current.tenant_id
    object_id = data.azuread_client_config.current.object_id

    key_permissions = [
      "Get", "List", "Update", "Create", "Import", "Delete", "Recover", "Backup", "Restore", "Purge",
    ]

    secret_permissions = [
      "Get", "List", "Delete", "Recover", "Backup", "Restore", "Set", "Purge",
    ]
  }

  access_policy {
    tenant_id = var.managed_identity_tenant_id
    object_id = var.managed_identity_principal_id

    secret_permissions = [
      "Get",
    ]
  }

  tags = {
    environment = var.environment
  }
}

# Add Cosmos DB connection string to key vault
resource "azurerm_key_vault_secret" "mongo_connection_string" {
  name         = var.mongo_connection_secret_name
  value        = var.mongo_connection_string
  key_vault_id = azurerm_key_vault.this.id
}

# Add network acl rules to key vault
resource "null_resource" "key_vault_network_rules" {
  triggers = {
    key_vault_secret = azurerm_key_vault_secret.mongo_connection_string.id
  }

  provisioner "local-exec" {
    command = "az keyvault network-rule add --name ${azurerm_key_vault.this.name} --subnet ${var.vnet_subnet_id} --resource-group ${var.resource_group_name}"
  }
}
