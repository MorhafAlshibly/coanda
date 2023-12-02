# S3 OCI provider
terraform {
  backend "s3" {
    skip_region_validation      = true
    skip_credentials_validation = true
    skip_requesting_account_id  = true
    use_path_style              = true
    skip_metadata_api_check     = true
    endpoints                   = { s3 = "https://lrk70fbaaokt.compat.objectstorage.uk-london-1.oraclecloud.com" }
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
}


# Compartment
resource "oci_identity_compartment" "this" {
  name           = format("compartment-%s-%s-%s", var.app_name, var.environment, var.location)
  description    = "Compartment for the ${var.app_name} application"
  compartment_id = var.parent_compartment_id
}

# Include the module that creates a container registry
module "container_registry" {
  source         = "./container_registry"
  compartment_id = oci_identity_compartment.this.id
  name           = format("acr%s%s%s", var.app_name, var.environment, var.location)
  environment    = var.environment
}
