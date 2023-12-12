# MongoDB connection string
output "mongo_connection_string" {
  value     = replace(replace(data.external.get_mongo_connection_string.result.mongo_connection_string, "[user:password@]", format("ADMIN:%s@", var.admin_password)), "[user]", "ADMIN")
  sensitive = true
}
