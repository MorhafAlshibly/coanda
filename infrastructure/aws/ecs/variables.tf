# Name
variable "name" {
  type = string
}

# Tags
variable "tags" {
  type = map(string)
}

# Containers
variable "containers" {
  type = list(object({
    name                 = string
    port                 = number
    host_port            = number
    endpoint             = string
    task_definition_name = string
    environment          = map(string)
    assign_public_ip     = bool
    public               = bool
  }))
}

# Repository ARN
variable "repository_arn" {
  type = string
}

# VPC ID
variable "vpc_id" {
  type = string
}

# Public subnets
variable "public_subnets" {
  type = list(string)
}

# Private subnets
variable "private_subnets" {
  type = list(string)
}

# Security groups
variable "security_groups" {
  type = list(string)
}
