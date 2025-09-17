# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

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

# IMPORT OPERATIONS - Import existing standalone instances
import {
  to = dirt_instance.imported_mysql
  id = "402e1dd567fd1a0eed7c96ebcdbc6f5c"
}

resource "dirt_instance" "imported_mysql" {
  project_id = "b4ebaf26de699d4a24de8fbb7f255b0b"
  name       = "standalone-instance-1"
  cpu        = 1
  memory_mb  = 512
  image      = "mysql:8.0"
  status     = "running"
}

import {
  to = dirt_instance.imported_mongo
  id = "b0946fae3349a3b9739c361fc3f26694"
}

resource "dirt_instance" "imported_mongo" {
  project_id = "b4ebaf26de699d4a24de8fbb7f255b0b"
  name       = "standalone-instance-2"
  cpu        = 2
  memory_mb  = 1024
  image      = "mongodb:6.0"
  status     = "running"
}

# UPDATE OPERATIONS - Modify existing resources
resource "dirt_instance" "existing_1" {
  project_id = "b4ebaf26de699d4a24de8fbb7f255b0b"
  name       = "updated-instance-name"  # Changed name - update in place
  cpu        = 2  # Changed CPU - update in place  
  memory_mb  = 1024
  image      = "alpine:3.18"
  status     = "running"
}

# REPLACE OPERATIONS - Change image to force replacement
resource "dirt_instance" "existing_2" {
  project_id = "b4ebaf26de699d4a24de8fbb7f255b0b"
  name       = "simple-instance-2"
  cpu        = 2
  memory_mb  = 2048
  image      = "nginx:1.25"  # Changed image version - forces replacement!
  status     = "running"
}

# MOVE OPERATIONS - Move resource to new name
moved {
  from = dirt_instance.to_be_moved
  to   = dirt_instance.moved_redis
}

resource "dirt_instance" "moved_redis" {
  project_id = "b4ebaf26de699d4a24de8fbb7f255b0b"
  name       = "redis-moved-instance"  # Updated name after move
  cpu        = 1
  memory_mb  = 512
  image      = "redis:7"
  status     = "running"
}

# CREATE OPERATIONS - Brand new resources
resource "dirt_instance" "brand_new" {
  project_id = "b4ebaf26de699d4a24de8fbb7f255b0b"
  name       = "newly-created-instance"
  cpu        = 1
  memory_mb  = 1024
  image      = "postgres:15"
  status     = "running"
}

# DESTROY OPERATIONS - Removed from config to demonstrate destroy
# resource "dirt_instance" "will_be_destroyed" was removed
