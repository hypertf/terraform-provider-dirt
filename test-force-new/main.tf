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

resource "dirt_instance" "test" {
  project_id = "b4ebaf26de699d4a24de8fbb7f255b0b"  # test-force-new-image project  
  name       = "final-test-instance"
  cpu        = 4
  memory_mb  = 2048
  image      = "ubuntu:20.04"  # Back to original image
  status     = "running"
}

output "instance_id" {
  value = dirt_instance.test.id
}

output "instance_image" {
  value = dirt_instance.test.image
}
