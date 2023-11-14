# Resource group name
output "resource_group_name" {
  value = azurerm_resource_group.this.name
}

# Cosmos DB account name
output "cosmosdb_account_name" {
  value = module.cosmosdb.account_name
}

# Cosmos DB main database name
output "cosmosdb_main_database_name" {
  value = module.cosmosdb.database_name
}

# Vault URI
output "vault_uri" {
  value = module.key_vault.vault_uri
}

# Storage account uri
output "storage_uri" {
  value = module.storage_account.storage_uri
}

# Container registry uri
output "container_registry_uri" {
  value = module.container_registry.uri
}
