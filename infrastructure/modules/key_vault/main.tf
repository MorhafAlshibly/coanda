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

  tags = {
    environment = var.environment
  }
}

# Grant access to the key vault by managed identity
resource "azurerm_role_assignment" "this" {
  scope                = azurerm_key_vault.this.id
  role_definition_name = "Key Vault Secrets User"
  principal_id         = var.managed_identity_principal_id
}
