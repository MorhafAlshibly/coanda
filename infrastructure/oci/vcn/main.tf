# Create a vnet
resource "oci_core_vcn" "this" {
  compartment_id = var.compartment_id
  cidr_blocks    = ["10.0.0.0/16"]
  display_name   = var.name
  freeform_tags = {
    "environment" : var.environment
  }
  dns_label = substr(replace(var.name, "-", ""), 0, 15)
}


# Create a network security group for the vcn
resource "oci_core_network_security_group" "this" {
  compartment_id = var.compartment_id
  vcn_id         = oci_core_vcn.this.id
  display_name   = format("nsg-%s", var.name)
  freeform_tags = {
    environment = var.environment
  }
}


# Create a private subnet
resource "oci_core_subnet" "private" {
  depends_on        = [oci_core_network_security_group.this]
  compartment_id    = var.compartment_id
  vcn_id            = oci_core_vcn.this.id
  cidr_block        = "10.0.0.0/24"
  display_name      = format("private-subnet-%s", var.name)
  dns_label         = "private"
  security_list_ids = [oci_core_security_list.private.id]
  route_table_id    = oci_core_route_table.private.id
  freeform_tags = {
    environment = var.environment
  }
}

# Create a public subnet
resource "oci_core_subnet" "public" {
  depends_on        = [oci_core_network_security_group.this]
  compartment_id    = var.compartment_id
  vcn_id            = oci_core_vcn.this.id
  cidr_block        = "10.0.1.0/24"
  display_name      = format("public-subnet-%s", var.name)
  dns_label         = "public"
  security_list_ids = [oci_core_security_list.public.id]
  route_table_id    = oci_core_route_table.public.id
  freeform_tags = {
    environment = var.environment
  }
}

# Create a public security list
resource "oci_core_security_list" "public" {
  compartment_id = var.compartment_id
  vcn_id         = oci_core_vcn.this.id
  display_name   = format("public-security-list-%s", var.name)
  freeform_tags = {
    environment = var.environment
  }
  // Allow all ingress traffic to SSH port 22
  ingress_security_rules {
    protocol  = "6"
    source    = "0.0.0.0/0"
    stateless = false
    tcp_options {
      source_port_range {
        max = 22
        min = 22
      }
    }
  }
  // Allow for ingress to ICMP port for error messages
  ingress_security_rules {
    protocol  = "1"
    source    = "0.0.0.0/0"
    stateless = false
    icmp_options {
      type = 3
      code = 4
    }
  }
  ingress_security_rules {
    protocol  = "1"
    source    = "10.0.0.0/16"
    stateless = false
    icmp_options {
      type = 3
    }
  }
  // Allow all egress traffic
  egress_security_rules {
    protocol    = "all"
    stateless   = false
    destination = "0.0.0.0/0"
  }
}

# Create a private security list
resource "oci_core_security_list" "private" {
  compartment_id = var.compartment_id
  vcn_id         = oci_core_vcn.this.id
  display_name   = format("private-security-list-%s", var.name)
  freeform_tags = {
    environment = var.environment
  }
  // Allow all ingress traffic to SSH port 22
  ingress_security_rules {
    protocol  = "6"
    source    = "10.0.0.0/16"
    stateless = false
    tcp_options {
      source_port_range {
        max = 22
        min = 22
      }
    }
  }
  // Allow for ingress to ICMP port for error messages
  ingress_security_rules {
    protocol  = "1"
    source    = "0.0.0.0/0"
    stateless = false
    icmp_options {
      type = 3
      code = 4
    }
  }
  ingress_security_rules {
    protocol  = "1"
    source    = "10.0.0.0/16"
    stateless = false
    icmp_options {
      type = 3
    }
  }
  // Allow all egress traffic
  egress_security_rules {
    protocol    = "all"
    stateless   = false
    destination = "0.0.0.0/0"
  }
}

# Create a internet gateway
resource "oci_core_internet_gateway" "this" {
  compartment_id = var.compartment_id
  vcn_id         = oci_core_vcn.this.id
  display_name   = format("gateway-%s", var.name)
  enabled        = true
  freeform_tags = {
    environment = var.environment
  }
}

# Set the route table for the internet gateway
resource "null_resource" "gateway_set_route_table" {
  triggers = {
    route_table_id = oci_core_route_table.public.id
  }
  provisioner "local-exec" {
    command = format("oci network internet-gateway update --ig-id %s --route-table-id %s", oci_core_internet_gateway.this.id, oci_core_route_table.public.id)
  }
}

# Create a public route table
resource "oci_core_route_table" "public" {
  compartment_id = var.compartment_id
  vcn_id         = oci_core_vcn.this.id
  display_name   = format("public-route-table-%s", var.name)
  freeform_tags = {
    environment = var.environment
  }
  route_rules {
    destination       = "0.0.0.0/0"
    destination_type  = "CIDR_BLOCK"
    network_entity_id = oci_core_internet_gateway.this.id
  }
}

# Create a NAT gateway for private subnet
resource "oci_core_nat_gateway" "this" {
  compartment_id = var.compartment_id
  vcn_id         = oci_core_vcn.this.id
  display_name   = format("nat-%s", var.name)
  freeform_tags = {
    environment = var.environment
  }
  block_traffic = false
}

# Create a service gateway
# resource "oci_core_service_gateway" "this" {
#   compartment_id = var.compartment_id
#   vcn_id         = oci_core_vcn.this.id
#   display_name   = format("service-gateway-%s", var.name)
#   freeform_tags = {
#     environment = var.environment
#   }
#   route_table_id = oci_core_route_table.private.id
#   services {
#     service_id = ""
#   }
# }

# Create a private route table
resource "oci_core_route_table" "private" {
  compartment_id = var.compartment_id
  vcn_id         = oci_core_vcn.this.id
  display_name   = format("private-route-table-%s", var.name)
  freeform_tags = {
    environment = var.environment
  }
  route_rules {
    destination       = "0.0.0.0/0"
    destination_type  = "CIDR_BLOCK"
    network_entity_id = oci_core_nat_gateway.this.id
  }
}

# Set the route table for the NAT gateway
resource "null_resource" "nat_set_route_table" {
  triggers = {
    route_table_id = oci_core_route_table.private.id
  }
  provisioner "local-exec" {
    command = format("oci network nat-gateway update --nat-gateway-id %s --route-table-id %s", oci_core_nat_gateway.this.id, oci_core_route_table.private.id)
  }
}
