package internal

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/labd/storyblok-go-sdk/sbmgmt"

	"github.com/labd/terraform-provider-storyblok/internal/utils"
)

// componentGroupResourceModel maps the resource schema data.
type componentGroupResourceModel struct {
	ID      types.String `tfsdk:"id"`
	GroupID types.Int64  `tfsdk:"group_id"`
	SpaceID types.Int64  `tfsdk:"space_id"`
	UUID    types.String `tfsdk:"uuid"`
	Name    types.String `tfsdk:"name"`
}

func (m *componentGroupResourceModel) toCreateInput() sbmgmt.ComponentGroupCreateInput {

	return sbmgmt.ComponentGroupCreateInput{
		ComponentGroup: sbmgmt.ComponentGroupBase{
			Name: m.Name.ValueString(),
		},
	}
}
func (m *componentGroupResourceModel) toUpdateInput() sbmgmt.ComponentGroupUpdateInput {

	return sbmgmt.ComponentGroupUpdateInput{
		ComponentGroup: sbmgmt.ComponentGroupBase{
			Name: m.Name.ValueString(),
		},
	}
}

func (m *componentGroupResourceModel) fromRemote(spaceID int64, c *sbmgmt.ComponentGroup) error {
	if c == nil {
		return fmt.Errorf("component-group is nil")
	}
	m.ID = types.StringValue(utils.CreateIdentifier(spaceID, c.Id))
	m.GroupID = types.Int64Value(c.Id)
	m.Name = types.StringValue(c.Name)
	m.UUID = types.StringValue(c.Uuid.String())
	return nil
}
