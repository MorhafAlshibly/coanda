# Create azure public ip address
resource "azurerm_public_ip" "this" {
  name                = var.name
  location            = var.location
  resource_group_name = var.resource_group_name
  allocation_method   = "Static"

  sku = "Standard"

  tags = {
    environment = var.environment
  }
}
