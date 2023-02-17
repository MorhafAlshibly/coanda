# Create a CosmosDB collection
resource "azurerm_cosmosdb_mongo_collection" "this" {
  name                = var.cosmosdb_replays_collection_name
  resource_group_name = var.resource_group_name
  account_name        = var.cosmosdb_account_name
  database_name       = var.cosmosdb_main_database_name

  default_ttl_seconds = "86400"
  shard_key           = "_id"
  throughput          = 400

  index {
    keys   = ["_id"]
    unique = true
  }
}
