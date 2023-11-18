# Environment
variable "environment" {
  type = string
}

# Container app environment name
variable "container_app_environment_name" {
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

# Managed identity id
variable "managed_identity_id" {
  type = string
}

# VNet subnet id
variable "vnet_subnet_id" {
  type = string
}

# Container registry url
variable "registry_uri" {
  type = string
}

# Location
variable "location" {
  type = string
}

# Log Analytics workspace id
variable "log_analytics_workspace_id" {
  type = string
}

# App Configuration endpoint
variable "app_configuration_endpoint" {
  type = string
}

# Managed identity client id
variable "managed_identity_client_id" {
  type = string
}

# Managed identity tenant id
variable "managed_identity_tenant_id" {
  type = string
}
