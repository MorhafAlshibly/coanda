# Elastic container registry
resource "aws_ecr_repository" "this" {
  for_each             = { for container in var.containers : container.name => container }
  name                 = each.value.repository_name
  image_tag_mutability = "MUTABLE"

  image_scanning_configuration {
    scan_on_push = true
  }
  tags = var.tags
}

# Have to define the docker provider here as it is third party
terraform {
  required_providers {
    docker = {
      source  = "kreuzwerker/docker"
      version = "3.0.2"
    }
  }
}

# Build docker images
resource "docker_image" "this" {
  for_each   = { for container in var.containers : container.name => container }
  depends_on = [aws_ecr_repository.this]
  name       = each.value.endpoint
  build {
    context = format("%s/../../", path.cwd)
    build_args = {
      PORT    = each.value.port
      SERVICE = each.value.name
    }
  }
}

# Push docker images
resource "docker_registry_image" "this" {
  for_each   = { for container in var.containers : container.name => container }
  depends_on = [docker_image.this]
  name       = each.value.endpoint
}
