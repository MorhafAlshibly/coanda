# Create a oracle container instance
resource "oci_container_instances_container_instance" "this" {
  availability_domain = var.availability_domain
  compartment_id      = var.compartment_id
  display_name        = var.name

  freeform_tags = {
    "environment" : var.environment
  }



  containers {
    image_url = format("%s/item:latest", var.registry_uri)
    environment_variables = {
      "ITEM_MONGOOVERTABLE" : true,
      "ITEM_MONGOCONN" : "mongodb://localhost:27017",
    }
  }

  shape = "CI.Standard.A1.Flex"

  shape_config {
    ocpus         = 4
    memory_in_gbs = 24
  }

  vnics {
    subnet_id             = var.subnet_id
    is_public_ip_assigned = true
  }

  state = "ACTIVE"
}
