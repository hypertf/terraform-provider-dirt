# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

terraform {
  required_providers {
    dirt = {
      source  = "hypertf/dirt"
      version = "0.1.0"
    }
  }
}

provider "dirt" {
  endpoint = "http://localhost:8080/v1"
  token    = "test-token"
}

# Create a project
resource "dirt_project" "example" {
  name = "test-project-from-published-provider"
}

output "project_id" {
  value = dirt_project.example.id
}
