# Create a virtual network
resource "azurerm_virtual_network" "this" {
  name                = var.name
  location            = var.location
  resource_group_name = var.resource_group_name
  address_space       = ["10.0.0.0/16"]

  tags = {
    environment = var.environment
  }
}

# Create a subnet
resource "azurerm_subnet" "this" {
  name                 = "internal"
  resource_group_name  = var.resource_group_name
  virtual_network_name = azurerm_virtual_network.this.name
  address_prefixes     = ["10.0.0.0/18"]
  service_endpoints    = ["Microsoft.AzureCosmosDB", "Microsoft.KeyVault", "Microsoft.Storage"]
}
