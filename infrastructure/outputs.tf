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

# Cosmos DB replays collection name
output "cosmosdb_replays_collection_name" {
  value = module.cosmosdb_replays_collection.collection_name
}

# Cosmos DB connection strings
output "cosmosdb_connection_strings" {
  value     = module.cosmosdb.connection_strings
  sensitive = true
}
