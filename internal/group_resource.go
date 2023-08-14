package internal

import (
	"context"

	"github.com/davecgh/go-spew/spew"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/labd/storyblok-go-sdk/sbmgmt"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &componentGroupResource{}
	_ resource.ResourceWithConfigure   = &componentGroupResource{}
	_ resource.ResourceWithImportState = &componentGroupResource{}
)

// NewComponentGroupResource is a helper function to simplify the provider implementation.
func NewComponentGroupResource() resource.Resource {
	return &componentGroupResource{}
}

// componentGroupResource is the resource implementation.
type componentGroupResource struct {
	client sbmgmt.ClientWithResponsesInterface
}

// Metadata returns the data source type name.
func (r *componentGroupResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_component_group"
}

// Schema defines the schema for the data source.
func (r *componentGroupResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manage a component.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The terraform ID of the component.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"group_id": schema.Int64Attribute{
				Description: "The ID of the component group.",
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"space_id": schema.Int64Attribute{
				Description: "The ID of the space.",
				Required:    true,
			},
			"uuid": schema.StringAttribute{
				Description: "The UUID of the component group.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description: "The name of the component group.",
				Required:    true,
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (r *componentGroupResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = getClient(req.ProviderData)
}

// Create creates the resource and sets the initial Terraform state.
func (r *componentGroupResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan componentGroupResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	input := plan.toCreateInput()
	spaceID := plan.SpaceID.ValueInt64()

	content, err := r.client.CreateComponentGroupWithResponse(ctx, spaceID, input)
	if d := checkCreateError("component_group", content, err); d != nil {
		resp.Diagnostics.Append(d)
		return
	}

	component := content.JSON201.ComponentGroup
	tflog.Debug(ctx, spew.Sdump(component))

	// Map response body to schema and populate Computed attribute values
	if err := plan.fromRemote(spaceID, component); err != nil {
		resp.Diagnostics.AddError(
			"Error creating component-group",
			"Could not create component-group, unexpected error: "+err.Error(),
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
func (r *componentGroupResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state componentGroupResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	spaceId, groupId := parseIdentifier(state.ID.ValueString())

	// Get refreshed order value from HashiCups
	content, err := r.client.GetComponentGroupWithResponse(ctx, spaceId, groupId)
	if d := checkGetError("component_group", content, err); d != nil {
		resp.Diagnostics.Append(d)
		return
	}

	group := content.JSON200.ComponentGroup

	// Overwrite items with refreshed state
	if err := state.fromRemote(spaceId, group); err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Storyblok Component",
			"Could not read Storyblok component ID "+state.ID.ValueString()+": "+err.Error(),
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
func (r *componentGroupResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan componentGroupResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	input := plan.toUpdateInput()
	spaceID := plan.SpaceID.ValueInt64()

	content, err := r.client.UpdateComponentGroupWithResponse(ctx, spaceID, plan.GroupID.ValueInt64(), input)
	if d := checkUpdateError("component_group", content, err); d != nil {
		resp.Diagnostics.Append(d)
		return
	}

	group := content.JSON200.ComponentGroup
	tflog.Debug(ctx, spew.Sdump(group))

	// Map response body to schema and populate Computed attribute values
	if err := plan.fromRemote(spaceID, group); err != nil {
		resp.Diagnostics.AddError(
			"Error creating component-group",
			"Could not create component-group, unexpected error: "+err.Error(),
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
func (r *componentGroupResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state componentGroupResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	spaceId, groupId := parseIdentifier(state.ID.ValueString())
	content, err := r.client.DeleteComponentGroupWithResponse(ctx, spaceId, groupId)
	if d := checkDeleteError("component_group", content, err); d != nil {
		resp.Diagnostics.Append(d)
		return
	}
}

func (r *componentGroupResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
