# ARN
output "arn" {
  value = aws_ecr_repository.this.arn
}

# Registry ID
output "registry_id" {
  value = aws_ecr_repository.this.registry_id
}

# Registry URL
output "registry_url" {
  value = aws_ecr_repository.this.repository_url
}
