resource "dirt_instance" "web_server" {
  project_id = dirt_project.example.id
  name       = "web-server-1"
  cpu        = 4
  memory_mb  = 4096
  image      = "ubuntu:22.04"
  status     = "running"
}

resource "dirt_instance" "database" {
  project_id = dirt_project.example.id
  name       = "postgres-db"
  cpu        = 8
  memory_mb  = 8192
  image      = "postgres:15"
  status     = "running"
}

# Example with minimal configuration (uses defaults)
resource "dirt_instance" "minimal" {
  project_id = dirt_project.example.id
  name       = "minimal-instance"
  # cpu defaults to 2
  # memory_mb defaults to 2048
  # image defaults to ubuntu:20.04
  # status defaults to running
}

# Reference the project from the project example
resource "dirt_project" "example" {
  name = "instance-example-project"
}