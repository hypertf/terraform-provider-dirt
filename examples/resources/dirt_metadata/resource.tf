resource "dirt_metadata" "app_config" {
  path  = "app/config/database_url"
  value = "postgres://user:pass@localhost:5432/mydb"
}

resource "dirt_metadata" "feature_flag" {
  path  = "features/new_ui_enabled"
  value = "true"
}

resource "dirt_metadata" "version" {
  path  = "app/version"
  value = "1.2.3"
}

# Example using dynamic values from other resources
resource "dirt_metadata" "instance_info" {
  path  = "instances/${dirt_instance.web_server.id}/description"
  value = "Web server instance with ${dirt_instance.web_server.cpu} CPUs and ${dirt_instance.web_server.memory_mb}MB RAM"
}

# Reference resources from other examples
resource "dirt_project" "example" {
  name = "metadata-example-project"
}

resource "dirt_instance" "web_server" {
  project_id = dirt_project.example.id
  name       = "web-server-metadata"
  cpu        = 2
  memory_mb  = 2048
}

# Outputs to show the metadata ID and other attributes
output "app_config_id" {
  description = "ID of the app config metadata"
  value       = dirt_metadata.app_config.id
}

output "app_config_path" {
  description = "Path of the app config metadata"
  value       = dirt_metadata.app_config.path
}

output "feature_flag_metadata" {
  description = "Complete feature flag metadata information"
  value = {
    id         = dirt_metadata.feature_flag.id
    path       = dirt_metadata.feature_flag.path
    value      = dirt_metadata.feature_flag.value
    created_at = dirt_metadata.feature_flag.created_at
    updated_at = dirt_metadata.feature_flag.updated_at
  }
}