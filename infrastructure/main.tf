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

# Development environment
module "dev_environment" {
  source = "./environments/dev"
}
