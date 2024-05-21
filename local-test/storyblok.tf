terraform {
  required_providers {
    storyblok = {
      source  = "labd/storyblok"
      version = "99.0.0"
    }
  }
}

provider "storyblok" {
  url   = "https://mapi.storyblok.com"
  token = "dmWs0jvlBUf3QBIJ8tAHrQtt-199599-q6zdnS7yUHrYYzEaiH8H"
}



resource "storyblok_component" "test" {
  name         = "test"
  space_id     = 233774
  is_root      = false
  is_nestable  = true
  display_name = "test Name"
  schema = {
    title = {
      position     = 0
      translatable = true
      display_name = "Titlee"
      required     = false
      type         = "text"
    }
  }
}
