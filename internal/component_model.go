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
	Schema      map[string]fieldModel `tfsdk:"schema"`
}

type fieldModel struct {
	Type     types.String `tfsdk:"type"`
	Position types.Int64  `tfsdk:"position"`
}

func (m *componentResourceModel) toRemoteInput() sbmgmt.ComponentInput {
	schema := make(map[string]sbmgmt.FieldInput, len(m.Schema))
	for name, item := range m.Schema {
		schema[name] = sbmgmt.FieldInput{
			Type: item.Type.ValueString(),
			Pos:  item.Position.ValueInt64(),
		}
	}

	return sbmgmt.ComponentInput{
		Name:   m.Name.ValueString(),
		Schema: ref(schema),
	}
}

func (m *componentResourceModel) fromRemote(spaceID int64, c *sbmgmt.Component) error {
	if c == nil {
		return fmt.Errorf("component is nil")
	}
	m.ID = types.StringValue(createIdentifier(spaceID, c.Id))
	m.ComponentID = types.Int64Value(int64(c.Id))
	m.CreatedAt = types.StringValue(c.CreatedAt.String())

	schema := make(map[string]fieldModel, len(*c.Schema))
	for name, field := range *c.Schema {
		schema[name] = fieldModel{
			Type:     types.StringValue(field.Type),
			Position: types.Int64Value(field.Pos),
		}
	}
	m.Schema = schema
	return nil
}
