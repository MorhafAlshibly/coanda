# Create a Cosmos DB account
resource "azurerm_cosmosdb_account" "this" {
  name                 = var.account_name
  location             = var.location
  resource_group_name  = var.resource_group_name
  enable_free_tier     = true
  offer_type           = "Standard"
  kind                 = "MongoDB"
  mongo_server_version = "4.2"

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

  is_virtual_network_filter_enabled = true
  virtual_network_rule {
    id                                   = var.vnet_subnet_id
    ignore_missing_vnet_service_endpoint = true
  }

  tags = {
    environment = var.environment
  }
}
