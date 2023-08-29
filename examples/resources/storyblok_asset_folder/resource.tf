resource "storyblok_asset_folder" "parent" {
  name = "parent"
}

resource "storyblok_asset_folder" "child" {
  name      = "child"
  parent_id = storyblok_asset_folder.parent.asset_folder_id
}
