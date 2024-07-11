// simple example
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


// advanced example
resource "storyblok_component" "advanced_component" {
  name        = "advanced-component"
  space_id    = "<my-space-id>"
  is_root     = true
  is_nestable = false

  schema = {
    title = {
      type        = "text"
      position    = 1
      required    = true                     // The field is required. Default is false.
      max_length  = 200                      // Set the max length of the input string
      description = "Title of the component" // Description shown in the editor interface
    }

    introduction = {
      type          = "rich_text"
      position      = 2
      rich_markdown = true // Enable rich markdown view by default
      description   = "Introduction text with rich text editor"
    }

    image = {
      type            = "image"
      position        = 3
      asset_folder_id = 1    // Default asset folder numeric id to store uploaded image of that field
      add_https       = true // Prepends https: to stop usage of relative protocol
      image_crop      = true // Activate force crop for images
    }

    release_date = {
      type         = "date"
      position     = 4
      disable_time = true // Disables time selection from date picker
      description  = "Release date of the content"
    }

    tags = {
      type            = "multi_option"
      position        = 5
      datasource_slug = "tags" // Define selectable datasources string
      description     = "Tags for the component"
    }

    rating = {
      type          = "number"
      position      = 6
      description   = "Rating of the content"
      default_value = "3" // Default value for the field
    }

    content = {
      type                = "bloks"
      position            = 7
      component_whitelist = ["text", "image", "video"] // Array of component/content type names
      maximum             = 10                         // Maximum amount of added bloks in this blok field
      description         = "Content blocks"
    }
  }
}

// conditional content
resource "storyblok_component" "conditional_settings_new" {
  name        = "conditional settings component"
  space_id    = "<your space id>"
  is_root     = false
  is_nestable = true


  schema = {
    content = {
      position     = 0
      translatable = true
      display_name = "Content"
      required     = true
      type         = "text"
    }

    more_content = {
      position     = 1
      translatable = true
      display_name = "more content"
      required     = true
      type         = "text"
    }

    conditionalContent = {
      position     = 2
      display_name = "conditinal content"
      required     = true
      type         = "text"

      conditional_settings = [
        {
          modifications = [
            {
              required = false
            }
          ]

          // make "conditional content" optional of either:
          // 1. content is empty
          // 2. more content equals "test"
          rule_match = "any"
          rule_conditions = [
            {
              validation = "empty"
              validated_object = {
                field_key = "content"
              }
            },
            {
              value      = "test"
              validation = "equals"
              validated_object = {
                field_key = "more_content"
              }
            }
          ]
        }
      ]
    }
  }
}
