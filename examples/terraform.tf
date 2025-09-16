terraform {
  required_providers {
    dirt = {
      version = "0.1.0"
    }
  }
}

provider "dirt" {
  # DirtCloud API endpoint (optional, defaults to http://localhost:8080/v1)
  # Can also be set via DIRT_ENDPOINT environment variable
  endpoint = "http://localhost:8080/v1"

  # DirtCloud API token for authentication (optional)
  # Can also be set via DIRT_TOKEN environment variable
  # token = "your-api-token-here"
}
