# Environment
variable "environment" {
  type = string
}

# Resource group name
variable "resource_group_name" {
  type = string
}

# Container app name
variable "name" {
  type = string
}

# Location
variable "location" {
  type = string
}

# Managed identity id
variable "managed_identity_id" {
  type = string
}

# Managed identity principal id
variable "managed_identity_principal_id" {
  type = string
}
