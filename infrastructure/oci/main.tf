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

# Compartment
resource "oci_identity_compartment" "this" {
  name           = format("compartment-%s-%s-%s", var.app_name, var.environment, var.region)
  description    = format("Compartment for the %s application in the %s environment", var.app_name, var.environment)
  compartment_id = var.parent_compartment_id

  freeform_tags = {
    "environment" : var.environment
  }
}

# Object storage namespace
data "oci_objectstorage_namespace" "this" {
  compartment_id = oci_identity_compartment.this.id
}

# Include the module that creates a container registry
module "container_registry" {
  source         = "./container_registry"
  compartment_id = oci_identity_compartment.this.id
  name           = format("cr-%s-%s-%s", var.app_name, var.environment, var.region)
  environment    = var.environment
  region         = var.region
  namespace      = data.oci_objectstorage_namespace.this.namespace
}
