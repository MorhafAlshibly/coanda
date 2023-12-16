# Parent compartment OCID
variable "parent_compartment_id" {
  type = string
}

# Environment
variable "environment" {
  type = string
}

# Region
variable "region" {
  type = string
}

# Name of the application
variable "app_name" {
  type = string
}

# Tenancy OCID
variable "tenancy_ocid" {
  type = string
}

# User OCID
variable "user_ocid" {
  type = string
}

# Authentication token
variable "auth_token" {
  type      = string
  sensitive = true
}

# Admin password
variable "admin_password" {
  type      = string
  sensitive = true
}
