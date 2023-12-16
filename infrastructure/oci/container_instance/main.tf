# Availability domains
data "oci_identity_availability_domains" "this" {
  compartment_id = var.compartment_id
}

# Create a oracle container instance
resource "oci_container_instances_container_instance" "this" {
  availability_domain = data.oci_identity_availability_domains.this.availability_domains[0].name
  compartment_id      = var.compartment_id
  display_name        = var.name

  freeform_tags = {
    environment = var.environment
  }

  containers {
    display_name = "item"
    image_url    = format("%s/item:latest", var.registry_uri)
    environment_variables = {
      "ITEM_MONGOOVERTABLE" : true,
      "ITEM_MONGOCONN" : var.mongo_connection_string,
    }
  }

  containers {
    display_name = "record"
    image_url    = format("%s/record:latest", var.registry_uri)
    environment_variables = {
      "RECORD_MONGOCONN" : var.mongo_connection_string,
    }
  }

  containers {
    display_name = "team"
    image_url    = format("%s/team:latest", var.registry_uri)
    environment_variables = {
      "TEAM_MONGOCONN" : var.mongo_connection_string,
    }
  }

  containers {
    display_name = "bff"
    image_url    = format("%s/bff:latest", var.registry_uri)
    environment_variables = {
      "BFF_ENABLEPLAYGROUND" : var.environment == "dev" ? true : false,
      "BFF_PORT" : 80,
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

  container_restart_policy = "ALWAYS"
  state                    = "ACTIVE"
}
