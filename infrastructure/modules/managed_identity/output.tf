output "principal_id" {
  value = azurerm_user_assigned_identity.this.principal_id
}

output "id" {
  value = data.azurerm_user_assigned_identity.after_acr_pull.id
}
