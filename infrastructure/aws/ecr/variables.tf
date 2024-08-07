# Name
variable "name" {
  type = string
}

# Containers
variable "containers" {
  type = list(object({
    name            = string
    endpoint        = string
    repository_name = string
    port            = number
  }))
}

# Tags
variable "tags" {
  type = map(string)
}
