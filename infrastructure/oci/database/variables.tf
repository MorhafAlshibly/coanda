# Display name
variable "name" {
  type = string
}

# Environment
variable "environment" {
  type = string
}

# Compartment OCID
variable "compartment_id" {
  type = string
}

# Admin password
variable "admin_password" {
  type      = string
  sensitive = true
}

# VCN OCID
variable "vcn_id" {
  type = string
}
