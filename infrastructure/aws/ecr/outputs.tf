# ARN glob
output "arn" {
  value = replace(aws_ecr_repository.this[keys(aws_ecr_repository.this)[0]].arn, keys(aws_ecr_repository.this)[0], "*")
}
