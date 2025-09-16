data "dirt_instance" "existing" {
  id = "instance-id-67890"
}

output "instance_details" {
  description = "Details of the existing instance"
  value = {
    name       = data.dirt_instance.existing.name
    project_id = data.dirt_instance.existing.project_id
    cpu        = data.dirt_instance.existing.cpu
    memory_mb  = data.dirt_instance.existing.memory_mb
    image      = data.dirt_instance.existing.image
    status     = data.dirt_instance.existing.status
  }
}

output "instance_created_at" {
  description = "When the instance was created"
  value       = data.dirt_instance.existing.created_at
}

output "instance_updated_at" {
  description = "When the instance was last updated"
  value       = data.dirt_instance.existing.updated_at
}

# Example using the data source to create metadata
resource "dirt_metadata" "instance_backup_schedule" {
  path  = "instances/${data.dirt_instance.existing.id}/backup_schedule"
  value = "daily-midnight"
}

resource "dirt_metadata" "instance_monitoring" {
  path  = "instances/${data.dirt_instance.existing.id}/monitoring_enabled"
  value = "true"
}