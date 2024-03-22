# VPC ID
output "vpc_id" {
  value = aws_vpc.this.id
}

# Public subnet IDs
output "public_subnet_ids" {
  value = values(aws_subnet.public)[*].id
}

# Private subnet IDs
output "private_subnet_ids" {
  value = values(aws_subnet.private)[*].id
}

# Public route table ID
output "public_route_table_id" {
  value = aws_route_table.public.id
}

# Private route table IDs
output "private_route_table_ids" {
  value = values(aws_route_table.private)[*].id
}

# Security group ID
output "security_group_id" {
  value = aws_security_group.this.id
}

# Internet gateway ID
output "internet_gateway_id" {
  value = aws_internet_gateway.this.id
}
