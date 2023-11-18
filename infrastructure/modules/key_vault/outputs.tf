# Key vault id
output "id" {
  value = azurerm_key_vault.this.id
}

# URI
output "uri" {
  value = azurerm_key_vault.this.vault_uri
}
