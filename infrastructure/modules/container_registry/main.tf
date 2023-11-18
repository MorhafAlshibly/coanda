# Create a container registry
resource "azurerm_container_registry" "this" {
  name                = var.name
  resource_group_name = var.resource_group_name
  location            = var.location
  sku                 = "Standard"
  admin_enabled       = false

  tags = {
    environment = var.environment
  }
}

# Docker compose run
resource "null_resource" "docker_compose" {
  depends_on = [azurerm_container_registry.this]
  triggers = {
    registry_uri = azurerm_container_registry.this.login_server
  }
  provisioner "local-exec" {
    command = format("task push ENV=%s", var.environment)
  }
}

data "azurerm_container_registry" "after_docker_compose" {
  depends_on          = [null_resource.docker_compose]
  name                = var.name
  resource_group_name = var.resource_group_name
}
