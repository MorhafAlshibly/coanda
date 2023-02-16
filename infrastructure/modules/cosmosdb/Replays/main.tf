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

# Cosmos DB Replays partition key
variable "cosmosdb_replays_partition_key" {
  type = string
}

# Create a CosmosDB container
resource "azurerm_cosmosdb_sql_container" "this" {
  name                = var.cosmosdb_replays_container_name
  resource_group_name = var.resource_group_name
  account_name        = var.cosmosdb_account_name
  database_name       = var.cosmosdb_replays_database_name
  partition_key_path  = var.cosmosdb_replays_partition_key
  throughput          = 400

  indexing_policy {
    indexing_mode = "consistent"
    automatic     = true
    included_paths {
      path    = "/*"
      indexes = ["Range", "Spatial"]
    }
  }
}

output "cosmosdb_replays_container_endpoint" {
  value = azurerm_cosmosdb_sql_container.this.endpoint
}

output "cosmosdb_replays_container_read_write_keys" {
  value = azurerm_cosmosdb_sql_container.this.primary_master_key
}

output "cosmosdb_replays_container_id" {
  value = azurerm_cosmosdb_sql_container.this.id
}
