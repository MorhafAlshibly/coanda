# Create the container registries
resource "oci_artifacts_container_repository" "this" {
  for_each       = var.registries
  compartment_id = var.compartment_id
  display_name   = format("%s/%s", var.name, each.key)

  freeform_tags = {
    environment = var.environment
  }
  is_immutable = false
  is_public    = false
}

# Docker compose run
resource "null_resource" "docker_compose" {
  depends_on = [oci_artifacts_container_repository.this]
  triggers = {
    registry_uri = format("%s.ocir.io/%s/%s", var.region, var.namespace, var.name)
    registry_id  = oci_artifacts_container_repository.this[keys(oci_artifacts_container_repository.this)[0]].id
  }
  provisioner "local-exec" {
    command = format("task oci:push ENV=%s NAMESPACE=%s REPO_NAME=%s USERNAME=%s PLATFORM=%s", var.environment, var.namespace, var.name, var.username, "linux/arm64")
  }
}
