# Managed identity principal id
output "principal_id" {
  value = azurerm_user_assigned_identity.this.principal_id
}

# Managed identity id
output "id" {
  value = azurerm_user_assigned_identity.this.id
}

# Managed identity client id
output "client_id" {
  value = azurerm_user_assigned_identity.this.client_id
}

# Managed identity tenant id
output "tenant_id" {
  value = azurerm_user_assigned_identity.this.tenant_id
}
