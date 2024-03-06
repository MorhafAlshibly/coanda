# build docker image
resource "docker_image" "this" {
  name = format("%s/%s:latest", data.aws_ecr_authorization_token.token.proxy_endpoint, var.name)
  build {
    context = "."
  }
  platform = "linux/arm64"
}

# push image to ecr repo
resource "docker_registry_image" "this" {
  name = docker_image.this.name
}
