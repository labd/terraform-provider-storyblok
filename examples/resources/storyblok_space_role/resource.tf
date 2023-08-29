resource "storyblok_space_role" "my_role" {
  space_id          = "<my-space-id>"
  role              = "Role"
  subtitle          = "A role description"
  permissions       = ["access_tasks"]
  field_permissions = ["component_name.field_name"]
  allowed_languages = ["default"]
  allowed_paths     = [1]
  external_id       = "1234"
}
