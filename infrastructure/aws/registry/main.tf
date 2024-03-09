# Container registry
resource "aws_ecr_repository" "this" {
  for_each             = { for container in var.containers : container.name => container }
  name                 = format("%s-%s", var.name, each.value.name)
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
  name       = format("%s/%s:latest", replace(var.endpoint, "https://", ""), format("%s-%s", var.name, each.value.name))
  build {
    context    = format("%s/../../", path.cwd)
    dockerfile = format("./cmd/%s/Dockerfile", each.value.name)
    build_args = {
      "SERVICE" = each.value.name
    }
  }
}

# Push docker images
resource "docker_registry_image" "this" {
  for_each   = { for container in var.containers : container.name => container }
  depends_on = [docker_image.this]
  name       = format("%s/%s:latest", replace(var.endpoint, "https://", ""), format("%s-%s", var.name, each.value.name))
}
