package internal

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/labd/storyblok-go-sdk/sbmgmt"
)

// componentResourceModel maps the resource schema data.
type componentResourceModel struct {
	ID                 types.String          `tfsdk:"id"`
	ComponentID        types.Int64           `tfsdk:"component_id"`
	SpaceID            types.Int64           `tfsdk:"space_id"`
	CreatedAt          types.String          `tfsdk:"created_at"`
	Name               types.String          `tfsdk:"name"`
	IsRoot             types.Bool            `tfsdk:"is_root"`
	IsNestable         types.Bool            `tfsdk:"is_nestable"`
	ComponentGroupUUID types.String          `tfsdk:"component_group_uuid"`
	Schema             map[string]fieldModel `tfsdk:"schema"`
}

type fieldModel struct {
	Type     types.String `tfsdk:"type"`
	Position types.Int64  `tfsdk:"position"`

	AddHttps             types.Bool     `tfsdk:"add_https"`
	AssetFolderId        types.Int64    `tfsdk:"asset_folder_id"`
	CanSync              types.Bool     `tfsdk:"can_sync"`
	ComponentWhitelist   []types.String `tfsdk:"component_whitelist"`
	DatasourceSlug       types.String   `tfsdk:"datasource_slug"`
	DefaultValue         types.String   `tfsdk:"default_value"`
	Description          types.String   `tfsdk:"description"`
	DisableTime          types.Bool     `tfsdk:"disable_time"`
	DisplayName          types.String   `tfsdk:"display_name"`
	ExternalDatasource   types.String   `tfsdk:"external_datasource"`
	FieldType            types.String   `tfsdk:"field_type"`
	Filetypes            []types.String `tfsdk:"filetypes"`
	FolderSlug           types.String   `tfsdk:"folder_slug"`
	ImageCrop            types.Bool     `tfsdk:"image_crop"`
	ImageHeight          types.String   `tfsdk:"image_height"`
	ImageWidth           types.String   `tfsdk:"image_width"`
	KeepImageSize        types.Bool     `tfsdk:"keep_image_size"`
	Keys                 []types.String `tfsdk:"keys"`
	MaxLength            types.Int64    `tfsdk:"max_length"`
	Maximum              types.Int64    `tfsdk:"maximum"`
	NoTranslate          types.Bool     `tfsdk:"no_translate"`
	Options              []optionModel  `tfsdk:"options"`
	PreviewField         types.Bool     `tfsdk:"preview_field"`
	Regex                types.String   `tfsdk:"regex"`
	Required             types.Bool     `tfsdk:"required"`
	RestrictComponents   types.Bool     `tfsdk:"restrict_components"`
	RestrictContentTypes types.Bool     `tfsdk:"restrict_content_types"`
	RichMarkdown         types.Bool     `tfsdk:"rich_markdown"`
	Rtl                  types.Bool     `tfsdk:"rtl"`
	Source               types.String   `tfsdk:"source"`
	Translatable         types.Bool     `tfsdk:"translatable"`
	Tooltip              types.Bool     `tfsdk:"tooltip"`
	UseUuid              types.Bool     `tfsdk:"use_uuid"`
}

type optionModel struct {
	Name  types.String `tfsdk:"name,omitempty"`
	Value types.String `tfsdk:"value,omitempty"`
}

func (m *componentResourceModel) toRemoteInput() sbmgmt.ComponentInput {

	raw := make(map[string]sbmgmt.FieldInput, len(m.Schema))
	for name := range m.Schema {
		item := m.Schema[name]
		raw[name] = sbmgmt.FieldInput{
			Type: item.Type.ValueString(),
			Pos:  item.Position.ValueInt64(),

			AddHttps:             item.AddHttps.ValueBoolPointer(),
			AssetFolderId:        item.AssetFolderId.ValueInt64Pointer(),
			CanSync:              item.CanSync.ValueBoolPointer(),
			DatasourceSlug:       item.DatasourceSlug.ValueStringPointer(),
			DefaultValue:         item.DefaultValue.ValueStringPointer(),
			Description:          item.Description.ValueStringPointer(),
			DisplayName:          item.DisplayName.ValueStringPointer(),
			ComponentWhitelist:   convertToPointerStringSlice(item.ComponentWhitelist),
			ExternalDatasource:   item.ExternalDatasource.ValueStringPointer(),
			FieldType:            item.FieldType.ValueStringPointer(),
			Filetypes:            convertToPointerStringSlice(item.Filetypes),
			FolderSlug:           item.FolderSlug.ValueStringPointer(),
			ImageCrop:            item.ImageCrop.ValueBoolPointer(),
			ImageHeight:          item.ImageHeight.ValueStringPointer(),
			ImageWidth:           item.ImageWidth.ValueStringPointer(),
			KeepImageSize:        item.KeepImageSize.ValueBoolPointer(),
			Keys:                 convertToPointerStringSlice(item.Keys),
			Maximum:              item.Maximum.ValueInt64Pointer(),
			NoTranslate:          item.NoTranslate.ValueBoolPointer(),
			Options:              deserializeOptionsModel(item.Options),
			PreviewField:         item.PreviewField.ValueBoolPointer(),
			Regex:                item.Regex.ValueStringPointer(),
			Required:             item.Required.ValueBoolPointer(),
			RestrictComponents:   item.RestrictComponents.ValueBoolPointer(),
			RestrictContentTypes: item.RestrictContentTypes.ValueBoolPointer(),
			RichMarkdown:         item.RichMarkdown.ValueBoolPointer(),
			Rtl:                  item.Rtl.ValueBoolPointer(),
			Source:               item.Source.ValueStringPointer(),
			Tooltip:              item.Tooltip.ValueBoolPointer(),
			Translatable:         item.Translatable.ValueBoolPointer(),
			UseUuid:              item.UseUuid.ValueBoolPointer(),
		}
	}

	// Sort the fields by position. Storyblok has a position field but ends up
	// using the ordering of the json...
	schema := sortComponentFields(raw)

	return sbmgmt.ComponentInput{
		Name:               m.Name.ValueString(),
		Schema:             schema,
		IsNestable:         m.IsNestable.ValueBoolPointer(),
		IsRoot:             m.IsRoot.ValueBoolPointer(),
		ComponentGroupUuid: ref(asUUID(m.ComponentGroupUUID)),
	}
}

func (m *componentResourceModel) fromRemote(spaceID int64, c *sbmgmt.Component) error {
	if c == nil {
		return fmt.Errorf("component is nil")
	}
	m.ID = types.StringValue(createIdentifier(spaceID, c.Id))
	m.ComponentID = types.Int64Value(c.Id)
	m.CreatedAt = types.StringValue(c.CreatedAt.String())
	m.IsRoot = types.BoolPointerValue(c.IsRoot)
	m.IsNestable = types.BoolPointerValue(c.IsNestable)
	m.ComponentGroupUUID = fromUUID(c.ComponentGroupUuid)

	schema := make(map[string]fieldModel, c.Schema.Len())
	for pair := c.Schema.Oldest(); pair != nil; pair = pair.Next() {
		name := pair.Key
		field := pair.Value

		schema[name] = fieldModel{
			Type:     types.StringValue(field.Type),
			Position: types.Int64Value(field.Pos),

			AddHttps:             types.BoolPointerValue(field.AddHttps),
			AssetFolderId:        types.Int64PointerValue(field.AssetFolderId),
			CanSync:              types.BoolPointerValue(field.CanSync),
			ComponentWhitelist:   convertToStringSlice(field.ComponentWhitelist),
			DatasourceSlug:       types.StringPointerValue(field.DatasourceSlug),
			DefaultValue:         types.StringPointerValue(field.DefaultValue),
			Description:          types.StringPointerValue(field.Description),
			DisableTime:          types.BoolPointerValue(field.DisableTime),
			DisplayName:          types.StringPointerValue(field.DisplayName),
			ExternalDatasource:   types.StringPointerValue(field.ExternalDatasource),
			FieldType:            types.StringPointerValue(field.FieldType),
			Filetypes:            convertToStringSlice(field.Filetypes),
			FolderSlug:           types.StringPointerValue(field.FolderSlug),
			ImageCrop:            types.BoolPointerValue(field.ImageCrop),
			ImageHeight:          types.StringPointerValue(field.ImageHeight),
			ImageWidth:           types.StringPointerValue(field.ImageWidth),
			KeepImageSize:        types.BoolPointerValue(field.KeepImageSize),
			Keys:                 convertToStringSlice(field.Keys),
			MaxLength:            types.Int64PointerValue(field.MaxLength),
			Maximum:              types.Int64PointerValue(field.Maximum),
			NoTranslate:          types.BoolPointerValue(field.NoTranslate),
			Options:              serializeOptionsModel(field.Options),
			PreviewField:         types.BoolPointerValue(field.PreviewField),
			Regex:                types.StringPointerValue(field.Regex),
			Required:             types.BoolPointerValue(field.Required),
			RestrictComponents:   types.BoolPointerValue(field.RestrictComponents),
			RestrictContentTypes: types.BoolPointerValue(field.RestrictContentTypes),
			RichMarkdown:         types.BoolPointerValue(field.RichMarkdown),
			Rtl:                  types.BoolPointerValue(field.Rtl),
			Source:               types.StringPointerValue(field.Source),
			Tooltip:              types.BoolPointerValue(field.Tooltip),
			Translatable:         types.BoolPointerValue(field.Translatable),
			UseUuid:              types.BoolPointerValue(field.UseUuid),
		}
	}
	m.Schema = schema
	return nil
}

func getComponentTypes() map[string]string {
	return map[string]string{
		"bloks":      "Blocks: a field to interleave other components in your current one",
		"text":       "Text: a text field",
		"textarea":   "Textarea: a text area",
		"markdown":   "Markdown: write markdown with a text area and additional formatting options",
		"number":     "Number: a number field",
		"datetime":   "Date/Time: a date- and time picker",
		"boolean":    "Boolean: a checkbox - true/false",
		"options":    "Multi-Options: a list of checkboxes",
		"option":     "Single-Option: a single dropdown",
		"asset":      "Asset: Single asset (images, videos, audio, and documents)",
		"multiasset": "Multi-Assets: (images, videos, audio, and documents)",
		"multilink":  "Link: an input field for internal linking to other stories",
		"section":    "Group: no input possibility - allows you to group fields in sections",
		"custom":     "Plugin: Extend the editor yourself with a color picker or similar - Check out: Creating a Storyblok field type plugin",
		"image":      "Image (old): a upload field for a single image with cropping possibilities",
		"file":       "File (old): a upload field for a single file",
	}
}

func serializeOptionsModel(options *[]sbmgmt.FieldOption) []optionModel {
	if options == nil {
		return nil
	}

	optionModels := make([]optionModel, len(*options))
	for i, option := range *options {
		optionModels[i] = optionModel{
			Name:  types.StringValue(option.Name),
			Value: types.StringValue(option.Value),
		}
	}

	return optionModels
}

func deserializeOptionsModel(options []optionModel) *[]sbmgmt.FieldOption {
	if options == nil {
		return nil
	}

	optionModels := make([]sbmgmt.FieldOption, len(options))
	for i, option := range options {
		optionModels[i] = sbmgmt.FieldOption{
			Name:  option.Name.ValueString(),
			Value: option.Value.ValueString(),
		}
	}

	return &optionModels
}
