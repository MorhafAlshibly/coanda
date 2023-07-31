# Azure provider
terraform {
  backend "local" {
    path = "terraform.tfstate"
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
  name     = var.resource_group_name
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
  key_vault_name      = var.key_vault_name
}

# Include the module that creates a Cosmos DB account and database
module "cosmosdb" {
  source              = "./modules/cosmosdb"
  resource_group_name = azurerm_resource_group.this.name
  account_name        = var.cosmosdb_account_name
  database_name       = var.cosmosdb_main_database_name
  environment         = var.environment
  location            = var.location
  key_vault_id        = module.key_vault.id
  secret_name         = var.cosmosdb_secret_name
}

# Include the module that creates the Cosmos DB replay collection
module "cosmosdb_replays_collection" {
  source              = "./modules/cosmosdb/Replays"
  resource_group_name = azurerm_resource_group.this.name
  account_name        = module.cosmosdb.account_name
  collection_name     = var.cosmosdb_replays_collection_name
  database_name       = module.cosmosdb.database_name
  default_ttl_seconds = var.cosmosdb_replays_collection_default_ttl_seconds
}
