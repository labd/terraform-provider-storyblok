package internal

import (
	"context"

	"github.com/davecgh/go-spew/spew"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/labd/storyblok-go-sdk/sbmgmt"

	"github.com/labd/terraform-provider-storyblok/internal/utils"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &spaceRoleResource{}
	_ resource.ResourceWithConfigure   = &spaceRoleResource{}
	_ resource.ResourceWithImportState = &spaceRoleResource{}
)

// NewSpaceRoleResource is a helper function to simplify the provider implementation.
func NewSpaceRoleResource() resource.Resource {
	return &spaceRoleResource{}
}

// spaceRoleResource is the resource implementation.
type spaceRoleResource struct {
	client sbmgmt.ClientWithResponsesInterface
}

// Metadata returns the data source type name.
func (r *spaceRoleResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_space_role"
}

// Schema defines the schema for the data source.
func (r *spaceRoleResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Space roles are custom permission sets that can be attached to collaborators to define their roles " +
			"and permissions in a specific space.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The terraform ID of the space role. This is a composite ID, " +
					"and should not be used as reference",
				Computed: true,
			},
			"role_id": schema.Int64Attribute{
				Description: "The ID of the role.",
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"space_id": schema.Int64Attribute{
				Description: "The ID of the space.",
				Required:    true,
			},
			"external_id": schema.StringAttribute{
				Description: "External ID (used for SSO)",
				Optional:    true,
			},
			"role": schema.StringAttribute{
				Description: "Name used in the interface",
				Required:    true,
			},
			"subtitle": schema.StringAttribute{
				Description: "A short description of the role",
				Optional:    true,
			},
			"allowed_languages": schema.ListAttribute{
				Description: "Add languages the user should have access to (acts as allow list). If no item is selected " +
					"the user has rights to edit all content.",
				Optional:    true,
				ElementType: types.StringType,
			},
			"allowed_paths": schema.ListAttribute{
				Description: "Story ids the user should have access to (acts as whitelist). If no item is selected the " +
					"user has rights to access all content items.",
				Optional:    true,
				ElementType: types.StringType,
			},
			"resolved_allowed_paths": schema.ListAttribute{
				Description: "Resolved allowed_paths for displaying paths",
				Optional:    true,
				ElementType: types.StringType,
			},
			"field_permissions": schema.ListAttribute{
				Description: "Hide specific fields for this user with an array of strings with the schema: " +
					"\"component_name.field_name\"",
				Optional:    true,
				ElementType: types.StringType,
			},
			"permissions": schema.ListAttribute{
				Description: "Allow specific actions in interface by adding the permission as array of strings",
				Optional:    true,
				ElementType: types.StringType,
			},
			"readonly_field_permissions": schema.ListAttribute{
				Description: "Read only field permissions",
				Optional:    true,
				ElementType: types.StringType,
			},
			"branch_ids": schema.ListAttribute{
				Description: "Branch ids that the role is allowed access to",
				Optional:    true,
				ElementType: types.Int64Type,
			},
			"component_ids": schema.ListAttribute{
				Description: "Component ids that the role is allowed access to",
				Optional:    true,
				ElementType: types.Int64Type,
			},
			"datasource_ids": schema.ListAttribute{
				Description: "Datasource ids that the role is allowed access to",
				Optional:    true,
				ElementType: types.Int64Type,
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (r *spaceRoleResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = utils.GetClient(req.ProviderData)
}

// Create creates the resource and sets the initial Terraform state.
func (r *spaceRoleResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan spaceRoleResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	input := plan.toCreateInput()
	spaceID := plan.SpaceID.ValueInt64()

	content, err := r.client.CreateSpaceRoleWithResponse(ctx, spaceID, input)
	if d := utils.CheckCreateError("space_role", content, err); d != nil {
		resp.Diagnostics.Append(d)
		return
	}

	component := content.JSON201.SpaceRole
	tflog.Debug(ctx, spew.Sdump(component))

	// Map response body to schema and populate Computed attribute values
	if err := plan.fromRemote(spaceID, component); err != nil {
		resp.Diagnostics.AddError(
			"Error creating Space Role",
			"Could not create Space Role, unexpected error: "+err.Error(),
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
func (r *spaceRoleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state spaceRoleResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	spaceId, groupId := utils.ParseIdentifier(state.ID.ValueString())

	content, err := r.client.GetSpaceRoleWithResponse(ctx, spaceId, groupId)
	if d := utils.CheckGetError("space_role", groupId, content, err); d != nil {
		resp.Diagnostics.Append(d)
		return
	}

	spaceRole := content.JSON200.SpaceRole

	// Overwrite items with refreshed state
	if err := state.fromRemote(spaceId, spaceRole); err != nil {
		resp.Diagnostics.AddError(
			"Error reading Space Role",
			"Could not read Space Role "+state.ID.ValueString()+": "+err.Error(),
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
func (r *spaceRoleResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan spaceRoleResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	input := plan.toUpdateInput()
	spaceID := plan.SpaceID.ValueInt64()

	content, err := r.client.UpdateSpaceRoleWithResponse(ctx, spaceID, plan.RoleID.ValueInt64(), input)
	if d := utils.CheckUpdateError("space_role", content, err); d != nil {
		resp.Diagnostics.Append(d)
		return
	}

	spaceRole := content.JSON200.SpaceRole
	tflog.Debug(ctx, spew.Sdump(spaceRole))

	// Map response body to schema and populate Computed attribute values
	if err := plan.fromRemote(spaceID, spaceRole); err != nil {
		resp.Diagnostics.AddError(
			"Error creating Space Role",
			"Could not create Space role, unexpected error: "+err.Error(),
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
func (r *spaceRoleResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state spaceRoleResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	spaceId, spaceRoleId := utils.ParseIdentifier(state.ID.ValueString())
	content, err := r.client.DeleteSpaceRoleWithResponse(ctx, spaceId, spaceRoleId)
	if d := utils.CheckDeleteError("space_role", content, err); d != nil {
		resp.Diagnostics.Append(d)
		return
	}
}

func (r *spaceRoleResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
