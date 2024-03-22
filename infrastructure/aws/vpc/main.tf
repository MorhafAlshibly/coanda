resource "aws_vpc" "this" {
  tags       = merge(var.tags, { Name = var.name })
  cidr_block = "10.0.0.0/16"
}

resource "aws_subnet" "public" {
  for_each          = { for i, zone in var.availability_zones : zone => zone }
  vpc_id            = aws_vpc.this.id
  cidr_block        = cidrsubnet(aws_vpc.this.cidr_block, 8, 2 * index(var.availability_zones, each.key))
  availability_zone = each.value
  tags              = merge(var.tags, { Name = format("%s-%s-public", var.name, each.value) })
}

resource "aws_subnet" "private" {
  for_each   = { for zone in var.availability_zones : zone => zone }
  vpc_id     = aws_vpc.this.id
  cidr_block = cidrsubnet(aws_vpc.this.cidr_block, 8, (2 * index(var.availability_zones, each.key)) + 1)
  tags       = merge(var.tags, { Name = format("%s-%s-private", var.name, each.value) })
}

resource "aws_route_table" "public" {
  vpc_id = aws_vpc.this.id
  tags   = merge(var.tags, { Name = format("%s-public", var.name) })
}

resource "aws_route_table" "private" {
  for_each = aws_subnet.private
  vpc_id   = aws_vpc.this.id
  tags     = merge(var.tags, { Name = format("%s-%s-private", var.name, each.key) })
}

resource "aws_route_table_association" "public" {
  for_each       = aws_subnet.public
  subnet_id      = each.value.id
  route_table_id = aws_route_table.public.id
}

resource "aws_route_table_association" "private" {
  for_each       = aws_subnet.private
  subnet_id      = each.value.id
  route_table_id = aws_route_table.private[each.key].id
}

resource "aws_security_group" "this" {
  vpc_id = aws_vpc.this.id
  tags   = merge(var.tags, { Name = format("%s-security-group", var.name) })
}

resource "aws_vpc_security_group_egress_rule" "ipv4" {
  security_group_id = aws_security_group.this.id
  ip_protocol       = "-1"
  cidr_ipv4         = "0.0.0.0/0"
}

resource "aws_vpc_security_group_egress_rule" "ipv6" {
  security_group_id = aws_security_group.this.id
  ip_protocol       = "-1"
  cidr_ipv6         = "::/0"
}

resource "aws_vpc_security_group_ingress_rule" "ipv4" {
  security_group_id = aws_security_group.this.id
  ip_protocol       = "-1"
  cidr_ipv4         = "0.0.0.0/0"
}

resource "aws_vpc_security_group_ingress_rule" "ipv6" {
  security_group_id = aws_security_group.this.id
  ip_protocol       = "-1"
  cidr_ipv6         = "::/0"
}

resource "aws_vpc_security_group_ingress_rule" "bff_ipv4" {
  security_group_id = aws_security_group.this.id
  from_port         = 8080
  to_port           = 8080
  ip_protocol       = "tcp"
  cidr_ipv4         = "0.0.0.0/0"
  description       = "Allow ipv4 traffic from the internet to the BFF"
}

resource "aws_vpc_security_group_ingress_rule" "bff_ipv6" {
  security_group_id = aws_security_group.this.id
  from_port         = 8080
  to_port           = 8080
  ip_protocol       = "tcp"
  cidr_ipv6         = "::/0"
  description       = "Allow ipv6 traffic from the internet to the BFF"
}

resource "aws_internet_gateway" "this" {
  vpc_id = aws_vpc.this.id
  tags   = merge(var.tags, { Name = format("%s-internet-gateway", var.name) })
}

resource "aws_route" "public_ipv4" {
  route_table_id         = aws_route_table.public.id
  destination_cidr_block = "0.0.0.0/0"
  gateway_id             = aws_internet_gateway.this.id
}

resource "aws_route" "public_ipv6" {
  route_table_id              = aws_route_table.public.id
  destination_ipv6_cidr_block = "::/0"
  gateway_id                  = aws_internet_gateway.this.id
}
