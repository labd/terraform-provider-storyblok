package internal

import (
	"context"
	"fmt"
	"net/http"

	"github.com/davecgh/go-spew/spew"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/labd/storyblok-go-sdk/sbmgmt"

	"github.com/labd/terraform-provider-storyblok/internal/utils"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &assetFolderResource{}
	_ resource.ResourceWithConfigure   = &assetFolderResource{}
	_ resource.ResourceWithImportState = &assetFolderResource{}
)

// NewAssetFolderResource is a helper function to simplify the provider implementation.
func NewAssetFolderResource() resource.Resource {
	return &assetFolderResource{}
}

// assetFolderResource is the resource implementation.
type assetFolderResource struct {
	client sbmgmt.ClientWithResponsesInterface
}

// Metadata returns the data source type name.
func (r *assetFolderResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_asset_folder"
}

// Schema defines the schema for the data source.
func (r *assetFolderResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Asset folder allow you to group your assets. Besides the overall root folder you can define nested " +
			"folder structures.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The terraform ID of the space role. This is a composite ID, " +
					"and should not be used as reference",
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"asset_folder_id": schema.Int64Attribute{
				Description: "The ID of the asset folder.",
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description: "The technical name of the asset folder.",
				Required:    true,
			},
			"space_id": schema.Int64Attribute{
				Description: "The ID of the space.",
				Required:    true,
			},
			"parent_id": schema.Int64Attribute{
				Description: "The ID of the parent asset folder.",
				Optional:    true,
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (r *assetFolderResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = utils.GetClient(req.ProviderData)
}

// Create creates the resource and sets the initial Terraform state.
func (r *assetFolderResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan assetFolderResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	input := plan.toCreateInput()
	spaceID := plan.SpaceID.ValueInt64()

	content, err := r.client.CreateAssetFolderWithResponse(ctx, spaceID, input)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating assetFolder",
			"Could not create assetFolder, unexpected error: "+err.Error(),
		)
		return
	}
	if content.StatusCode() != http.StatusCreated {
		resp.Diagnostics.AddError(
			"Error creating assetFolder",
			fmt.Sprintf(
				"Could not create assetFolder, status code %d error: %s",
				content.StatusCode(), string(content.Body)),
		)
		return
	}

	assetFolder := content.JSON201.AssetFolder
	tflog.Debug(ctx, spew.Sdump(assetFolder))

	// Map response body to schema and populate Computed attribute values
	if err := plan.fromRemote(spaceID, assetFolder); err != nil {
		resp.Diagnostics.AddError(
			"Error creating assetFolder",
			"Could not create assetFolder, unexpected error: "+err.Error(),
		)
		return
	}

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *assetFolderResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state assetFolderResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	spaceId, id := utils.ParseIdentifier(state.ID.ValueString())

	content, err := r.client.GetAssetFolderWithResponse(ctx, spaceId, id)
	if d := utils.CheckGetError("assetFolder", id, content, err); d != nil {
		resp.Diagnostics.Append(d)
		return
	}

	assetFolder := content.JSON200.AssetFolder

	// Overwrite items with refreshed state
	if err := state.fromRemote(spaceId, assetFolder); err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Storyblok AssetFolder",
			"Could not read Storyblok assetFolder ID "+state.ID.ValueString()+": "+err.Error(),
		)
		return
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *assetFolderResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan assetFolderResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	input := plan.toUpdateInput()
	spaceID := plan.SpaceID.ValueInt64()

	content, err := r.client.UpdateAssetFolderWithResponse(ctx, spaceID, plan.AssetFolderID.ValueInt64(), input)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating assetFolder",
			"Could not update assetFolder, unexpected error: "+err.Error(),
		)
		return
	}
	if content.StatusCode() != http.StatusNoContent {
		resp.Diagnostics.AddError(
			"Error updating assetFolder",
			fmt.Sprintf(
				"Could not update assetFolder, status code %d error: %s",
				content.StatusCode(), string(content.Body)),
		)
		return
	}

	afResp, err := r.client.GetAssetFolderWithResponse(ctx, spaceID, plan.AssetFolderID.ValueInt64())
	if d := utils.CheckGetError("assetFolder", plan.AssetFolderID.ValueInt64(), afResp, err); d != nil {
		resp.Diagnostics.Append(d)
		return
	}

	assetFolder := afResp.JSON200.AssetFolder

	// Map response body to schema and populate Computed attribute values
	if err := plan.fromRemote(spaceID, assetFolder); err != nil {
		resp.Diagnostics.AddError(
			"Error updating assetFolder",
			"Could not update assetFolder, unexpected error: "+err.Error(),
		)
		return
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *assetFolderResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state assetFolderResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	spaceId, assetFolderId := utils.ParseIdentifier(state.ID.ValueString())
	content, err := r.client.DeleteAssetFolderWithResponse(ctx, spaceId, assetFolderId)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting assetFolder",
			"Could not delete assetFolder, unexpected error: "+err.Error(),
		)
		return
	}
	if content.StatusCode() != http.StatusNoContent {
		resp.Diagnostics.AddError(
			"Error deleting assetFolder",
			fmt.Sprintf(
				"Could not delete assetFolder, status code %d error: %s",
				content.StatusCode(), string(content.Body)),
		)
		return
	}
}

func (r *assetFolderResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
