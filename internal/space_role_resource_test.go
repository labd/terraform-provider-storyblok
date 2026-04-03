package internal

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/labd/terraform-provider-storyblok/internal/utils"
)

func TestSpaceRoleResourceBasic(t *testing.T) {
	f, stop := ProviderFactories("./assets/space_role")
	defer func() {
		_ = stop()
	}()

	id := "test"
	rn := fmt.Sprintf("storyblok_space_role.%s", id)
	spaceId := 233252

	resource.Test(t, resource.TestCase{
		PreCheck:                 TestAccPreCheck(t),
		ProtoV6ProviderFactories: f,
		Steps: []resource.TestStep{
			{
				Config: testSpaceRoleConfig(id, spaceId),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(rn, "role", "tester"),
					resource.TestCheckResourceAttr(rn, "subtitle", "A test group"),
					resource.TestCheckResourceAttr(rn, "permissions.#", "1"),
					resource.TestCheckResourceAttr(rn, "permissions.0", "access_tasks"),
					resource.TestCheckResourceAttr(rn, "field_permissions.#", "1"),
					resource.TestCheckResourceAttr(rn, "field_permissions.0", "component_name.field_name"),
					resource.TestCheckResourceAttr(rn, "allowed_languages.#", "1"),
					resource.TestCheckResourceAttr(rn, "allowed_languages.0", "default"),
					resource.TestCheckResourceAttr(rn, "allowed_paths.#", "1"),
					resource.TestCheckResourceAttr(rn, "allowed_paths.0", "1"),
					resource.TestCheckResourceAttr(rn, "external_id", "FizBuz"),
				),
			},
			{
				Config: testSpaceRoleConfigUpdate(id, spaceId),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(rn, "role", "new-tester"),
					resource.TestCheckResourceAttr(rn, "subtitle", "A new test group"),
					resource.TestCheckResourceAttr(rn, "external_id", "BuzFiz"),
				),
			},
		},
	})
}

func testSpaceRoleConfig(identifier string, spaceId int) string {
	return utils.HCLTemplate(`
		resource "storyblok_space_role" "{{ .identifier }}" {
		  space_id          = "{{ .spaceId }}"
		  role              = "tester"
		  subtitle          = "A test group"
		  permissions       = ["access_tasks"]
		  field_permissions = ["component_name.field_name"]
		  allowed_languages = ["default"]
		  allowed_paths     = [1]
		  external_id       = "FizBuz"
		}
	`, map[string]any{
		"identifier": identifier,
		"spaceId":    spaceId,
	})
}

// TestSpaceRoleResourceNumericPaths verifies that the provider handles API
// responses where allowed_paths contains integers instead of strings.
func TestSpaceRoleResourceNumericPaths(t *testing.T) {
	f, stop := ProviderFactories("./assets/space_role_numeric_paths")
	defer func() {
		_ = stop()
	}()

	id := "numeric"
	rn := fmt.Sprintf("storyblok_space_role.%s", id)
	spaceId := 233252

	resource.Test(t, resource.TestCase{
		PreCheck:                 TestAccPreCheck(t),
		ProtoV6ProviderFactories: f,
		Steps: []resource.TestStep{
			{
				Config: testSpaceRoleConfigNumericPaths(id, spaceId),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(rn, "role", "bug-test-role"),
					resource.TestCheckResourceAttr(rn, "allowed_paths.#", "2"),
					resource.TestCheckResourceAttr(rn, "allowed_paths.0", "325468791"),
					resource.TestCheckResourceAttr(rn, "allowed_paths.1", "343319186"),
				),
			},
		},
	})
}

func testSpaceRoleConfigNumericPaths(identifier string, spaceId int) string {
	return utils.HCLTemplate(`
		resource "storyblok_space_role" "{{ .identifier }}" {
		  space_id          = "{{ .spaceId }}"
		  role              = "bug-test-role"
		  subtitle          = ""
		  permissions       = ["read_stories","save_stories"]
		  field_permissions = []
		  allowed_languages = ["en-au","en-nz"]
		  allowed_paths     = ["325468791","343319186"]
		}
	`, map[string]any{
		"identifier": identifier,
		"spaceId":    spaceId,
	})
}

func testSpaceRoleConfigUpdate(identifier string, spaceId int) string {
	return utils.HCLTemplate(`
		resource "storyblok_space_role" "{{ .identifier }}" {
		  space_id          = "{{ .spaceId }}"
		  role              = "new-tester"
		  subtitle          = "A new test group"
		  permissions       = ["access_tasks"]
		  field_permissions = ["component_name.field_name"]
		  allowed_languages = ["default"]
		  allowed_paths     = [1]
		  external_id       = "BuzFiz"
		}
	`, map[string]any{
		"identifier": identifier,
		"spaceId":    spaceId,
	})
}
