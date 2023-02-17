# Cosmos DB replays container name
output "cosmosdb_replays_container_name" {
  value = azurerm_cosmosdb_sql_container.this.name
}
