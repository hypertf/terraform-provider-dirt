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
