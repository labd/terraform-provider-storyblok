package internal

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/labd/terraform-provider-storyblok/internal/utils"
)

func TestComponentGroupResourceBasic(t *testing.T) {
	f, stop := ProviderFactories("./assets/component_group")
	defer func() {
		_ = stop()
	}()

	id := "test"
	rn := fmt.Sprintf("storyblok_component_group.%s", id)
	spaceId := 233252

	resource.Test(t, resource.TestCase{
		PreCheck:                 TestAccPreCheck(t),
		ProtoV6ProviderFactories: f,
		Steps: []resource.TestStep{
			{
				Config: testComponentGroupConfig(id, spaceId),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(rn, "name", "test-component-group"),
				),
			},
			{
				Config: testComponentGroupConfigUpdate(id, spaceId),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(rn, "name", "new-test-component-group"),
				),
			},
		},
	})
}

func testComponentGroupConfig(identifier string, spaceId int) string {
	return utils.HCLTemplate(`
		resource "storyblok_component_group" "{{ .identifier }}" {
		  space_id = "{{ .spaceId }}"
		  name     = "test-component-group"
		}
	`, map[string]any{
		"identifier": identifier,
		"spaceId":    spaceId,
	})
}

func testComponentGroupConfigUpdate(identifier string, spaceId int) string {
	return utils.HCLTemplate(`
		resource "storyblok_component_group" "{{ .identifier }}" {
		  space_id = "{{ .spaceId }}"
		  name     = "new-test-component-group"
		}
	`, map[string]any{
		"identifier": identifier,
		"spaceId":    spaceId,
	})
}
