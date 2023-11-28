# Storage account name
output "name" {
  value = azurerm_storage_account.this.name
}

# Storage account id
output "id" {
  value = azurerm_storage_account.this.id
}
