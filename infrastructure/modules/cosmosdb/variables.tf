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

# Key vault id to store secret
variable "key_vault_id" {
  type = string
}

# Name of secret to store connection string
variable "secret_name" {
  type = string
}

# VNet subnet id
variable "vnet_subnet_id" {
  type = string
}

# App configuration id
variable "app_configuration_id" {
  type = string
}
