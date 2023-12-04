# Private subnet id
output "private_subnet_id" {
  value = oci_core_subnet.private.id
}

# Public subnet id
output "public_subnet_id" {
  value = oci_core_subnet.public.id
}
