# Private subnet id
output "private_subnet_id" {
  value = oci_core_subnet.private.id
}

# Public subnet id
output "public_subnet_id" {
  value = oci_core_subnet.public.id
}

# VCN id
output "id" {
  value = oci_core_vcn.this.id
}
