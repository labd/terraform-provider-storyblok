package webhook

import (
	"context"
	"fmt"
	"net/http"

	"github.com/davecgh/go-spew/spew"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/labd/storyblok-go-sdk/sbmgmt"

	"github.com/labd/terraform-provider-storyblok/internal/utils"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &webhookResource{}
	_ resource.ResourceWithConfigure   = &webhookResource{}
	_ resource.ResourceWithImportState = &webhookResource{}
)

// NewWebhookResource is a helper function to simplify the provider implementation.
func NewWebhookResource() resource.Resource {
	return &webhookResource{}
}

// webhookResource is the resource implementation.
type webhookResource struct {
	client sbmgmt.ClientWithResponsesInterface
}

// Metadata returns the data source type name.
func (r *webhookResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_webhook"
}

// Schema defines the schema for the data source.
func (r *webhookResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Webhooks are used to send Storyblok events to other applications. There are some default " +
			"Storyblok events that you can listen to when they are triggered. Read about " +
			"[Available Triggers](https://www.storyblok.com/docs/concepts/webhooks#setup) to learn more.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The terraform ID of the space role. This is a composite ID, " +
					"and should not be used as reference",
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"webhook_id": schema.Int64Attribute{
				Description: "The ID of the webhook.",
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description: "The technical name of the webhook.",
				Required:    true,
			},
			"space_id": schema.Int64Attribute{
				Description: "Numeric ID of a space.",
				Required:    true,
			},
			"endpoint": schema.StringAttribute{
				Description: "The endpoint URL to send the webhook to.",
				Required:    true,
			},
			"activated": schema.BoolAttribute{
				Description: "Whether the webhook is activated.",
				Required:    false,
				Default:     booldefault.StaticBool(true),
				Computed:    true,
			},
			"actions": schema.ListAttribute{
				Description: "The actions that should trigger the webhook.",
				Required:    true,
				ElementType: types.StringType,
			},
			"secret": schema.StringAttribute{
				Description: "The secret to sign the webhook payload with.",
				Optional:    true,
				Sensitive:   true,
			},
			"description": schema.StringAttribute{
				Description: "The description of the webhook.",
				Optional:    true,
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (r *webhookResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = utils.GetClient(req.ProviderData)
}

// Create creates the resource and sets the initial Terraform state.
func (r *webhookResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan WebhookModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	input := plan.toCreateInput()
	spaceID := plan.SpaceID.ValueInt64()

	content, err := r.client.CreateWebhookWithResponse(ctx, spaceID, input)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating webhook",
			"Could not create webhook, unexpected error: "+err.Error(),
		)
		return
	}
	if content.StatusCode() != http.StatusCreated {
		resp.Diagnostics.AddError(
			"Error creating webhook",
			fmt.Sprintf(
				"Could not create webhook, status code %d error: %s",
				content.StatusCode(), string(content.Body)),
		)
		return
	}

	webhook := content.JSON201.WebhookEndpoint
	tflog.Debug(ctx, spew.Sdump(webhook))

	// Map response body to schema and populate Computed attribute values
	if err := plan.fromRemote(spaceID, webhook); err != nil {
		resp.Diagnostics.AddError(
			"Error creating webhook",
			"Could not create webhook, unexpected error: "+err.Error(),
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
func (r *webhookResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state WebhookModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	spaceId, id := utils.ParseIdentifier(state.ID.ValueString())

	content, err := r.client.GetWebhookWithResponse(ctx, spaceId, id)
	if d := utils.CheckGetError("webhook", id, content, err); d != nil {
		resp.Diagnostics.Append(d)
		return
	}

	webhook := content.JSON200.WebhookEndpoint

	// Overwrite items with refreshed state
	if err := state.fromRemote(spaceId, webhook); err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Storyblok Webhook",
			"Could not read Storyblok webhook ID "+state.ID.ValueString()+": "+err.Error(),
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
func (r *webhookResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan WebhookModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	input := plan.toUpdateInput()
	spaceID := plan.SpaceID.ValueInt64()

	content, err := r.client.UpdateWebhookWithResponse(ctx, spaceID, plan.WebhookID.ValueInt64(), input)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating webhook",
			"Could not update webhook, unexpected error: "+err.Error(),
		)
		return
	}
	if content.StatusCode() != http.StatusOK {
		resp.Diagnostics.AddError(
			"Error updating webhook",
			fmt.Sprintf(
				"Could not update webhook, status code %d error: %s",
				content.StatusCode(), string(content.Body)),
		)
		return
	}

	afResp, err := r.client.GetWebhookWithResponse(ctx, spaceID, plan.WebhookID.ValueInt64())
	if d := utils.CheckGetError("webhook", plan.WebhookID.ValueInt64(), afResp, err); d != nil {
		resp.Diagnostics.Append(d)
		return
	}

	webhook := afResp.JSON200.WebhookEndpoint

	// Map response body to schema and populate Computed attribute values
	if err := plan.fromRemote(spaceID, webhook); err != nil {
		resp.Diagnostics.AddError(
			"Error updating webhook",
			"Could not update webhook, unexpected error: "+err.Error(),
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
func (r *webhookResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state WebhookModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	spaceId, webhookId := utils.ParseIdentifier(state.ID.ValueString())
	content, err := r.client.DeleteWebhookWithResponse(ctx, spaceId, webhookId)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting webhook",
			"Could not delete webhook, unexpected error: "+err.Error(),
		)
		return
	}
	if content.StatusCode() != http.StatusNoContent {
		resp.Diagnostics.AddError(
			"Error deleting webhook",
			fmt.Sprintf(
				"Could not delete webhook, status code %d error: %s",
				content.StatusCode(), string(content.Body)),
		)
		return
	}
}

func (r *webhookResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
