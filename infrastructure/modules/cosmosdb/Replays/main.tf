# Create a CosmosDB container
resource "azurerm_cosmosdb_sql_container" "this" {
  name                  = var.cosmosdb_replays_container_name
  resource_group_name   = var.resource_group_name
  account_name          = var.cosmosdb_account_name
  database_name         = var.cosmosdb_main_database_name
  partition_key_path    = var.cosmosdb_replays_partition_key
  partition_key_version = 1
  throughput            = 400
  default_ttl           = 86400

  indexing_policy {
    indexing_mode = "consistent"
    included_path {
      path = "/*"
    }
    included_path {
      path = "/included/?"
    }
    excluded_path {
      path = "/excluded/?"
    }
  }
  unique_key {
    paths = ["/definition/idlong", "/definition/idshort"]
  }
}
