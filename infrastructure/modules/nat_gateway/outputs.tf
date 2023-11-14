# NAT gateway ID
output "id" {
  value = data.azurerm_nat_gateway.after_associations.id
}
