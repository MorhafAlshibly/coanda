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

# Cosmos DB account name
variable "cosmosdb_account_name" {
  type = string
}

# Create a Cosmos DB account
resource "azurerm_cosmosdb_account" "this" {
  name                = var.cosmosdb_account_name
  location            = var.location
  resource_group_name = var.resource_group_name
  enable_free_tier    = true
  offer_type          = "Standard"
  kind                = "GlobalDocumentDB"
  capacity {
    total_throughput_limit = 1000
  }
  geo_location {
    location          = var.location
    failover_priority = 0
  }
  consistency_policy {
    consistency_level       = "BoundedStaleness"
    max_interval_in_seconds = 300
    max_staleness_prefix    = 100000
  }
  tags = {
    environment = var.environment
  }
}
