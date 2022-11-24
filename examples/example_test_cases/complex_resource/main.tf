terraform {
  required_providers {
    tfcoremock = {
      source = "hashicorp/tfcoremock"
    }
  }
}

provider "tfcoremock" {}

data "tfcoremock_simple_resource" "simple_resource" {
  id = "192977d6-b169-4170-a9d4-ee1dcef7c6ea"
}

resource "tfcoremock_complex_resource" "complex_resource" {
  id = "d199d8ea-e8f8-4fb0-8276-3567a74d3db8"

  integer = data.tfcoremock_simple_resource.simple_resource.integer

  object = {
    string = "hello"
    bool = true
  }

  list = [
    {
      string = "one.one"
    },
    {
      string = "one.two"
    }
  ]

  list_block {
    string = "two.one"
  }

  list_block {
    string = "two.two"
  }

}
