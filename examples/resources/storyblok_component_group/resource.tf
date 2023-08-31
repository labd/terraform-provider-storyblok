resource "storyblok_component_group" "my_component_group" {
  space_id = "<my-space-id>"
  name     = "my-component-group"
}

resource "storyblok_component" "component1" {
  name     = "component1"
  space_id = storyblok_component_group.my_component_group.space_id
  component_group_uuid = storyblok_component_group.my_component_group.uuid

  schema = {
    field1 = {
      type     = "text"
      position = 1
    }
  }
}

resource "storyblok_component" "component2" {
  name     = "component2"
  space_id = storyblok_component_group.my_component_group.space_id
  component_group_uuid = storyblok_component_group.my_component_group.uuid

  schema = {
    field1 = {
      type     = "image"
      position = 1
    }
  }
}
