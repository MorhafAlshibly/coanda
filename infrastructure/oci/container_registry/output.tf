# Registry uri
output "registry_uri" {
  depends_on = [null_resource.docker_compose]
  value      = null_resource.docker_compose.triggers["registry_uri"]
}
