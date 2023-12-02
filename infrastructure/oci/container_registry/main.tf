# Create an oracle container registry instance
resource "oci_artifacts_container_repository" "this" {
  compartment_id = var.compartment_id
  display_name   = var.name

  freeform_tags = {
    "environment" : var.environment
  }
  is_immutable = false
  is_public    = true

}
