terraform {
  required_providers {
    dirt = {
      source  = "hypertf/dirt"
      version = ">= 0.2.0"
    }
  }
}

provider "dirt" {
  endpoint = "http://localhost:8080/v1"
}
