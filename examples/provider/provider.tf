terraform {
  required_providers {
    storyblok = {
      source  = "labd/storyblok"
      version = "1.0.2"
    }
  }
}


provider "storyblok" {
  url   = "https://mapi.storyblok.com"
  token = "<my-token>"
}


resource "storyblok_component" "banner" {
  name     = "my-banner"
  space_id = "<my-space-id>"
  schema = {

    title = {
      type     = "text"
      position = 1
    }

    intro = {
      type     = "text"
      position = 2
    }

    image = {
      type     = "image"
      position = 3
    }
  }
}
