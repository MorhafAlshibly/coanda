# Azure provider
terraform {
  backend "azurerm" {
    resource_group_name  = "tfstates"
    storage_account_name = "tfstatesaccount"
    container_name       = "tfstatedevops"
    key                  = "tfstatedevops.tfstate"
  }
  required_providers {
    azurerm = {
      source  = "hashicorp/azurerm"
      version = "=3.0.0"
    }
  }
}

provider "azurerm" {
  features {}
}

# Development resource group
resource "azurerm_resource_group" "this" {
  name     = var.resource_group_name
  location = var.location
  tags = {
    environment = "dev"
  }
}

# Include the module that creates a Cosmos DB account and database
module "cosmosdb" {
  source                = "../modules/cosmosdb"
  resource_group_name   = azurerm_resource_group.this.name
  cosmosdb_account_name = var.cosmosdb_account_name
  environment           = var.environment
  location              = var.location
}

# Include the module that creates the Cosmos DB replay collection
module "cosmosdb_replays_collection" {
  source                           = "../modules/cosmosdb/Replays"
  resource_group_name              = azurerm_resource_group.this.name
  cosmosdb_account_name            = module.cosmosdb.account_name
  cosmosdb_replays_collection_name = var.cosmosdb_replays_collection_name
  cosmosdb_main_database_name      = module.cosmosdb.database_name
}
