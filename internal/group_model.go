package internal

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/labd/storyblok-go-sdk/sbmgmt"
)

// componentGroupResourceModel maps the resource schema data.
type componentGroupResourceModel struct {
	ID        types.String `tfsdk:"id"`
	GroupID   types.Int64  `tfsdk:"group_id"`
	SpaceID   types.Int64  `tfsdk:"space_id"`
	CreatedAt types.String `tfsdk:"created_at"`
	UpdatedAt types.String `tfsdk:"updated_at"`
	UUID      types.String `tfsdk:"uuid"`
	Name      types.String `tfsdk:"name"`
}

func (m *componentGroupResourceModel) toRemoteInput() sbmgmt.ComponentGroupInput {

	return sbmgmt.ComponentGroupInput{
		Name: m.Name.ValueString(),
	}
}

func (m *componentGroupResourceModel) fromRemote(spaceID int64, c *sbmgmt.ComponentGroup) error {
	if c == nil {
		return fmt.Errorf("component-group is nil")
	}
	m.ID = types.StringValue(createIdentifier(spaceID, c.Id))
	m.GroupID = types.Int64Value(c.Id)
	m.CreatedAt = types.StringValue(c.CreatedAt.String())
	m.UpdatedAt = types.StringValue(c.UpdatedAt.String())
	m.Name = types.StringValue(c.Name)
	m.UUID = types.StringValue(c.Uuid.String())
	return nil
}
