# Environment
variable "environment" {
  type = string
}

# Location
variable "location" {
  type = string
}

# Resource group name
variable "resource_group_name" {
  type = string
}

# Cosmos DB account name
variable "account_name" {
  type = string
}

# Cosmos DB main database name
variable "database_name" {
  type = string
}

# Key vault id to store secret
variable "key_vault_id" {
  type = string
}

# Name of secret to store connection string
variable "secret_name" {
  type = string
}
