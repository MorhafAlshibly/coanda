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

# Include the module that creates a Key Vault
module "key_vault" {
  source              = "./modules/key_vault"
  resource_group_name = azurerm_resource_group.this.name
  environment         = var.environment
  location            = var.location
  key_vault_name      = format("kv-%s-%s-%s", var.app_name, var.environment, var.location)
}

# Include the module that creates a Cosmos DB account and database
module "cosmosdb" {
  source              = "./modules/cosmosdb"
  resource_group_name = azurerm_resource_group.this.name
  account_name        = format("cdb-%s-%s-%s", var.app_name, var.environment, var.location)
  database_name       = format("db-%s-%s-%s", var.app_name, var.environment, var.location)
  environment         = var.environment
  location            = var.location
  key_vault_id        = module.key_vault.id
  secret_name         = format("cdb-%s-secret-%s-%s", var.app_name, var.environment, var.location)
}

# Include the module that creates a container registry
module "container_registry" {
  source              = "./modules/container_registry"
  resource_group_name = azurerm_resource_group.this.name
  name                = format("acr%s%s%s", var.app_name, var.environment, var.location)
  environment         = var.environment
  location            = var.location
}

# Include the module that creates a virtual network
module "virtual_network" {
  source              = "./modules/virtual_network"
  resource_group_name = azurerm_resource_group.this.name
  name                = format("vnet-%s-%s-%s", var.app_name, var.environment, var.location)
  environment         = var.environment
  location            = var.location
}

# Include the module that creates a managed identity
module "managed_identity" {
  source              = "./modules/managed_identity"
  resource_group_name = azurerm_resource_group.this.name
  name                = format("mi-%s-%s-%s", var.app_name, var.environment, var.location)
  environment         = var.environment
  location            = var.location
}
