# Key vault id
output "id" {
  value = azurerm_key_vault.this.id
}

# URI
output "uri" {
  value = azurerm_key_vault.this.vault_uri
}

# MongoDB connection string secret name
output "mongo_connection_secret_name" {
  value = azurerm_key_vault_secret.mongo_connection_string.name
}
