package internal

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/labd/storyblok-go-sdk/sbmgmt"

	"github.com/labd/terraform-provider-storyblok/internal/utils"
)

// spaceRoleResourceModel maps the resource schema data.
type spaceRoleResourceModel struct {
	ID      types.String `tfsdk:"id"`
	RoleID  types.Int64  `tfsdk:"role_id"`
	SpaceID types.Int64  `tfsdk:"space_id"`

	AllowedLanguages         []types.String `tfsdk:"allowed_languages"`
	AllowedPaths             []types.String `tfsdk:"allowed_paths"`
	BranchIds                []types.Int64  `tfsdk:"branch_ids"`
	ComponentIds             []types.Int64  `tfsdk:"component_ids"`
	DatasourceIds            []types.Int64  `tfsdk:"datasource_ids"`
	ExternalID               types.String   `tfsdk:"external_id"`
	FieldPermissions         []types.String `tfsdk:"field_permissions"`
	Permissions              []types.String `tfsdk:"permissions"`
	ReadonlyFieldPermissions []types.String `tfsdk:"readonly_field_permissions"`
	ResolvedAllowedPaths     []types.String `tfsdk:"resolved_allowed_paths"`
	Role                     types.String   `tfsdk:"role"`
	Subtitle                 types.String   `tfsdk:"subtitle"`
}

func (m *spaceRoleResourceModel) toCreateInput() sbmgmt.SpaceRoleCreateInput {
	return sbmgmt.SpaceRoleCreateInput{
		SpaceRole: &sbmgmt.SpaceRoleBase{
			AllowedLanguages:         utils.ConvertToPointerStringSlice(m.AllowedLanguages),
			AllowedPaths:             utils.ConvertToPointerStringSlice(m.AllowedPaths),
			BranchIds:                utils.ConvertToPointerIntSlice(m.BranchIds),
			ComponentIds:             utils.ConvertToPointerIntSlice(m.ComponentIds),
			DatasourceIds:            utils.ConvertToPointerIntSlice(m.DatasourceIds),
			ExtId:                    m.ExternalID.ValueStringPointer(),
			FieldPermissions:         utils.ConvertToPointerStringSlice(m.FieldPermissions),
			Permissions:              utils.ConvertToPointerStringSlice(m.Permissions),
			ReadonlyFieldPermissions: utils.ConvertToPointerStringSlice(m.ReadonlyFieldPermissions),
			ResolvedAllowedPaths:     utils.ConvertToPointerStringSlice(m.ResolvedAllowedPaths),
			Role:                     m.Role.ValueString(),
			Subtitle:                 m.Subtitle.ValueStringPointer(),
		},
	}
}
func (m *spaceRoleResourceModel) toUpdateInput() sbmgmt.SpaceRoleUpdateInput {
	return sbmgmt.SpaceRoleUpdateInput{
		SpaceRole: &sbmgmt.SpaceRoleBase{
			AllowedLanguages:         utils.ConvertToPointerStringSlice(m.AllowedLanguages),
			AllowedPaths:             utils.ConvertToPointerStringSlice(m.AllowedPaths),
			BranchIds:                utils.ConvertToPointerIntSlice(m.BranchIds),
			ComponentIds:             utils.ConvertToPointerIntSlice(m.ComponentIds),
			DatasourceIds:            utils.ConvertToPointerIntSlice(m.DatasourceIds),
			ExtId:                    m.ExternalID.ValueStringPointer(),
			FieldPermissions:         utils.ConvertToPointerStringSlice(m.FieldPermissions),
			Permissions:              utils.ConvertToPointerStringSlice(m.Permissions),
			ReadonlyFieldPermissions: utils.ConvertToPointerStringSlice(m.ReadonlyFieldPermissions),
			ResolvedAllowedPaths:     utils.ConvertToPointerStringSlice(m.ResolvedAllowedPaths),
			Role:                     m.Role.ValueString(),
			Subtitle:                 m.Subtitle.ValueStringPointer(),
		},
	}
}

func (m *spaceRoleResourceModel) fromRemote(spaceId int64, c *sbmgmt.SpaceRole) error {
	if c == nil {
		return fmt.Errorf("space role is nil")
	}
	m.ID = types.StringValue(utils.CreateIdentifier(spaceId, int64(c.Id)))
	m.RoleID = types.Int64Value(int64(c.Id))
	return nil
}
