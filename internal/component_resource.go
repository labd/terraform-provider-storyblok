package internal

import (
	"context"
	"fmt"
	"net/http"

	"github.com/davecgh/go-spew/spew"
	"github.com/elliotchance/pie/v2"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/labd/storyblok-go-sdk/sbmgmt"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &componentResource{}
	_ resource.ResourceWithConfigure   = &componentResource{}
	_ resource.ResourceWithImportState = &componentResource{}
)

// NewcomponentResource is a helper function to simplify the provider implementation.
func NewComponentResource() resource.Resource {
	return &componentResource{}
}

// componentResource is the resource implementation.
type componentResource struct {
	client sbmgmt.ClientWithResponsesInterface
}

// Metadata returns the data source type name.
func (r *componentResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_component"
}

// Schema defines the schema for the data source.
func (r *componentResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
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
			"component_id": schema.Int64Attribute{
				Description: "The ID of the component.",
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"space_id": schema.Int64Attribute{
				Description: "The ID of the space.",
				Required:    true,
			},
			"created_at": schema.StringAttribute{
				Description: "The creation timestamp of the component.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description: "The technical name of the component.",
				Required:    true,
			},
			"is_root": schema.BoolAttribute{
				Description: "Component should be usable as a Content Type",
				Optional:    true,
				Computed:    true,
			},
			"is_nestable": schema.BoolAttribute{
				Description: "Component should be insertable in blocks field type fields",
				Optional:    true,
				Computed:    true,
			},
			"schema": schema.MapNestedAttribute{
				Description: "Schema of this component.",
				Required:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"type": schema.StringAttribute{
							Description: "The type of the field",
							Required:    true,
							Validators: []validator.String{
								stringvalidator.OneOf(pie.Keys(getComponentTypes())...),
							},
						},
						"position": schema.Int64Attribute{
							Description: "The position of the field",
							Required:    true,
						},
					},
				},
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (r *componentResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(sbmgmt.ClientWithResponsesInterface)
}

// Create creates the resource and sets the initial Terraform state.
func (r *componentResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan componentResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	input := plan.toRemoteInput()
	spaceID := plan.SpaceID.ValueInt64()

	content, err := r.client.CreateComponentWithResponse(ctx, spaceID, input)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating component",
			"Could not create component, unexpected error: "+err.Error(),
		)
		return
	}
	if content.StatusCode() != http.StatusCreated {
		resp.Diagnostics.AddError(
			"Error creating component",
			fmt.Sprintf(
				"Could not create component, status code %d error: %s",
				content.StatusCode(), string(content.Body)),
		)
		return
	}

	component := content.JSON201.Component
	tflog.Debug(ctx, spew.Sdump(component))

	// Map response body to schema and populate Computed attribute values
	if err := plan.fromRemote(spaceID, component); err != nil {
		resp.Diagnostics.AddError(
			"Error creating component",
			"Could not create component, unexpected error: "+err.Error(),
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
func (r *componentResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state componentResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	spaceId, componentId := parseIdentifier(state.ID.ValueString())

	// Get refreshed order value from HashiCups
	content, err := r.client.GetComponentWithResponse(ctx, spaceId, componentId)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Storyblok Component",
			"Could not read Storyblok component ID "+state.ID.ValueString()+": "+err.Error(),
		)
		return
	}

	component := content.JSON200.Component

	// Overwrite items with refreshed state
	if err := state.fromRemote(spaceId, component); err != nil {
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
func (r *componentResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan componentResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	input := plan.toRemoteInput()
	spaceID := plan.SpaceID.ValueInt64()

	content, err := r.client.UpdateComponentWithResponse(ctx, spaceID, plan.ComponentID.ValueInt64(), input)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating component",
			"Could not update component, unexpected error: "+err.Error(),
		)
		return
	}
	if content.StatusCode() != http.StatusOK {
		resp.Diagnostics.AddError(
			"Error updating component",
			fmt.Sprintf(
				"Could not update component, status code %d error: %s",
				content.StatusCode(), string(content.Body)),
		)
		return
	}

	component := content.JSON200.Component
	tflog.Debug(ctx, spew.Sdump(component))

	// Map response body to schema and populate Computed attribute values
	if err := plan.fromRemote(spaceID, component); err != nil {
		resp.Diagnostics.AddError(
			"Error creating component",
			"Could not create component, unexpected error: "+err.Error(),
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
func (r *componentResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state componentResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	spaceId, componentId := parseIdentifier(state.ID.ValueString())
	content, err := r.client.DeleteComponentWithResponse(ctx, spaceId, componentId)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting component",
			"Could not delete component, unexpected error: "+err.Error(),
		)
		return
	}
	if content.StatusCode() != http.StatusOK {
		resp.Diagnostics.AddError(
			"Error deleting component",
			fmt.Sprintf(
				"Could not delete component, status code %d error: %s",
				content.StatusCode(), string(content.Body)),
		)
		return
	}

}

func (r *componentResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
