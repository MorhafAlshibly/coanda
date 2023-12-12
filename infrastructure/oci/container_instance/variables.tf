# Compartment OCID
variable "compartment_id" {
  type = string
}

# Registry URI
variable "registry_uri" {
  type = string
}

# Environment
variable "environment" {
  type = string
}

# VCN subnet id
variable "subnet_id" {
  type = string
}

# Name
variable "name" {
  type = string
}

# Username
variable "username" {
  type = string
}

# Namespace
variable "namespace" {
  type = string
}

# Password
variable "password" {
  type      = string
  sensitive = true
}

# Mongo connection string
variable "mongo_connection_string" {
  type      = string
  sensitive = true
}
