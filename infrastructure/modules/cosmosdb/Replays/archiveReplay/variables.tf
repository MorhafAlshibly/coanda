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
