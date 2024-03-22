# Name
variable "name" {
  type = string
}

# Availability zones
variable "availability_zones" {
  type = list(string)
}

# Tags
variable "tags" {
  type = map(string)
}
