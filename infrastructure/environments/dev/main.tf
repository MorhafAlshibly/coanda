resource "azurerm_resource_group" "coanda-rg" {
  name     = var.resource_group_name
  location = var.location

  tags = {
    environment = "dev"
  }
}

# Include the module that creates a Cosmos DB account
module "cosmosdb" {
  source                = "../../modules/cosmosdb"
  resource_group_name   = var.resource_group_name
  cosmosdb_account_name = var.cosmosdb_account_name
  environment           = var.environment
  location              = var.location
}

# Include the module that creates the Cosmos DB replay container
module "cosmosdb_replays_container" {
  source                          = "../../modules/cosmosdb/Replays"
  resource_group_name             = var.resource_group_name
  cosmosdb_account_name           = var.cosmosdb_account_name
  cosmosdb_replays_container_name = var.cosmosdb_replays_container_name
  cosmosdb_replays_database_name  = var.cosmosdb_replays_database_name
  cosmosdb_replays_partition_key  = var.cosmosdb_replays_partition_key
}

# Include the module that creates the Cosmos DB replay container
module "archiveReplay_stored_procedure" {
  source                              = "../../modules/cosmosdb/Replays/archiveReplay"
  resource_group_name                 = var.resource_group_name
  cosmosdb_account_name               = var.cosmosdb_account_name
  cosmosdb_replays_container_name     = var.cosmosdb_replays_container_name
  cosmosdb_replays_database_name      = var.cosmosdb_replays_database_name
  archiveReplay_stored_procedure_id   = var.archiveReplay_stored_procedure_id
  archiveReplay_stored_procedure_code = var.archiveReplay_stored_procedure_code
}
