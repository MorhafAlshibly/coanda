# Create an oracle container registry instance
resource "oci_artifacts_container_repository" "this" {
  compartment_id = var.compartment_id
  display_name   = var.name

  freeform_tags = {
    "environment" : var.environment
  }
  is_immutable = false
  is_public    = false

}

# Docker compose run
resource "null_resource" "docker_compose" {
  depends_on = [oci_artifacts_container_repository.this]
  triggers = {
    registry_uri = format("%s.ocir.io/%s/%s", var.region, var.namespace, var.name)
  }
  provisioner "local-exec" {
    command = format("task oci:push ENV=%s NAMESPACE=%s REPO_NAME=%s USERNAME=%s", var.environment, var.namespace, var.name, var.username)
  }
}
