terraform {
  required_providers {
    mock = {
      source = "liamcervante/mock"
    }
  }
}

provider "mock" {}

resource "mock_simple_resource" "simple_resource" {
  id = "192977d6-b169-4170-a9d4-ee1dcef7c6ea"
  integer = 1
}
