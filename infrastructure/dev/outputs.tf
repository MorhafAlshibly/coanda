# Resource group name
output "resource_group_name" {
  value = azurerm_resource_group.this.name
}

# Cosmos DB account name
output "cosmosdb_account_name" {
  value = module.cosmosdb.cosmosdb_account_name
}

# Cosmos DB main database name
output "cosmosdb_main_database_name" {
  value = module.cosmosdb.cosmosdb_main_database_name
}

# Cosmos DB replays container name
output "cosmosdb_replays_container_name" {
  value = module.cosmosdb_replays_container.cosmosdb_replays_container_name
}

# archiveReplay resource id
output "archiveReplay_resource_id" {
  value = module.archiveReplay_stored_procedure.archiveReplay_resource_id
}
