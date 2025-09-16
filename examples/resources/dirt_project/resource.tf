resource "dirt_project" "example" {
  name = "my-dirt-project"
}

# Example of creating multiple projects
resource "dirt_project" "web_project" {
  name = "web-servers"
}

resource "dirt_project" "db_project" {
  name = "database-cluster"
}