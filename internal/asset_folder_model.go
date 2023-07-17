package internal

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/labd/storyblok-go-sdk/sbmgmt"
)

// assetFolderResourceModel maps the resource schema data.
type assetFolderResourceModel struct {
	ID            types.String `tfsdk:"id"`
	AssetFolderID types.Int64  `tfsdk:"asset_folder_id"`
	SpaceID       types.Int64  `tfsdk:"space_id"`
	Name          types.String `tfsdk:"name"`
	ParentID      types.Int64  `tfsdk:"parent_id"`
}

func (m *assetFolderResourceModel) toCreateInput() sbmgmt.AssetFolderCreateInput {
	return sbmgmt.AssetFolderCreateInput{
		AssetFolder: sbmgmt.AssetFolderBase{
			Name:     m.Name.ValueString(),
			ParentId: m.ParentID.ValueInt64Pointer(),
		},
	}
}
func (m *assetFolderResourceModel) toUpdateInput() sbmgmt.AssetFolderUpdateInput {
	return sbmgmt.UpdateAssetFolderJSONRequestBody{
		AssetFolder: sbmgmt.AssetFolderBase{
			Name:     m.Name.ValueString(),
			ParentId: m.ParentID.ValueInt64Pointer(),
		},
	}
}

func (m *assetFolderResourceModel) fromRemote(spaceID int64, f *sbmgmt.AssetFolder) error {
	if f == nil {
		return fmt.Errorf("asset folder is nil")
	}
	m.ID = types.StringValue(createIdentifier(spaceID, f.Id))
	m.SpaceID = types.Int64Value(spaceID)
	m.Name = types.StringValue(f.Name)
	if f.ParentId != nil {
		m.ParentID = types.Int64Value(*f.ParentId)
	}
	return nil
}
