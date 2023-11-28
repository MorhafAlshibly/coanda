# Resource group name
output "resource_group_name" {
  value = azurerm_resource_group.this.name
}

# Storage account name
output "storage_account_name" {
  value = module.storage_account.name
}

# Container registry uri
output "container_registry_uri" {
  value = module.container_registry.uri
}
