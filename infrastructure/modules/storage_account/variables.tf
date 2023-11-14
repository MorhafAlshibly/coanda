# Storage account name
variable "name" {
  type = string
}

# Resource group name
variable "resource_group_name" {
  type = string
}

# Environment
variable "environment" {
  type = string
}

# Location
variable "location" {
  type = string
}

# Managed identity principal ID
variable "managed_identity_principal_id" {
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
