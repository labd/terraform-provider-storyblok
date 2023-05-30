package customvalidators

import (
	"context"

	"github.com/gofrs/uuid"
	"github.com/hashicorp/terraform-plugin-framework-validators/helpers/validatordiag"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

var _ validator.String = uuidValidator{}

// uuidValidator validates that a string Attribute's value is a valid UUID.
type uuidValidator struct{}

// Description describes the validation in plain text formatting.
func (validator uuidValidator) Description(_ context.Context) string {
	return "value must be a valid uuid"
}

// MarkdownDescription describes the validation in Markdown formatting.
func (validator uuidValidator) MarkdownDescription(ctx context.Context) string {
	return validator.Description(ctx)
}

// Validate performs the validation.
func (v uuidValidator) ValidateString(ctx context.Context, request validator.StringRequest, response *validator.StringResponse) {
	if request.ConfigValue.IsNull() || request.ConfigValue.IsUnknown() {
		return
	}

	value := request.ConfigValue.ValueString()

	if _, err := uuid.FromString(value); err != nil {
		response.Diagnostics.Append(validatordiag.InvalidAttributeValueMatchDiagnostic(
			request.Path,
			v.Description(ctx),
			value,
		))
	}
}

// UUID returns an AttributeValidator which ensures that any configured
// attribute value is a valid UUID.
//
// Null (unconfigured) and unknown (known after apply) values are skipped.
func UUID() validator.String {
	return uuidValidator{}
}
