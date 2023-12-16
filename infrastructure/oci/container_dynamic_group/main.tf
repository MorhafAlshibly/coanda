# Create a dynamic group to allow the container instance to pull images from the registry
resource "oci_identity_dynamic_group" "this" {
  compartment_id = var.tenancy_ocid
  description    = "Dynamic group for container instances"
  matching_rule  = "ALL {resource.type='computecontainerinstance'}"
  name           = var.name
  freeform_tags = {
    environment = var.environment
  }
}

# Create a policy to allow the dynamic group to pull images from the registry
resource "oci_identity_policy" "this" {
  compartment_id = var.tenancy_ocid
  description    = "Policy for container instances"
  name           = format("policy-%s", var.name)
  statements = [
    format("Allow dynamic-group %s to read repos in tenancy", oci_identity_dynamic_group.this.name),
  ]
}
