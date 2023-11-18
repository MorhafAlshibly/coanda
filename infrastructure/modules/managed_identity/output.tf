output "principal_id" {
  value = azurerm_user_assigned_identity.this.principal_id
}

output "id" {
  value = data.azurerm_user_assigned_identity.after_permissions.id
}

output "client_id" {
  value = data.azurerm_user_assigned_identity.after_permissions.client_id
}

output "tenant_id" {
  value = data.azurerm_user_assigned_identity.after_permissions.tenant_id
}
