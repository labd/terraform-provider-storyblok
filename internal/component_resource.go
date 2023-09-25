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
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/labd/storyblok-go-sdk/sbmgmt"

	"github.com/labd/terraform-provider-storyblok/internal/customvalidators"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &componentResource{}
	_ resource.ResourceWithConfigure   = &componentResource{}
	_ resource.ResourceWithImportState = &componentResource{}
)

// NewComponentResource is a helper function to simplify the provider implementation.
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
		Description: "A component is a standalone entity that is meaningful in its own right. While components (or " +
			"blocks) can be nested in each other, semantically they remain equal. Each component is a small piece " +
			"of your data structure which can be filled with content or nested by your content editor. One component can " +
			"consist of as many field types as required.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The terraform ID of the space role. This is a composite ID, " +
					"and should not be used as reference",
				Computed: true,
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
			"component_group_uuid": schema.StringAttribute{
				Description: "The UUID of the component group.",
				Optional:    true,
				Validators: []validator.String{
					customvalidators.UUID(),
				},
			},
			"icon": schema.StringAttribute{
				Description: "The Icon of the component",
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.OneOf(getComponentIcons()...),
				},
			},
			"image": schema.StringAttribute{
				Description: "An image url of the component",
				Optional:    true,
			},
			"preview_tmpl": schema.StringAttribute{
				Description: "The preview template of the component",
				Optional:    true,
			},
			"preview_field": schema.StringAttribute{
				Description: "A preview field of the component",
				Optional:    true,
			},
			"color": schema.StringAttribute{
				Description: "The background color for the icon of the component",
				Optional:    true,
			},
			"display_name": schema.StringAttribute{
				Description: "The display name of the component",
				Optional:    true,
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
						"add_https": schema.BoolAttribute{
							Description: "Prepends https: to stop usage of relative protocol",
							Optional:    true,
						},
						"allow_target_blank": schema.BoolAttribute{
							Description: "Allows to open links in a new tab for Richtext; Default: false",
							Optional:    true,
						},
						"asset_folder_id": schema.Int64Attribute{
							Description: "Default asset folder numeric id to store uploaded image of that field",
							Optional:    true,
						},
						"can_sync": schema.BoolAttribute{
							Description: "Advanced usage to sync with field in preview; Default: false",
							Optional:    true,
						},
						"customize_toolbar": schema.BoolAttribute{
							Description: "Allow to customize the Markdown or Richtext toolbar; Default: false",
							Optional:    true,
						},
						"component_whitelist": schema.ListAttribute{
							Description: "Array of component/content type names: [\"post\",\"page\",\"product\"]",
							Optional:    true,
							ElementType: types.StringType,
						},
						"datasource_slug": schema.StringAttribute{
							Description: "Define selectable datasources string; Effects editor only if source=internal",
							Optional:    true,
						},
						"default_value": schema.StringAttribute{
							Description: "Default value for the field; Can be an escaped JSON object",
							Optional:    true,
						},
						"description": schema.StringAttribute{
							Description: "Description shown in the editor interface",
							Optional:    true,
						},
						"disable_time": schema.BoolAttribute{
							Description: "Disables time selection from date picker; Default: false",
							Optional:    true,
						},
						"display_name": schema.StringAttribute{
							Description: "Display name shown in the editor interface",
							Optional:    true,
						},
						"external_datasource": schema.StringAttribute{
							Description: "Define external datasource JSON Url; Effects editor only if source=external",
							Optional:    true,
						},
						"field_type": schema.StringAttribute{
							Description: "Name of the custom field type plugin",
							Optional:    true,
						},
						"filetypes": schema.ListAttribute{
							Description: "Array of file type names: [\"images\", \"videos\", \"audios\", \"texts\"]",
							Optional:    true,
							ElementType: types.StringType,
						},
						"force_link_scope": schema.BoolAttribute{
							Description: "Force link scope to be internal_stories; Default: false",
							Optional:    true,
						},
						"folder_slug": schema.StringAttribute{
							Description: "Filter on selectable stories path; Effects editor only if source=internal_stories; In case you have a multi-language folder structure you can add the '{0}' placeholder and the path will be adapted dynamically. Examples: \"{0}/categories/\", {0}/{1}/categories/",
							Optional:    true,
						},
						"image_crop": schema.BoolAttribute{
							Description: "Activate force crop for images: (true/false)",
							Optional:    true,
						},
						"image_height": schema.StringAttribute{
							Description: "Define height in px or height ratio if keep_image_size is enabled",
							Optional:    true,
						},
						"image_width": schema.StringAttribute{
							Description: "Define width in px or width ratio if keep_image_size is enabled",
							Optional:    true,
						},
						"keep_image_size": schema.BoolAttribute{
							Description: "Keep original size: (true/false)",
							Optional:    true,
						},
						"keys": schema.ListAttribute{
							Description: "Array of field keys to include in this section",
							Optional:    true,
							ElementType: types.StringType,
						},
						"link_scope": schema.StringAttribute{
							Description: "A path to a folder to restrict the link scope",
							Optional:    true,
						},
						"max_length": schema.Int64Attribute{
							Description: "Set the max length of the input string",
							Optional:    true,
						},
						"maximum": schema.Int64Attribute{
							Description: "Maximum amount of added bloks in this blok field",
							Optional:    true,
						},
						"minimum": schema.Int64Attribute{
							Description: "Minimum amount of added bloks in this blok field",
							Optional:    true,
						},
						"no_translate": schema.BoolAttribute{
							Description: "Should be excluded in translation export",
							Optional:    true,
						},
						"options": schema.ListNestedAttribute{
							Description: "Array of datasource entries [{name:\"\", value:\"\"}]; Effects editor only if source=undefined",
							Optional:    true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"name": schema.StringAttribute{
										Description: "Name of the datasource entry",
										Required:    true,
									},
									"value": schema.StringAttribute{
										Description: "Value of the datasource entry",
										Required:    true,
									},
								},
							},
						},
						"regex": schema.StringAttribute{
							Description: "Client Regex validation for the field",
							Optional:    true,
						},
						"required": schema.BoolAttribute{
							Description: "Is field required; Default: false",
							Optional:    true,
						},
						"restrict_components": schema.BoolAttribute{
							Description: "Activate restriction nestable component option; Default: false",
							Optional:    true,
						},
						"restrict_content_types": schema.BoolAttribute{
							Description: "Activate restriction content type option",
							Optional:    true,
						},
						"rich_markdown": schema.BoolAttribute{
							Description: "Enable rich markdown view by default (true/false)",
							Optional:    true,
						},
						"rtl": schema.BoolAttribute{
							Description: "Enable global RTL for this field",
							Optional:    true,
						},
						"source": schema.StringAttribute{
							Description: "Possible values: undefined: Self; internal_stories: Stories; internal: Datasource; external: API Endpoint in Datasource Entries Array Format",
							Optional:    true,
						},
						"tooltip": schema.BoolAttribute{
							Description: "Show the description as a tooltip",
							Optional:    true,
						},
						"translatable": schema.BoolAttribute{
							Description: "Can field be translated; Default: false",
							Optional:    true,
						},
						"toolbar": schema.ListAttribute{
							Description: "Array of toolbar keys to include in the Richtext or Markdown toolbar",
							Optional:    true,
							ElementType: types.StringType,
						},
						"use_uuid": schema.BoolAttribute{
							Description: "Default: true; available in option and source=internal_stories",
							Optional:    true,
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

	r.client = getClient(req.ProviderData)
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
	if d := checkGetError("component", componentId, content, err); d != nil {
		resp.Diagnostics.Append(d)
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
	input := plan.toUpdateInput()
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
