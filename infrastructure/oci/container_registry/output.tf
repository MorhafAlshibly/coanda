# Registry uri
output "registry_uri" {
  value = null_resource.docker_compose.triggers["registry_uri"]
}
