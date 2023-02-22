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
  type = string
}

# Cosmos DB main database name
variable "cosmosdb_main_database_name" {
  type = string
}

# Cosmos DB replay collection name
variable "cosmosdb_replays_collection_name" {
  type = string
}

# Key vault name
variable "key_vault_name" {
  type = string
}

# CosmosDB connection string secret name
variable "cosmosdb_secret_name" {
  type = string
}

# Replay collection default TTL seconds
variable "cosmosdb_replays_collection_default_ttl_seconds" {
  type = number
}
