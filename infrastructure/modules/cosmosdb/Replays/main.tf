# Create a CosmosDB collection
resource "azurerm_cosmosdb_mongo_collection" "Replays" {
  name                = var.collection_name
  resource_group_name = var.resource_group_name
  account_name        = var.account_name
  database_name       = var.database_name

  default_ttl_seconds = var.default_ttl_seconds
  shard_key           = "_id"
  throughput          = 400

  index {
    keys   = ["_id"]
    unique = true
  }
}
