# The app configuration id
output "id" {
  value = azurerm_app_configuration.this.id
}

# The app configuration endpoint
output "endpoint" {
  value = azurerm_app_configuration.this.endpoint
}
