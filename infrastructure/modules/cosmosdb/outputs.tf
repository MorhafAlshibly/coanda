# Cosmos DB account name
output "account_name" {
  value = azurerm_cosmosdb_account.this.name
}

# MongoDB connection string
output "connection_string" {
  value = tolist(azurerm_cosmosdb_account.this.connection_strings)[0]
}
