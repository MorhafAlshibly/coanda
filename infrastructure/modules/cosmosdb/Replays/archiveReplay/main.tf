# Create a stored procedure in a Cosmos DB container
resource "azurerm_cosmosdb_sql_stored_procedure" "this" {
  name                = var.archiveReplay_stored_procedure_id
  resource_group_name = var.resource_group_name
  account_name        = var.cosmosdb_account_name
  database_name       = var.cosmosdb_replays_database_name
  container_name      = var.cosmosdb_replays_container_name
  body                = var.archiveReplay_stored_procedure_code
}
