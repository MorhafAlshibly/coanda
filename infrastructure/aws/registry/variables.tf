# Name
variable "name" {
  type = string
}

# Endpoint
variable "endpoint" {
  type = string
}

# Containers
variable "containers" {
  type = list(object({
    name = string
  }))
}

# Tags
variable "tags" {
  type = map(string)
}
