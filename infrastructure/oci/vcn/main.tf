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
# resource "oci_core_network_security_group" "this" {
#   compartment_id = var.compartment_id
#   vcn_id         = oci_core_vcn.this.id
#   display_name   = format("nsg-%s", var.name)
#   freeform_tags = {
#     environment = var.environment
#   }
# }


# Create a private subnet
resource "oci_core_subnet" "private" {
  compartment_id            = var.compartment_id
  vcn_id                    = oci_core_vcn.this.id
  cidr_block                = "10.0.1.0/24"
  display_name              = format("private-subnet-%s", var.name)
  dns_label                 = "private"
  security_list_ids         = [oci_core_security_list.private.id]
  route_table_id            = oci_core_route_table.private.id
  prohibit_internet_ingress = true
  freeform_tags = {
    environment = var.environment
  }
}

# Create a public subnet
resource "oci_core_subnet" "public" {
  compartment_id    = var.compartment_id
  vcn_id            = oci_core_vcn.this.id
  cidr_block        = "10.0.0.0/24"
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
    stateless = false
    source    = "0.0.0.0/0"
    tcp_options {
      max = 22
      min = 22
    }
  }
  // Allow for ingress to ICMP port for error messages
  ingress_security_rules {
    protocol  = "1"
    stateless = false
    source    = "0.0.0.0/0"
    icmp_options {
      type = 3
      code = 4
    }
  }
  ingress_security_rules {
    protocol  = "1"
    stateless = false
    source    = "10.0.0.0/16"
    icmp_options {
      type = 3
    }
  }
  // Allow ingress traffic to HTTP port 80
  ingress_security_rules {
    protocol  = "6"
    stateless = false
    source    = "0.0.0.0/0"
    tcp_options {
      max = 80
      min = 80
    }
  }
  // Allow ingress traffic to HTTPS port 443
  ingress_security_rules {
    protocol  = "6"
    stateless = false
    source    = "0.0.0.0/0"
    tcp_options {
      max = 443
      min = 443
    }
  }
  // Allow eggress traffic to all ports
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
  display_name   = format("internet-gateway-%s", var.name)
  enabled        = true
  freeform_tags = {
    environment = var.environment
  }
}

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
  route_rules {
    destination       = "all-lhr-services-in-oracle-services-network"
    destination_type  = "SERVICE_CIDR_BLOCK"
    network_entity_id = oci_core_service_gateway.this.id
  }
}


# Create a NAT gateway for private subnet
resource "oci_core_nat_gateway" "this" {
  compartment_id = var.compartment_id
  vcn_id         = oci_core_vcn.this.id
  display_name   = format("nat-gateway-%s", var.name)
  freeform_tags = {
    environment = var.environment
  }
  block_traffic = false
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


# All LHR services in Oracle Services Network
data "oci_core_services" "this" {}

# Create a service gateway
resource "oci_core_service_gateway" "this" {
  compartment_id = var.compartment_id
  vcn_id         = oci_core_vcn.this.id
  display_name   = format("service-gateway-%s", var.name)
  freeform_tags = {
    environment = var.environment
  }
  services {
    service_id = data.oci_core_services.this.services[index(data.oci_core_services.this.services.*.cidr_block, "all-lhr-services-in-oracle-services-network")].id
  }
}
