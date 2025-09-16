data "dirt_metadata" "app_config" {
  path = "app/config/database_url"
}

data "dirt_metadata" "feature_flag" {
  path = "app/features/new_ui_enabled"
}

output "database_url" {
  description = "Database URL from metadata"
  value       = data.dirt_metadata.app_config.value
  sensitive   = true
}

output "database_url_path" {
  description = "Path of the database URL metadata"
  value       = data.dirt_metadata.app_config.path
}

output "new_ui_enabled" {
  description = "Whether the new UI feature is enabled"
  value       = data.dirt_metadata.feature_flag.value
}

output "feature_flag_details" {
  description = "Complete feature flag metadata details"
  value = {
    id         = data.dirt_metadata.feature_flag.id
    path       = data.dirt_metadata.feature_flag.path
    value      = data.dirt_metadata.feature_flag.value
    created_at = data.dirt_metadata.feature_flag.created_at
    updated_at = data.dirt_metadata.feature_flag.updated_at
  }
}

# Example using metadata value in conditional logic
resource "dirt_instance" "conditional_instance" {
  count      = data.dirt_metadata.feature_flag.value == "true" ? 1 : 0
  project_id = dirt_project.example.id
  name       = "feature-test-instance"
  cpu        = 2
  memory_mb  = 2048
}

# Example using metadata for configuration
resource "dirt_metadata" "derived_config" {
  path  = "app/config/full_database_config"
  value = "Database: ${data.dirt_metadata.app_config.value}, Pool: 10"
}

# Reference project for the conditional instance
resource "dirt_project" "example" {
  name = "metadata-datasource-example"
}