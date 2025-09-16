# Create projects for different environments
resource "dirt_project" "web_project" {
  name = "web-application"
}

resource "dirt_project" "database_project" {
  name = "database-cluster"
}

# Create web server instances
resource "dirt_instance" "web_server_1" {
  project_id = dirt_project.web_project.id
  name       = "web-server-1"
  cpu        = 4
  memory_mb  = 4096
  image      = "nginx:latest"
  status     = "running"
}

resource "dirt_instance" "web_server_2" {
  project_id = dirt_project.web_project.id
  name       = "web-server-2"
  cpu        = 4
  memory_mb  = 4096
  image      = "nginx:latest"
  status     = "running"
}

# Create load balancer instance
resource "dirt_instance" "load_balancer" {
  project_id = dirt_project.web_project.id
  name       = "load-balancer"
  cpu        = 2
  memory_mb  = 2048
  image      = "haproxy:latest"
  status     = "running"
}

# Create database instances
resource "dirt_instance" "database_primary" {
  project_id = dirt_project.database_project.id
  name       = "postgres-primary"
  cpu        = 8
  memory_mb  = 16384
  image      = "postgres:15"
  status     = "running"
}

resource "dirt_instance" "database_replica" {
  project_id = dirt_project.database_project.id
  name       = "postgres-replica"
  cpu        = 8
  memory_mb  = 16384
  image      = "postgres:15"
  status     = "running"
}

# Store configuration in metadata
resource "dirt_metadata" "database_connection_string" {
  path  = "app/config/database_url"
  value = "postgres://app:password@${dirt_instance.database_primary.id}:5432/appdb"
}

resource "dirt_metadata" "database_replica_url" {
  path  = "app/config/database_replica_url"
  value = "postgres://app:password@${dirt_instance.database_replica.id}:5432/appdb"
}

resource "dirt_metadata" "load_balancer_config" {
  path  = "loadbalancer/backends"
  value = "${dirt_instance.web_server_1.id}:80,${dirt_instance.web_server_2.id}:80"
}

# Application configuration
resource "dirt_metadata" "app_version" {
  path  = "app/version"
  value = "1.0.0"
}

resource "dirt_metadata" "feature_flags" {
  path  = "features/redis_cache"
  value = "enabled"
}

resource "dirt_metadata" "monitoring_config" {
  path = "monitoring/instances"
  value = jsonencode({
    web_servers = [
      dirt_instance.web_server_1.id,
      dirt_instance.web_server_2.id
    ],
    load_balancer = dirt_instance.load_balancer.id,
    databases = [
      dirt_instance.database_primary.id,
      dirt_instance.database_replica.id
    ]
  })
}

# Instance-specific metadata
resource "dirt_metadata" "web_server_1_role" {
  path  = "instances/${dirt_instance.web_server_1.id}/role"
  value = "web-server"
}

resource "dirt_metadata" "web_server_2_role" {
  path  = "instances/${dirt_instance.web_server_2.id}/role"
  value = "web-server"
}

resource "dirt_metadata" "database_primary_role" {
  path  = "instances/${dirt_instance.database_primary.id}/role"
  value = "database-primary"
}

resource "dirt_metadata" "database_replica_role" {
  path  = "instances/${dirt_instance.database_replica.id}/role"
  value = "database-replica"
}

# Read existing configuration (example of data sources)
data "dirt_metadata" "external_api_key" {
  path = "external/api_key"
}

# Use data source values in new metadata
resource "dirt_metadata" "app_external_config" {
  path  = "app/config/external_api_config"
  value = "api_key=${data.dirt_metadata.external_api_key.value}&timeout=30"
}

# Outputs
output "web_project_id" {
  description = "ID of the web project"
  value       = dirt_project.web_project.id
}

output "database_project_id" {
  description = "ID of the database project"
  value       = dirt_project.database_project.id
}

output "web_server_instances" {
  description = "Web server instance IDs"
  value = [
    dirt_instance.web_server_1.id,
    dirt_instance.web_server_2.id
  ]
}

output "database_instances" {
  description = "Database instance details"
  value = {
    primary = {
      id        = dirt_instance.database_primary.id
      name      = dirt_instance.database_primary.name
      cpu       = dirt_instance.database_primary.cpu
      memory_mb = dirt_instance.database_primary.memory_mb
    }
    replica = {
      id        = dirt_instance.database_replica.id
      name      = dirt_instance.database_replica.name
      cpu       = dirt_instance.database_replica.cpu
      memory_mb = dirt_instance.database_replica.memory_mb
    }
  }
}

output "load_balancer_endpoint" {
  description = "Load balancer instance ID"
  value       = dirt_instance.load_balancer.id
}

output "metadata_ids" {
  description = "IDs of created metadata entries"
  value = {
    database_connection = dirt_metadata.database_connection_string.id
    app_version         = dirt_metadata.app_version.id
    monitoring_config   = dirt_metadata.monitoring_config.id
    feature_flags       = dirt_metadata.feature_flags.id
  }
}
