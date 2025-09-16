# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

terraform {
  required_providers {
    dirt = {
      version = "0.1.0"
    }
  }
}

provider "dirt" {
  endpoint = "http://localhost:8080/v1"
}

# First create a metadata resource
resource "dirt_metadata" "test" {
  path  = "test/path/example"
  value = "test-value-123"
}

# Then try to read it back with the data source using path
data "dirt_metadata" "test_read" {
  path = dirt_metadata.test.path
}

output "created_metadata" {
  value = {
    id    = dirt_metadata.test.id
    path  = dirt_metadata.test.path
    value = dirt_metadata.test.value
  }
}

output "read_metadata" {
  value = {
    id    = data.dirt_metadata.test_read.id
    path  = data.dirt_metadata.test_read.path
    value = data.dirt_metadata.test_read.value
  }
}
