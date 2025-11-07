package internal

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/labd/terraform-provider-storyblok/internal/utils"
)

func TestComponentResourceBasic(t *testing.T) {
	f, stop := ProviderFactories("./assets/component")
	defer func() {
		_ = stop()
	}()

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
					resource.TestCheckResourceAttr(rn, "schema.image.conditional_settings.0.modifications.0.required", "true"),
					resource.TestCheckResourceAttr(rn, "schema.image.conditional_settings.0.rule_match", "all"),
					resource.TestCheckResourceAttr(rn, "schema.image.conditional_settings.0.rule_conditions.0.validation", "empty"),
					resource.TestCheckResourceAttr(rn, "schema.image.conditional_settings.0.rule_conditions.0.validated_object.field_key", "intro"),
					resource.TestCheckResourceAttr(rn, "schema.link.type", "multilink"),
					resource.TestCheckResourceAttr(rn, "schema.link.description", "Link to a page"),
					resource.TestCheckResourceAttr(rn, "schema.link.translatable", "true"),
					resource.TestCheckResourceAttr(rn, "schema.link.required", "true"),
					resource.TestCheckResourceAttr(rn, "schema.link.allow_external_url", "true"),
					resource.TestCheckResourceAttr(rn, "schema.link.restrict_content_types", "true"),
					resource.TestCheckResourceAttr(rn, "schema.link.allow_advanced_search", "true"),
					resource.TestCheckResourceAttr(rn, "schema.link.allow_custom_attributes", "true"),
					resource.TestCheckResourceAttr(rn, "schema.link.asset_link_type", "true"),
					resource.TestCheckResourceAttr(rn, "schema.link.show_anchor", "true"),
					resource.TestCheckResourceAttr(rn, "schema.link.email_link_type", "true"),
					resource.TestCheckResourceAttr(rn, "schema.link.position", "3"),
					resource.TestCheckResourceAttr(rn, "schema.link.allow_target_blank", "true"),
					resource.TestCheckResourceAttr(rn, "schema.link.link_scope", "{0}"),
					resource.TestCheckResourceAttr(rn, "schema.link.force_link_scope", "true"),
					resource.TestCheckResourceAttr(rn, "schema.link.tooltip", "true"),
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
					resource.TestCheckResourceAttr(rn, "schema.buttons.filter_content_type.0", "button"),
					resource.TestCheckResourceAttr(rn, "schema.buttons.type", "options"),
					resource.TestCheckResourceAttr(rn, "schema.buttons.position", "3"),
					resource.TestCheckResourceAttr(rn, "schema.link.description", "Other link"),
					resource.TestCheckResourceAttr(rn, "schema.link.translatable", "false"),
					resource.TestCheckResourceAttr(rn, "schema.link.required", "false"),
					resource.TestCheckResourceAttr(rn, "schema.link.allow_external_url", "false"),
					resource.TestCheckResourceAttr(rn, "schema.link.restrict_content_types", "false"),
					resource.TestCheckResourceAttr(rn, "schema.link.allow_advanced_search", "false"),
					resource.TestCheckResourceAttr(rn, "schema.link.allow_custom_attributes", "false"),
					resource.TestCheckResourceAttr(rn, "schema.link.asset_link_type", "false"),
					resource.TestCheckResourceAttr(rn, "schema.link.show_anchor", "false"),
					resource.TestCheckResourceAttr(rn, "schema.link.email_link_type", "false"),
					resource.TestCheckResourceAttr(rn, "schema.link.position", "1"),
					resource.TestCheckResourceAttr(rn, "schema.link.allow_target_blank", "false"),
					resource.TestCheckResourceAttr(rn, "schema.link.force_link_scope", "false"),
					resource.TestCheckResourceAttr(rn, "schema.link.tooltip", "false"),
				),
			},
		},
	})
}

func testComponentConfig(identifier string, spaceId int) string {
	return utils.HCLTemplate(`
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
	
				link = {
				  type                    = "multilink"
				  description 			  = "Link to a page"
				  translatable 			  = true
				  required 				  = true
				  allow_external_url 	  = true
				  restrict_content_types  = true
				  allow_advanced_search   = true
				  allow_custom_attributes = true
				  asset_link_type         = true
				  show_anchor             = true
				  email_link_type         = true
				  position                = 3
				  allow_target_blank      = true
				  link_scope              = "{0}"
				  force_link_scope        = true
				  tooltip 				  = true
				}
			
				image = {
					type     = "image"
					position = 3

					conditional_settings = [
						{
							modifications = [
								{
									required = true
								}
							]

							rule_match = "all"
							rule_conditions = [
								{
									validation = "empty"
									validated_object = {
										field_key = "intro"
									}
								}
							]
						}
					]
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
	return utils.HCLTemplate(`
		resource "storyblok_component" "{{ .identifier }}" {
		  name     = "new-test-banner"
		  space_id = "{{ .spaceId }}"
		  schema = {
			title = {
			  type     = "text"
			  position = 2
			}
	
			link = {
			  type                    = "multilink"
			  description 			  = "Other link"
			  translatable 			  = false
			  required 				  = false
			  allow_external_url 	  = false
			  restrict_content_types  = false
			  allow_advanced_search   = false
			  allow_custom_attributes = false
			  asset_link_type         = false
			  show_anchor             = false
			  email_link_type         = false
			  position                = 1
			  allow_target_blank      = false
			  link_scope              = "{0}"
			  force_link_scope        = false
			  tooltip 				  = false
			}
		
			intro = {
			  type     = "text"
			  position = 1
			}

			buttons = {
				type = "options"
				source = "internal_stories"
				position = 3
				filter_content_type = ["button"]
			}
		  }
			preview_tmpl = "<div></div>"
		}
	`, map[string]any{
		"identifier": identifier,
		"spaceId":    spaceId,
	})
}
