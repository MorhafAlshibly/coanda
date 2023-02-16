# Resource group name
variable "resource_group_name" {
  type = string
}

# Cosmos DB account name
variable "cosmosdb_account_name" {
  type = string
}

# Cosmos DB Replays container name
variable "cosmosdb_replays_container_name" {
  type = string
}

# Cosmos DB Replays database name
variable "cosmosdb_replays_database_name" {
  type = string
}

# Cosmos DB archiveReplay stored procedure name
variable "archiveReplay_stored_procedure_id" {
  type = string
}

# Cosmos DB archiveReplay stored procedure code
variable "archiveReplay_stored_procedure_code" {
  type = string
}

# Create a stored procedure in a Cosmos DB container
resource "azurerm_cosmosdb_sql_stored_procedure" "this" {
  name                = var.archiveReplay_stored_procedure_id
  resource_group_name = var.resource_group_name
  account_name        = var.cosmosdb_account_name
  database_name       = var.cosmosdb_replays_database_name
  container_name      = var.cosmosdb_replays_container_name
  body                = var.archiveReplay_stored_procedure_code
}

# Define output values
output "archivereplay_resource_id" {
  value = azurerm_cosmosdb_sql_stored_procedure.this.id
}
