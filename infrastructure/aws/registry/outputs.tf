# # ARN
# output "arn" {
#   depends_on = [docker_registry_image.this]
#   value      = aws_ecr_repository.this.arn
# }

# # Registry ID
# output "registry_id" {
#   depends_on = [docker_registry_image.this]
#   value      = aws_ecr_repository.this.registry_id
# }

# # Registry URL
# output "registry_url" {
#   depends_on = [docker_registry_image.this]
#   value      = aws_ecr_repository.this.repository_url
# }
