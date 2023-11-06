# Key vault id
output "id" {
  value = azurerm_key_vault.this.id
}

# Vault URI
output "vault_uri" {
  value = format("https://%s.vault.azure.net/", azurerm_key_vault.this.name)
}
