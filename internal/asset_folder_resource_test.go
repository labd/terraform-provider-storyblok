package internal

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/labd/terraform-provider-storyblok/internal/utils"
)

func TestAssetFolderResourceBasic(t *testing.T) {
	f, stop := ProviderFactories("./assets/asset_folder")
	defer func() {
		_ = stop()
	}()

	id := "test"
	rn := fmt.Sprintf("storyblok_asset_folder.%s", id)
	spaceId := 233252

	resource.Test(t, resource.TestCase{
		PreCheck:                 TestAccPreCheck(t),
		ProtoV6ProviderFactories: f,
		Steps: []resource.TestStep{
			{
				Config: testAssetFolderConfig(id, spaceId),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(rn, "name", "asset-folder-name"),
				),
			},
			{
				Config: testAssetFolderConfigUpdate(id, spaceId),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(rn, "name", "new-asset-folder-name"),
				),
			},
		},
	})
}

func testAssetFolderConfig(identifier string, spaceId int) string {
	return utils.HCLTemplate(`
		resource "storyblok_asset_folder" "{{ .identifier }}" {
		  space_id = {{ .spaceId }}
		  name = "asset-folder-name"
		}
	`, map[string]any{
		"identifier": identifier,
		"spaceId":    spaceId,
	})
}

func testAssetFolderConfigUpdate(identifier string, spaceId int) string {
	return utils.HCLTemplate(`
		resource "storyblok_asset_folder" "{{ .identifier }}" {
		  space_id = {{ .spaceId }}
		  name = "new-asset-folder-name"
		}
	`, map[string]any{
		"identifier": identifier,
		"spaceId":    spaceId,
	})
}
