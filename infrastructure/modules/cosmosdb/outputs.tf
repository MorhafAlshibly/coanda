# Cosmos DB account name
output "cosmosdb_account_name" {
  value = azurerm_cosmosdb_account.this.name
}

# Cosmos DB main database name
output "cosmosdb_main_database_name" {
  value = azurerm_cosmosdb_sql_database.this.name
}
