resource "storyblok_space_role" "test_role" {
  space_id          = 245461
  role              = "tester"
  subtitle          = "A test group"
  permissions       = ["access_tasks"]
  field_permissions = ["component_name.field_name"]
  allowed_languages = ["default"]
  allowed_paths     = [1]
  external_id       = "FizBuz"
}
