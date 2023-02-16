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
  offer_type          = "Standard"

  consistency_policy {
    consistency_level       = "BoundedStaleness"
    max_interval_in_seconds = 300
    max_staleness_prefix    = 100000
  }

  geo_location {
    location = var.location
  }

  tags = {
    environment = var.environment
  }
}

# Endpoint of the DB
output "cosmosdb_endpoint" {
  value = azurerm_cosmosdb_account.this.endpoint
}

# Cosmos DB key
output "cosmosdb_key" {
  value = azurerm_cosmosdb_account.this.primary_master_key
}
