---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "storyblok_component_group Resource - storyblok"
subcategory: ""
description: |-
  Manage a component.
---

# storyblok_component_group (Resource)

Manage a component.

## Example Usage

```terraform
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
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) The name of the component group.
- `space_id` (Number) The ID of the space.

### Read-Only

- `created_at` (String) The creation timestamp of the component group.
- `group_id` (Number) The ID of the component group.
- `id` (String) The terraform ID of the component.
- `updated_at` (String) The creation timestamp of the component group.
- `uuid` (String) The UUID of the component group.


