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

# Key vault name
variable "key_vault_name" {
  type = string
}

# VNet subnet id
variable "vnet_subnet_id" {
  type = string
}

# MongoDB connection string secret name
variable "mongo_connection_secret_name" {
  type = string
}

# MongoDB connection string
variable "mongo_connection_string" {
  type = string
}

# Managed identity tenant id
variable "managed_identity_tenant_id" {
  type = string
}

# Managed identity principal id
variable "managed_identity_principal_id" {
  type = string
}
