# Create an azure container apps service
resource "azurerm_container_app_environment" "this" {
  name                       = var.container_app_environment_name
  location                   = var.location
  resource_group_name        = var.resource_group_name
  log_analytics_workspace_id = var.log_analytics_workspace_id

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

  ingress {
    allow_insecure_connections = false
    target_port                = 8080
    external_enabled           = true
    traffic_weight {
      percentage = 100
    }
  }

  template {
    container {
      name   = "bff"
      image  = "${var.registry_uri}/bff:latest"
      cpu    = 0.25
      memory = "0.5Gi"
    }
    container {
      name   = "item"
      image  = "${var.registry_uri}/item:latest"
      cpu    = 0.25
      memory = "0.5Gi"
    }
    container {
      name   = "team"
      image  = "${var.registry_uri}/team:latest"
      cpu    = 0.25
      memory = "0.5Gi"
    }
    container {
      name   = "record"
      image  = "${var.registry_uri}/record:latest"
      cpu    = 0.25
      memory = "0.5Gi"
    }
    container {
      name   = "redis"
      image  = "redis:alpine"
      cpu    = 0.25
      memory = "0.5Gi"
    }
  }

  tags = {
    environment = var.environment
  }
}
