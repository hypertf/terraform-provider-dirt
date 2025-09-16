data "dirt_project" "existing" {
  id = "project-id-12345"
}

output "project_name" {
  description = "Name of the existing project"
  value       = data.dirt_project.existing.name
}

output "project_created_at" {
  description = "When the project was created"
  value       = data.dirt_project.existing.created_at
}

output "project_updated_at" {
  description = "When the project was last updated"
  value       = data.dirt_project.existing.updated_at
}

# Example using the data source to create related resources
resource "dirt_instance" "app_server" {
  project_id = data.dirt_project.existing.id
  name       = "app-server-in-existing-project"
  cpu        = 4
  memory_mb  = 4096
  image      = "ubuntu:22.04"
}