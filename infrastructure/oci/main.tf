# S3 OCI provider
terraform {
  backend "s3" {
    skip_region_validation      = true
    skip_credentials_validation = true
    force_path_style            = true
    skip_metadata_api_check     = true
    endpoint                    = var.endpoint
    bucket                      = var.bucket
    key                         = var.key
    region                      = var.region
    access_key                  = var.access_key
    secret_key                  = var.secret_key
  }
  required_providers {
    oci = {
      source  = "oracle/oci"
      version = ">= 4.38.0"
    }
  }
}

provider "oci" {
  auth   = "APIKey"
  region = var.region
}

# Object storage namespace
data "oci_objectstorage_namespace" "this" {
  compartment_id = var.parent_compartment_id
}

# User details for the authentication
data "oci_identity_user" "this" {
  user_id = var.user_ocid
}

# Compartment
resource "oci_identity_compartment" "this" {
  name           = format("compartment-%s-%s-%s", var.app_name, var.environment, var.region)
  description    = format("Compartment for the %s application in the %s environment", var.app_name, var.environment)
  compartment_id = var.parent_compartment_id

  freeform_tags = {
    environment = var.environment
  }
}

# Include the module that creates a container registry
module "container_registry" {
  source         = "./container_registry"
  compartment_id = oci_identity_compartment.this.id
  name           = format("cr-%s-%s-%s", var.app_name, var.environment, var.region)
  environment    = var.environment
  region         = var.region
  namespace      = data.oci_objectstorage_namespace.this.namespace
  username       = data.oci_identity_user.this.name
  registries     = ["bff", "item", "record", "team"]
}

# Include the module that creates a virtual cloud network
module "vcn" {
  source         = "./vcn"
  compartment_id = oci_identity_compartment.this.id
  name           = format("vcn-%s-%s-%s", var.app_name, var.environment, var.region)
  environment    = var.environment
}

# Include the module that creates a container instance
module "container_instance" {
  depends_on              = [module.container_dynamic_group]
  source                  = "./container_instance"
  compartment_id          = oci_identity_compartment.this.id
  name                    = format("ci-%s-%s-%s", var.app_name, var.environment, var.region)
  environment             = var.environment
  subnet_id               = module.vcn.public_subnet_id
  registry_uri            = module.container_registry.registry_uri
  mongo_connection_string = module.database.mongo_connection_string
}

# Include the module that creates a database
module "database" {
  source         = "./database"
  compartment_id = oci_identity_compartment.this.id
  name           = format("db%s%s%s", var.app_name, var.environment, replace(var.region, "-", ""))
  environment    = var.environment
  admin_password = var.admin_password
  vcn_id         = module.vcn.id
}

# Include the module that creates a container dynamic group
module "container_dynamic_group" {
  source       = "./container_dynamic_group"
  tenancy_ocid = var.tenancy_ocid
  name         = format("cdg-%s-%s-%s", var.app_name, var.environment, var.region)
  environment  = var.environment
}
