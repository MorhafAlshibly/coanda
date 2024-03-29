# Azure provider
terraform {
  backend "azurerm" {
    resource_group_name  = var.resource_group_name
    storage_account_name = var.storage_account_name
    container_name       = var.container_name
    key                  = var.key
  }
  required_providers {
    azurerm = {
      source  = "hashicorp/azurerm"
      version = ">=3.0.0"
    }
  }
}

provider "azurerm" {
  features {}
}

# Resource group
resource "azurerm_resource_group" "this" {
  name     = format("rg-%s-%s-%s", var.app_name, var.environment, var.location)
  location = var.location
  tags = {
    environment = var.environment
  }
}

# Include the module that creates a managed identity
module "managed_identity" {
  source                = "./managed_identity"
  resource_group_name   = azurerm_resource_group.this.name
  name                  = format("mi%s%s%s", var.app_name, var.environment, var.location)
  environment           = var.environment
  location              = var.location
  container_registry_id = module.container_registry.id
  storage_account_id    = module.storage_account.id
}

# Include the module that creates a Key Vault
#module "key_vault" {
#  source                        = "./key_vault"
#  resource_group_name           = azurerm_resource_group.this.name
#  environment                   = var.environment
#  location                      = var.location
#  key_vault_name                = format("kv-%s-%s-%s", var.app_name, var.environment, var.location)
#  vnet_subnet_id                = module.virtual_network.vnet_subnet_id
#  mongo_connection_string       = module.cosmosdb.connection_string
#  managed_identity_tenant_id    = module.managed_identity.tenant_id
#  managed_identity_principal_id = module.managed_identity.principal_id
#}

# Include the module that creates a Cosmos DB account and database
module "cosmosdb" {
  source              = "./cosmosdb"
  resource_group_name = azurerm_resource_group.this.name
  account_name        = format("cdb-%s-%s-%s", var.app_name, var.environment, var.location)
  environment         = var.environment
  location            = var.location
  vnet_subnet_id      = module.virtual_network.vnet_subnet_id
}

# Include the module that creates a container registry
module "container_registry" {
  source              = "./container_registry"
  resource_group_name = azurerm_resource_group.this.name
  name                = format("acr%s%s%s", var.app_name, var.environment, var.location)
  environment         = var.environment
  location            = var.location
  managed_identity_id = module.managed_identity.id
}

# Include the module that creates a virtual network
module "virtual_network" {
  source              = "./virtual_network"
  resource_group_name = azurerm_resource_group.this.name
  name                = format("vnet-%s-%s-%s", var.app_name, var.environment, var.location)
  environment         = var.environment
  location            = var.location
}

# Include the module that creates a storage account
module "storage_account" {
  source              = "./storage_account"
  resource_group_name = azurerm_resource_group.this.name
  name                = format("sa%s%s%s", var.app_name, var.environment, var.location)
  environment         = var.environment
  location            = var.location
  vnet_subnet_id      = module.virtual_network.vnet_subnet_id
}

# Include the module that creates a log analytics workspace
module "log_analytics_workspace" {
  source              = "./log_analytics_workspace"
  resource_group_name = azurerm_resource_group.this.name
  name                = format("law%s%s%s", var.app_name, var.environment, var.location)
  environment         = var.environment
  location            = var.location
}

# Include the module that creates a container app
module "container_app" {
  source                          = "./container_app"
  environment                     = var.environment
  container_app_environment_name  = format("cae-%s-%s-%s", var.app_name, var.environment, var.location)
  resource_group_name             = azurerm_resource_group.this.name
  name                            = format("ca-%s-%s-%s", var.app_name, var.environment, var.location)
  location                        = var.location
  managed_identity_id             = module.managed_identity.id
  vnet_subnet_id                  = module.virtual_network.vnet_subnet_id
  registry_uri                    = module.container_registry.uri
  log_analytics_workspace_id      = module.log_analytics_workspace.id
  storage_table_connection_string = format("https://%s.table.core.windows.net", module.storage_account.name)
  mongo_connection_secret_name    = format("cdb-secret-%s-%s-%s", var.app_name, var.environment, var.location)
  mongo_connection_string         = module.cosmosdb.connection_string
  managed_identity_client_id      = module.managed_identity.client_id
}

# Include the module that creates a NAT gateway
#module "nat_gateway" {
#  source              = "./nat_gateway"
#  environment         = var.environment
#  resource_group_name = azurerm_resource_group.this.name
#  name                = format("nat-%s-%s-%s", var.app_name, var.environment, var.location)
#  vnet_subnet_id      = module.virtual_network.vnet_subnet_id
#  location            = var.location
#  ip_address_id       = module.ip_address.id
#}

# Include the module that creates an ip address
#module "ip_address" {
#  source              = "./ip_address"
#  environment         = var.environment
#  resource_group_name = azurerm_resource_group.this.name
#  name                = format("ip-%s-%s-%s", var.app_name, var.environment, var.location)
#  location            = var.location
#}

# Include the module that creates an app configuration
#module "app_configuration" {
#  source                        = "./app_configuration"
#  environment                   = var.environment
#  resource_group_name           = azurerm_resource_group.this.name
#  name                          = format("ac-%s-%s-%s", var.app_name, var.environment, var.location)
#  location                      = var.location
#  managed_identity_id           = module.managed_identity.id
#  managed_identity_principal_id = module.managed_identity.principal_id
#}
