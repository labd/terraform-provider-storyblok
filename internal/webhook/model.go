package webhook

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/labd/storyblok-go-sdk/sbmgmt"

	"github.com/labd/terraform-provider-storyblok/internal/utils"
)

type WebhookModel struct {
	ID          types.String   `tfsdk:"id"`
	WebhookID   types.Int64    `tfsdk:"webhook_id"`
	Description types.String   `tfsdk:"description"`
	SpaceID     types.Int64    `tfsdk:"space_id"`
	Actions     []types.String `tfsdk:"actions"`
	Name        types.String   `tfsdk:"name"`
	Endpoint    types.String   `tfsdk:"endpoint"`
	Activated   types.Bool     `tfsdk:"activated"`
	Secret      types.String   `tfsdk:"secret"`
}

func (m *WebhookModel) toCreateInput() sbmgmt.CreateWebhookJSONRequestBody {
	return sbmgmt.CreateWebhookJSONRequestBody{
		WebhookEndpoint: sbmgmt.WebhookCreateInput{
			Activated:   m.Activated.ValueBool(),
			Name:        m.Name.ValueString(),
			Endpoint:    m.Endpoint.ValueString(),
			Actions:     *utils.ConvertToPointerStringSlice(m.Actions),
			Secret:      m.Secret.ValueString(),
			Description: m.Description.ValueStringPointer(),
		},
	}
}
func (m *WebhookModel) toUpdateInput() sbmgmt.UpdateWebhookJSONRequestBody {
	return sbmgmt.UpdateWebhookJSONRequestBody{
		WebhookEndpoint: sbmgmt.WebhookUpdateInput{
			Activated:   m.Activated.ValueBool(),
			Name:        m.Name.ValueString(),
			Endpoint:    m.Endpoint.ValueString(),
			Actions:     *utils.ConvertToPointerStringSlice(m.Actions),
			Secret:      m.Secret.ValueString(),
			Description: m.Description.ValueStringPointer(),
		},
	}
}

func (m *WebhookModel) fromRemote(spaceID int64, i sbmgmt.Webhook) error {
	m.ID = types.StringValue(utils.CreateIdentifier(spaceID, i.Id))
	m.SpaceID = types.Int64Value(spaceID)
	m.WebhookID = types.Int64Value(i.Id)
	m.Name = types.StringValue(i.Name)
	return nil
}
