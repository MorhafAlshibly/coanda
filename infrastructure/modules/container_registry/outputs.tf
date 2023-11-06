# Container registry uri
output "container_registry_uri" {
  value = format("%s.azurecr.io", azurerm_container_registry.this.name)
}
