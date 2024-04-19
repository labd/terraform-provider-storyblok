package internal

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestComponentResourceBasic(t *testing.T) {
	f, stop := ProviderFactories("./assets/component")
	defer stop()

	id := "test"
	rn := fmt.Sprintf("storyblok_component.%s", id)
	spaceId := 233252

	resource.Test(t, resource.TestCase{
		PreCheck:                 TestAccPreCheck(t),
		ProtoV6ProviderFactories: f,
		Steps: []resource.TestStep{
			{
				Config: testComponentConfig(id, spaceId),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(rn, "name", "test-banner"),
					resource.TestCheckResourceAttr(rn, "schema.title.position", "1"),
					resource.TestCheckResourceAttr(rn, "schema.title.type", "text"),
					resource.TestCheckResourceAttr(rn, "schema.intro.position", "2"),
					resource.TestCheckResourceAttr(rn, "schema.intro.type", "text"),
					resource.TestCheckResourceAttr(rn, "schema.image.position", "3"),
					resource.TestCheckResourceAttr(rn, "schema.image.type", "image"),
				),
			},
			{
				Config: testComponentConfigUpdate(id, spaceId),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(rn, "name", "new-test-banner"),
					resource.TestCheckResourceAttr(rn, "schema.intro.position", "1"),
					resource.TestCheckResourceAttr(rn, "schema.intro.type", "text"),
					resource.TestCheckResourceAttr(rn, "schema.title.position", "2"),
					resource.TestCheckResourceAttr(rn, "schema.title.type", "text"),
				),
			},
		},
	})
}

func testComponentConfig(identifier string, spaceId int) string {
	return HCLTemplate(`
		resource "storyblok_component" "{{ .identifier }}" {
		  name     = "test-banner"
		  space_id = "{{ .spaceId }}"
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
	      preview_tmpl = "<div></div>"
		}
	`, map[string]any{
		"identifier": identifier,
		"spaceId":    spaceId,
	})
}

func testComponentConfigUpdate(identifier string, spaceId int) string {
	return HCLTemplate(`
		resource "storyblok_component" "{{ .identifier }}" {
		  name     = "new-test-banner"
		  space_id = "{{ .spaceId }}"
		  schema = {
			title = {
			  type     = "text"
			  position = 2
			}
		
			intro = {
			  type     = "text"
			  position = 1
			}
		  }
		}
	`, map[string]any{
		"identifier": identifier,
		"spaceId":    spaceId,
	})
}
