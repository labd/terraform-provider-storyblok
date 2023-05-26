package internal

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/labd/storyblok-go-sdk/sbmgmt"
)

// componentResourceModel maps the resource schema data.
type componentResourceModel struct {
	ID          types.String          `tfsdk:"id"`
	ComponentID types.Int64           `tfsdk:"component_id"`
	SpaceID     types.Int64           `tfsdk:"space_id"`
	CreatedAt   types.String          `tfsdk:"created_at"`
	Name        types.String          `tfsdk:"name"`
	IsRoot      types.Bool            `tfsdk:"is_root"`
	IsNestable  types.Bool            `tfsdk:"is_nestable"`
	Schema      map[string]fieldModel `tfsdk:"schema"`
}

type fieldModel struct {
	Type     types.String `tfsdk:"type"`
	Position types.Int64  `tfsdk:"position"`
}

func (m *componentResourceModel) toRemoteInput() sbmgmt.ComponentInput {

	raw := make(map[string]sbmgmt.FieldInput, len(m.Schema))
	for name, item := range m.Schema {
		raw[name] = sbmgmt.FieldInput{
			Type: item.Type.ValueString(),
			Pos:  item.Position.ValueInt64(),
		}
	}

	// Sort the fields by position. Storyblok has a position field but ends up
	// using the ordering of the json...
	schema := sortComponentFields(raw)

	return sbmgmt.ComponentInput{
		Name:       m.Name.ValueString(),
		Schema:     schema,
		IsNestable: m.IsNestable.ValueBoolPointer(),
		IsRoot:     m.IsRoot.ValueBoolPointer(),
	}
}

func (m *componentResourceModel) fromRemote(spaceID int64, c *sbmgmt.Component) error {
	if c == nil {
		return fmt.Errorf("component is nil")
	}
	m.ID = types.StringValue(createIdentifier(spaceID, c.Id))
	m.ComponentID = types.Int64Value(c.Id)
	m.CreatedAt = types.StringValue(c.CreatedAt.String())
	m.IsRoot = types.BoolPointerValue(c.IsRoot)
	m.IsNestable = types.BoolPointerValue(c.IsNestable)

	schema := make(map[string]fieldModel, c.Schema.Len())
	for pair := c.Schema.Oldest(); pair != nil; pair = pair.Next() {
		name := pair.Key
		field := pair.Value

		schema[name] = fieldModel{
			Type:     types.StringValue(field.Type),
			Position: types.Int64Value(field.Pos),
		}
	}
	m.Schema = schema
	return nil
}

func getComponentTypes() map[string]string {
	return map[string]string{
		"bloks":      "Blocks: a field to interleave other components in your current one",
		"text":       "Text: a text field",
		"textarea":   "Textarea: a text area",
		"markdown":   "Markdown: write markdown with a text area and additional formatting options",
		"number":     "Number: a number field",
		"datetime":   "Date/Time: a date- and time picker",
		"boolean":    "Boolean: a checkbox - true/false",
		"options":    "Multi-Options: a list of checkboxes",
		"option":     "Single-Option: a single dropdown",
		"asset":      "Asset: Single asset (images, videos, audio, and documents)",
		"multiasset": "Multi-Assets: (images, videos, audio, and documents)",
		"multilink":  "Link: an input field for internal linking to other stories",
		"section":    "Group: no input possibility - allows you to group fields in sections",
		"custom":     "Plugin: Extend the editor yourself with a color picker or similar - Check out: Creating a Storyblok field type plugin",
		"image":      "Image (old): a upload field for a single image with cropping possibilities",
		"file":       "File (old): a upload field for a single file",
	}
}
