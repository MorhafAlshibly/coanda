# Cosmos DB account name
output "account_name" {
  value = azurerm_cosmosdb_account.this.name
}

# Cosmos DB main database name
output "database_name" {
  value = azurerm_cosmosdb_mongo_database.main.name
}

# Cosmos DB connection strings
output "connection_strings" {
  value     = azurerm_cosmosdb_account.this.connection_strings
  sensitive = true
}
