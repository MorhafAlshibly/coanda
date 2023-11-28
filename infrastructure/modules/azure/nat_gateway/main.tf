# Create a NAT gateway
resource "azurerm_nat_gateway" "this" {
  name                = var.name
  location            = var.location
  resource_group_name = var.resource_group_name
  sku_name            = "Standard"

  tags = {
    environment = var.environment
  }
}

# Create a NAT gateway public ip association
resource "azurerm_nat_gateway_public_ip_association" "this" {
  nat_gateway_id       = azurerm_nat_gateway.this.id
  public_ip_address_id = var.ip_address_id
}

# Create a NAT gateway subnet association
resource "azurerm_subnet_nat_gateway_association" "this" {
  nat_gateway_id = azurerm_nat_gateway.this.id
  subnet_id      = var.vnet_subnet_id
}

# Get NAT gateway id and output after associations
data "azurerm_nat_gateway" "after_associations" {
  depends_on          = [azurerm_nat_gateway_public_ip_association.this, azurerm_subnet_nat_gateway_association.this]
  name                = azurerm_nat_gateway.this.name
  resource_group_name = azurerm_nat_gateway.this.resource_group_name
}
