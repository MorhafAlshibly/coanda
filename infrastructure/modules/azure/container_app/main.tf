# Create an azure container apps service
resource "azurerm_container_app_environment" "this" {
  name                       = var.container_app_environment_name
  location                   = var.location
  resource_group_name        = var.resource_group_name
  log_analytics_workspace_id = var.log_analytics_workspace_id

  internal_load_balancer_enabled = false
  infrastructure_subnet_id       = var.vnet_subnet_id

  tags = {
    environment = var.environment
  }
}

# Create an azure container apps service
resource "azurerm_container_app" "this" {
  name                         = var.name
  resource_group_name          = var.resource_group_name
  container_app_environment_id = azurerm_container_app_environment.this.id
  revision_mode                = "Single"

  identity {
    type = "UserAssigned"
    identity_ids = [
      var.managed_identity_id
    ]
  }

  registry {
    server   = var.registry_uri
    identity = var.managed_identity_id
  }

  secret {
    name  = var.mongo_connection_secret_name
    value = var.mongo_connection_string
  }

  ingress {
    allow_insecure_connections = false
    target_port                = 8080
    external_enabled           = true
    traffic_weight {
      latest_revision = true
      percentage      = 100
    }
  }

  template {
    container {
      name   = "bff"
      image  = "${var.registry_uri}/bff:latest"
      cpu    = 0.25
      memory = "0.5Gi"
      env {
        name  = "BFF_ENABLEPLAYGROUND"
        value = var.environment == "dev" ? "true" : "false"
      }
    }
    container {
      name   = "item"
      image  = "${var.registry_uri}/item:latest"
      cpu    = 0.25
      memory = "0.5Gi"
      env {
        name  = "ITEM_TABLECONN"
        value = var.storage_table_connection_string
      }
      env {
        name  = "AZURE_CLIENT_ID"
        value = var.managed_identity_client_id
      }
    }
    container {
      name   = "team"
      image  = "${var.registry_uri}/team:latest"
      cpu    = 0.25
      memory = "0.5Gi"
      env {
        name        = "TEAM_MONGOCONN"
        secret_name = var.mongo_connection_secret_name
      }
      env {
        name  = "AZURE_CLIENT_ID"
        value = var.managed_identity_client_id
      }
    }
    container {
      name   = "record"
      image  = "${var.registry_uri}/record:latest"
      cpu    = 0.25
      memory = "0.5Gi"
      env {
        name        = "RECORD_MONGOCONN"
        secret_name = var.mongo_connection_secret_name
      }
      env {
        name  = "AZURE_CLIENT_ID"
        value = var.managed_identity_client_id
      }
    }
    container {
      name   = "redis"
      image  = "docker.io/redis:alpine"
      cpu    = 0.25
      memory = "0.5Gi"
    }
  }

  tags = {
    environment = var.environment
  }
}
