# Create an Autonomous JSON Database with MongoDB API
resource "oci_database_autonomous_database" "this" {
  compartment_id = var.compartment_id
  admin_password = var.admin_password
  db_name        = var.name
  db_version     = "21c"
  db_workload    = "AJD"
  display_name   = var.name
  freeform_tags = {
    environment = var.environment
  }
  is_free_tier    = true
  license_model   = "LICENSE_INCLUDED"
  whitelisted_ips = [var.vcn_id]
}

# Get the connection string for the Autonomous JSON Database
data "external" "get_mongo_connection_string" {
  depends_on = [oci_database_autonomous_database.this]
  program    = ["task", "oci:mongo", "OCID=${oci_database_autonomous_database.this.id}"]
  query = {
    mongo_connection_string = "mongo_connection_string"
  }
}
