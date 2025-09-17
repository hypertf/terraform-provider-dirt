terraform {
  required_providers {
    dirt = {
      source = "hashicorp/dirt"
    }
  }
}

provider "dirt" {
  endpoint = "http://localhost:8080/v1"
}

# Import block - import existing instance
import {
  to = dirt_instance.imported
  id = "4b5515de9539041aba87e52c527bedeb"
}

# Imported instance - will show as import in plan
resource "dirt_instance" "imported" {
  project_id = "b4ebaf26de699d4a24de8fbb7f255b0b"
  name       = "instance-to-import"
  cpu        = 1
  memory_mb  = 1024
  image      = "alpine:latest"
  status     = "running"
}

# Import the existing test instance so we can replace it
import {
  to = dirt_instance.replaced
  id = "bf5c15b1ba6af7d78c34531a23467bc3"
}

# Existing instance from previous test - will be replaced due to image change
resource "dirt_instance" "replaced" {
  project_id = "b4ebaf26de699d4a24de8fbb7f255b0b"
  name       = "test-instance"
  cpu        = 2
  memory_mb  = 2048
  image      = "debian:11"  # Changed from ubuntu:22.04 - Forces replacement
  status     = "running"
}

# Instance to be updated in-place (change memory)
resource "dirt_instance" "updated" {
  project_id = "b4ebaf26de699d4a24de8fbb7f255b0b"
  name       = "updated-instance"
  cpu        = 4  # Already updated
  memory_mb  = 4096  # Changed from 2048 - will update in-place
  image      = "ubuntu:20.04"
  status     = "running"
}

# New instance - this was moved to renamed_instance below

# Instance was removed from config - will be destroyed

# Another new instance - will be created
resource "dirt_instance" "another_new" {
  project_id = "b4ebaf26de699d4a24de8fbb7f255b0b"
  name       = "another-new-instance"
  cpu        = 2
  memory_mb  = 1024
  image      = "postgres:14"
  status     = "running"
}

# Move the new instance to a different name for demonstration
moved {
  from = dirt_instance.new
  to   = dirt_instance.renamed_instance
}

# The moved resource with new name
resource "dirt_instance" "renamed_instance" {
  project_id = "b4ebaf26de699d4a24de8fbb7f255b0b"
  name       = "renamed-instance"  # Changed name
  cpu        = 1
  memory_mb  = 512
  image      = "nginx:latest"
  status     = "running"
}

# Trigger replacement by changing image
resource "dirt_instance" "force_replace" {
  project_id = "b4ebaf26de699d4a24de8fbb7f255b0b"
  name       = "replacement-demo"
  cpu        = 1
  memory_mb  = 1024
  image      = "ubuntu:22.04"  # Changed from ubuntu:18.04 - forces replacement
  status     = "running"
}
