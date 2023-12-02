# Container registry uri
output "uri" {
  value = data.azurerm_container_registry.after_docker_compose.login_server
}

# Container registry id
output "id" {
  value = data.azurerm_container_registry.after_docker_compose.id
}
