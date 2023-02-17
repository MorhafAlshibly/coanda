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

# Cosmos DB main database name
variable "cosmosdb_main_database_name" {
  type = string
}

# Cosmos DB Replays partition key
variable "cosmosdb_replays_partition_key" {
  type = string
}
