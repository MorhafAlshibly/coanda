# Resource group name
variable "resource_group_name" {
  type = string
}

# Cosmos DB account name
variable "account_name" {
  type = string
}

# Cosmos DB Replays collection name
variable "collection_name" {
  type = string
}

# Cosmos DB main database name
variable "database_name" {
  type = string
}

# Collection default TTL seconds
variable "default_ttl_seconds" {
  type = number
}
