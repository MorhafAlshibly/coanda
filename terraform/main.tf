terraform {
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


resource "azurerm_resource_group" "coanda-rg" {
  name     = "coanda-resources"
  location = "East US"

  tags = {
    environment = "dev"
  }
}

resource "azurerm_virtual_network" "coanda-vn" {
  name                = "coanda-network"
  resource_group_name = azurerm_resource_group.coanda-rg.name
  location            = azurerm_resource_group.coanda-rg.location
  address_space       = ["10.0.0.0/16"]

  tags = {
    environment = "dev"
  }
}
