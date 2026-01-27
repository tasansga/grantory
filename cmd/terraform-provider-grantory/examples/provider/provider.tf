terraform {
  required_providers {
    grantory = {
      source = "tasansga/grantory"
    }
  }
}

provider "grantory" {
  server = "http://localhost:8080"
}
