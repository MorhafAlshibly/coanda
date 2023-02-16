# Environment
variable "environment" {
  type    = string
  default = "dev"
}

# Location
variable "location" {
  type    = string
  default = "eastus"
}

# Resource group name
variable "resource_group_name" {
  type    = string
  default = "coanda-resources"
}

# Cosmos DB account name
variable "cosmosdb_account_name" {
  type    = string
  default = "coanda-cosmosdb"
}

# Cosmos DB replay container name
variable "cosmosdb_replays_container_name" {
  type    = string
  default = "Replays"
}

# Cosmos DB Replays database name
variable "cosmosdb_replays_database_name" {
  type    = string
  default = "coanda-cosmosdb-main"
}

# Cosmos DB Replays partition key
variable "cosmosdb_replays_partition_key" {
  type    = string
  default = "/id"
}

# Cosmos DB archiveReplay stored procedure name
variable "archiveReplay_stored_procedure_id" {
  type    = string
  default = "archiveReplay"
}

# Cosmos DB archiveReplay stored procedure code
variable "archiveReplay_stored_procedure_code" {
  type    = string
  default = "temp"
}
