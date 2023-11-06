# Storage account uri
output "storage_uri" {
  value = format("https://%s.table.core.windows.net", azurerm_storage_account.this.name)
}
