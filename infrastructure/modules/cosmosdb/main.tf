# Create a Cosmos DB account
resource "azurerm_cosmosdb_account" "this" {
  name                = var.account_name
  location            = var.location
  resource_group_name = var.resource_group_name
  enable_free_tier    = true
  offer_type          = "Standard"
  kind                = "MongoDB"

  capabilities {
    name = "mongoEnableDocLevelTTL"
  }
  capabilities {
    name = "MongoDBv3.4"
  }
  capabilities {
    name = "EnableMongo"
  }

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

# Create a CosmosDB database
resource "azurerm_cosmosdb_mongo_database" "main" {
  name                = var.database_name
  resource_group_name = var.resource_group_name
  account_name        = azurerm_cosmosdb_account.this.name
}

# Add Cosmos DB connection string to key vault
resource "azurerm_key_vault_secret" "cosmosdb_connection_string" {
  name         = var.secret_name
  value        = tolist(azurerm_cosmosdb_account.this.connection_strings)[0]
  key_vault_id = var.key_vault_id
}
