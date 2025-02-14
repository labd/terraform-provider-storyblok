package internal

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/labd/storyblok-go-sdk/sbmgmt"

	"github.com/labd/terraform-provider-storyblok/internal/utils"
)

// componentResourceModel maps the resource schema data.
type componentResourceModel struct {
	ID                 types.String          `tfsdk:"id"`
	ComponentID        types.Int64           `tfsdk:"component_id"`
	SpaceID            types.Int64           `tfsdk:"space_id"`
	CreatedAt          types.String          `tfsdk:"created_at"`
	DisplayName        types.String          `tfsdk:"display_name"`
	Color              types.String          `tfsdk:"color"`
	Icon               types.String          `tfsdk:"icon"`
	Image              types.String          `tfsdk:"image"`
	PreviewTmpl        types.String          `tfsdk:"preview_tmpl"`
	PreviewField       types.String          `tfsdk:"preview_field"`
	Name               types.String          `tfsdk:"name"`
	IsRoot             types.Bool            `tfsdk:"is_root"`
	IsNestable         types.Bool            `tfsdk:"is_nestable"`
	ComponentGroupUUID types.String          `tfsdk:"component_group_uuid"`
	Schema             map[string]fieldModel `tfsdk:"schema"`
}

type fieldModel struct {
	Type     types.String `tfsdk:"type"`
	Position types.Int64  `tfsdk:"position"`

	AddHttps                types.Bool                 `tfsdk:"add_https"`
	AllowTargetBlank        types.Bool                 `tfsdk:"allow_target_blank"`
	AssetFolderId           types.Int64                `tfsdk:"asset_folder_id"`
	CanSync                 types.Bool                 `tfsdk:"can_sync"`
	ComponentWhitelist      []types.String             `tfsdk:"component_whitelist"`
	ConditionalSettings     []conditionalSettingsModel `tfsdk:"conditional_settings"`
	CustomizeToolbar        types.Bool                 `tfsdk:"customize_toolbar"`
	DatasourceSlug          types.String               `tfsdk:"datasource_slug"`
	DefaultValue            types.String               `tfsdk:"default_value"`
	Description             types.String               `tfsdk:"description"`
	DisableTime             types.Bool                 `tfsdk:"disable_time"`
	DisplayName             types.String               `tfsdk:"display_name"`
	ExternalDatasource      types.String               `tfsdk:"external_datasource"`
	FieldType               types.String               `tfsdk:"field_type"`
	Filetypes               []types.String             `tfsdk:"filetypes"`
	FilterContentType       []types.String             `tfsdk:"filter_content_type"`
	FolderSlug              types.String               `tfsdk:"folder_slug"`
	ForceLinkScope          types.Bool                 `tfsdk:"force_link_scope"`
	ImageCrop               types.Bool                 `tfsdk:"image_crop"`
	ImageHeight             types.String               `tfsdk:"image_height"`
	ImageWidth              types.String               `tfsdk:"image_width"`
	KeepImageSize           types.Bool                 `tfsdk:"keep_image_size"`
	Keys                    []types.String             `tfsdk:"keys"`
	LinkScope               types.String               `tfsdk:"link_scope"`
	Maximum                 types.Int64                `tfsdk:"maximum"`
	MaxLength               types.Int64                `tfsdk:"max_length"`
	MaxOptions              types.Int64                `tfsdk:"max_options"`
	Minimum                 types.Int64                `tfsdk:"minimum"`
	MinOptions              types.Int64                `tfsdk:"min_options"`
	NoTranslate             types.Bool                 `tfsdk:"no_translate"`
	Options                 []optionModel              `tfsdk:"options"`
	Regex                   types.String               `tfsdk:"regex"`
	Required                types.Bool                 `tfsdk:"required"`
	RestrictComponents      types.Bool                 `tfsdk:"restrict_components"`
	RestrictContentTypes    types.Bool                 `tfsdk:"restrict_content_types"`
	RichMarkdown            types.Bool                 `tfsdk:"rich_markdown"`
	Rtl                     types.Bool                 `tfsdk:"rtl"`
	Source                  types.String               `tfsdk:"source"`
	Toolbar                 []types.String             `tfsdk:"toolbar"`
	Tooltip                 types.Bool                 `tfsdk:"tooltip"`
	Translatable            types.Bool                 `tfsdk:"translatable"`
	UseUuid                 types.Bool                 `tfsdk:"use_uuid"`
	Steps                   types.Int64                `tfsdk:"steps"`
	ShowAnchor              types.Bool                 `tfsdk:"show_anchor"`
	RestrictType            types.String               `tfsdk:"restrict_type"`
	ComponentGroupWhitelist []types.String             `tfsdk:"component_group_whitelist"`
	ComponentTagWhitelist   []types.Int64              `tfsdk:"component_tag_whitelist"`
	AssetLinkType           types.Bool                 `tfsdk:"asset_link_type"`
	AllowAdvancedSearch     types.Bool                 `tfsdk:"allow_advanced_search"`
	AllowCustomAttributes   types.Bool                 `tfsdk:"allow_custom_attributes"`
	AllowExternalUrl        types.Bool                 `tfsdk:"allow_external_url"`
	AllowMultiline          types.Bool                 `tfsdk:"allow_multiline"`
	Decimals                types.Int64                `tfsdk:"decimals"`
	EmailLinkType           types.Bool                 `tfsdk:"email_link_type"`
	EntryAppearance         types.String               `tfsdk:"entry_appearance"`
	ExcludeEmptyOption      types.Bool                 `tfsdk:"exclude_empty_option"`
	ExcludeFromMerge        types.Bool                 `tfsdk:"exclude_from_merge"`
	ExcludeFromOverwrite    types.Bool                 `tfsdk:"exclude_from_overwrite"`
	ForceMerge              types.Bool                 `tfsdk:"force_merge"`
	InlineLabel             types.Bool                 `tfsdk:"inline_label"`
	IsReferenceType         types.Bool                 `tfsdk:"is_reference_type"`
	MaxValue                types.Int64                `tfsdk:"max_value"`
	MinValue                types.Int64                `tfsdk:"min_value"`
}

type conditionalSettingsModel struct {
	Modifications  []modificationModel  `tfsdk:"modifications"`
	RuleMatch      types.String         `tfsdk:"rule_match"`
	RuleConditions []ruleConditionModel `tfsdk:"rule_conditions"`
}

type ruleConditionModel struct {
	Validation      types.String         `tfsdk:"validation"`
	Value           types.String         `tfsdk:"value"`
	ValidatedObject validatedObjectModel `tfsdk:"validated_object"`
}

type validatedObjectModel struct {
	FieldKey types.String `tfsdk:"field_key"`
}

type modificationModel struct {
	Display  types.String `tfsdk:"display"`
	Required types.Bool   `tfsdk:"required"`
}

type optionModel struct {
	Name  types.String `tfsdk:"name"`
	Value types.String `tfsdk:"value"`
}

func (m *componentResourceModel) toRemoteInput() sbmgmt.ComponentCreateInput {

	raw := make(map[string]sbmgmt.FieldInput, len(m.Schema))
	for name := range m.Schema {
		item := m.Schema[name]
		raw[name] = toFieldInput(item)
	}

	// Sort the fields by position. Storyblok has a position field but ends up
	// using the ordering of the json...
	schema := sortComponentFields(raw)

	componentGroupUuid := asUUIDPointer(m.ComponentGroupUUID)

	return sbmgmt.ComponentCreateInput{
		Component: sbmgmt.ComponentBase{
			Color:              m.Color.ValueStringPointer(),
			ComponentGroupUuid: componentGroupUuid,
			DisplayName:        m.DisplayName.ValueStringPointer(),
			Icon:               (*sbmgmt.ComponentBaseIcon)(m.Icon.ValueStringPointer()),
			Image:              m.Image.ValueStringPointer(),
			IsNestable:         m.IsNestable.ValueBoolPointer(),
			IsRoot:             m.IsRoot.ValueBoolPointer(),
			Name:               m.Name.ValueString(),
			PreviewTmpl:        m.PreviewTmpl.ValueStringPointer(),
			PreviewField:       m.PreviewField.ValueStringPointer(),
			Schema:             schema,
		},
	}
}
func (m *componentResourceModel) toUpdateInput() sbmgmt.ComponentUpdateInput {

	raw := make(map[string]sbmgmt.FieldInput, len(m.Schema))
	for name := range m.Schema {
		item := m.Schema[name]
		raw[name] = toFieldInput(item)
	}

	// Sort the fields by position. Storyblok has a position field but ends up
	// using the ordering of the json...
	schema := sortComponentFields(raw)

	componentGroupUuid := asUUIDPointer(m.ComponentGroupUUID)

	return sbmgmt.ComponentUpdateInput{
		Component: sbmgmt.ComponentBase{
			Color:              m.Color.ValueStringPointer(),
			ComponentGroupUuid: componentGroupUuid,
			DisplayName:        m.DisplayName.ValueStringPointer(),
			Icon:               (*sbmgmt.ComponentBaseIcon)(m.Icon.ValueStringPointer()),
			Image:              m.Image.ValueStringPointer(),
			IsNestable:         m.IsNestable.ValueBoolPointer(),
			IsRoot:             m.IsRoot.ValueBoolPointer(),
			Name:               m.Name.ValueString(),
			PreviewTmpl:        m.PreviewTmpl.ValueStringPointer(),
			PreviewField:       m.PreviewField.ValueStringPointer(),
			Schema:             schema,
		},
	}
}

func toFieldInput(item fieldModel) sbmgmt.FieldInput {
	return sbmgmt.FieldInput{
		Type: item.Type.ValueString(),
		Pos:  item.Position.ValueInt64(),

		AddHttps:                item.AddHttps.ValueBoolPointer(),
		AllowTargetBlank:        item.AllowTargetBlank.ValueBoolPointer(),
		AssetFolderId:           item.AssetFolderId.ValueInt64Pointer(),
		CanSync:                 item.CanSync.ValueBoolPointer(),
		ComponentWhitelist:      utils.ConvertToPointerStringSlice(item.ComponentWhitelist),
		ConditionalSettings:     deserializeConditionalSettings(item.ConditionalSettings),
		CustomizeToolbar:        item.CustomizeToolbar.ValueBoolPointer(),
		DatasourceSlug:          item.DatasourceSlug.ValueStringPointer(),
		DefaultValue:            item.DefaultValue.ValueStringPointer(),
		Description:             item.Description.ValueStringPointer(),
		DisableTime:             item.DisableTime.ValueBoolPointer(),
		DisplayName:             item.DisplayName.ValueStringPointer(),
		ExternalDatasource:      item.ExternalDatasource.ValueStringPointer(),
		FieldType:               item.FieldType.ValueStringPointer(),
		Filetypes:               utils.ConvertToPointerStringSlice(item.Filetypes),
		FilterContentType:       utils.ConvertToPointerStringSlice(item.FilterContentType),
		FolderSlug:              item.FolderSlug.ValueStringPointer(),
		ForceLinkScope:          item.ForceLinkScope.ValueBoolPointer(),
		ImageCrop:               item.ImageCrop.ValueBoolPointer(),
		ImageHeight:             item.ImageHeight.ValueStringPointer(),
		ImageWidth:              item.ImageWidth.ValueStringPointer(),
		KeepImageSize:           item.KeepImageSize.ValueBoolPointer(),
		Keys:                    utils.ConvertToPointerStringSlice(item.Keys),
		LinkScope:               item.LinkScope.ValueStringPointer(),
		Maximum:                 item.Maximum.ValueInt64Pointer(),
		MaxLength:               item.MaxLength.ValueInt64Pointer(),
		MaxOptions:              item.MaxOptions.ValueInt64Pointer(),
		Minimum:                 item.Minimum.ValueInt64Pointer(),
		MinOptions:              item.MinOptions.ValueInt64Pointer(),
		NoTranslate:             item.NoTranslate.ValueBoolPointer(),
		Options:                 deserializeOptionsModel(item.Options),
		Regex:                   item.Regex.ValueStringPointer(),
		Required:                item.Required.ValueBoolPointer(),
		RestrictComponents:      item.RestrictComponents.ValueBoolPointer(),
		RestrictContentTypes:    item.RestrictContentTypes.ValueBoolPointer(),
		RichMarkdown:            item.RichMarkdown.ValueBoolPointer(),
		Rtl:                     item.Rtl.ValueBoolPointer(),
		Source:                  item.Source.ValueStringPointer(),
		Toolbar:                 utils.ConvertToPointerStringSlice(item.Toolbar),
		Tooltip:                 item.Tooltip.ValueBoolPointer(),
		Translatable:            item.Translatable.ValueBoolPointer(),
		UseUuid:                 item.UseUuid.ValueBoolPointer(),
		Steps:                   item.Steps.ValueInt64Pointer(),
		ShowAnchor:              item.ShowAnchor.ValueBoolPointer(),
		RestrictType:            item.RestrictType.ValueStringPointer(),
		ComponentGroupWhitelist: utils.ConvertToPointerStringSlice(item.ComponentGroupWhitelist),
		ComponentTagWhitelist:   utils.ConvertToPointerIntSlice(item.ComponentTagWhitelist),
		AssetLinkType:           item.AssetLinkType.ValueBoolPointer(),
		AllowAdvancedSearch:     item.AllowAdvancedSearch.ValueBoolPointer(),
		AllowCustomAttributes:   item.AllowCustomAttributes.ValueBoolPointer(),
		AllowExternalUrl:        item.AllowExternalUrl.ValueBoolPointer(),
		AllowMultiline:          item.AllowMultiline.ValueBoolPointer(),
		Decimals:                item.Decimals.ValueInt64Pointer(),
		EmailLinkType:           item.EmailLinkType.ValueBoolPointer(),
		EntryAppearance:         item.EntryAppearance.ValueStringPointer(),
		ExcludeEmptyOption:      item.ExcludeEmptyOption.ValueBoolPointer(),
		ExcludeFromMerge:        item.ExcludeFromMerge.ValueBoolPointer(),
		ExcludeFromOverwrite:    item.ExcludeFromOverwrite.ValueBoolPointer(),
		ForceMerge:              item.ForceMerge.ValueBoolPointer(),
		InlineLabel:             item.InlineLabel.ValueBoolPointer(),
		IsReferenceType:         item.IsReferenceType.ValueBoolPointer(),
		MaxValue:                item.MaxValue.ValueInt64Pointer(),
		MinValue:                item.MinValue.ValueInt64Pointer(),
	}
}

func (m *componentResourceModel) fromRemote(spaceID int64, c *sbmgmt.Component) error {
	if c == nil {
		return fmt.Errorf("component is nil")
	}
	m.ID = types.StringValue(utils.CreateIdentifier(spaceID, c.Id))
	m.ComponentID = types.Int64Value(c.Id)
	m.CreatedAt = types.StringValue(c.CreatedAt.String())
	m.IsRoot = types.BoolPointerValue(c.IsRoot)
	m.IsNestable = types.BoolPointerValue(c.IsNestable)
	m.ComponentGroupUUID = fromUUID(c.ComponentGroupUuid)
	m.Color = fromStringPointer(c.Color)
	m.DisplayName = fromStringPointer(c.DisplayName)
	m.Image = fromStringPointer(c.Image)
	m.PreviewField = fromStringPointer(c.PreviewField)
	m.PreviewTmpl = fromStringPointer(c.PreviewTmpl)
	if c.Icon != nil {
		m.Icon = types.StringValue(string(*c.Icon))
	}

	schema := make(map[string]fieldModel, c.Schema.Len())
	for pair := c.Schema.Oldest(); pair != nil; pair = pair.Next() {
		name := pair.Key
		field := pair.Value

		schema[name] = toFieldModel(field)
	}
	m.Schema = schema
	return nil
}

func toFieldModel(field sbmgmt.FieldInput) fieldModel {
	return fieldModel{
		Type:     types.StringValue(field.Type),
		Position: types.Int64Value(field.Pos),

		AddHttps:                types.BoolPointerValue(field.AddHttps),
		AllowTargetBlank:        types.BoolPointerValue(field.AllowTargetBlank),
		AssetFolderId:           types.Int64PointerValue(field.AssetFolderId),
		CanSync:                 types.BoolPointerValue(field.CanSync),
		ComponentWhitelist:      utils.ConvertToStringSlice(field.ComponentWhitelist),
		CustomizeToolbar:        types.BoolPointerValue(field.CustomizeToolbar),
		DatasourceSlug:          types.StringPointerValue(field.DatasourceSlug),
		DefaultValue:            types.StringPointerValue(field.DefaultValue),
		Description:             types.StringPointerValue(field.Description),
		DisableTime:             types.BoolPointerValue(field.DisableTime),
		DisplayName:             types.StringPointerValue(field.DisplayName),
		ExternalDatasource:      types.StringPointerValue(field.ExternalDatasource),
		FieldType:               types.StringPointerValue(field.FieldType),
		Filetypes:               utils.ConvertToStringSlice(field.Filetypes),
		FilterContentType:       utils.ConvertToStringSlice(field.FilterContentType),
		FolderSlug:              types.StringPointerValue(field.FolderSlug),
		ForceLinkScope:          types.BoolPointerValue(field.ForceLinkScope),
		ImageCrop:               types.BoolPointerValue(field.ImageCrop),
		ImageHeight:             types.StringPointerValue(field.ImageHeight),
		ImageWidth:              types.StringPointerValue(field.ImageWidth),
		KeepImageSize:           types.BoolPointerValue(field.KeepImageSize),
		Keys:                    utils.ConvertToStringSlice(field.Keys),
		LinkScope:               types.StringPointerValue(field.LinkScope),
		Maximum:                 types.Int64PointerValue(field.Maximum),
		MaxLength:               types.Int64PointerValue(field.MaxLength),
		MaxOptions:              types.Int64PointerValue(field.MaxOptions),
		Minimum:                 types.Int64PointerValue(field.Minimum),
		MinOptions:              types.Int64PointerValue(field.MinOptions),
		NoTranslate:             types.BoolPointerValue(field.NoTranslate),
		Options:                 serializeOptionsModel(field.Options),
		Regex:                   types.StringPointerValue(field.Regex),
		Required:                types.BoolPointerValue(field.Required),
		RestrictComponents:      types.BoolPointerValue(field.RestrictComponents),
		RestrictContentTypes:    types.BoolPointerValue(field.RestrictContentTypes),
		RichMarkdown:            types.BoolPointerValue(field.RichMarkdown),
		Rtl:                     types.BoolPointerValue(field.Rtl),
		Source:                  types.StringPointerValue(field.Source),
		Toolbar:                 utils.ConvertToStringSlice(field.Toolbar),
		Tooltip:                 types.BoolPointerValue(field.Tooltip),
		Translatable:            types.BoolPointerValue(field.Translatable),
		UseUuid:                 types.BoolPointerValue(field.UseUuid),
		ConditionalSettings:     serializeConditionalSettings(field.ConditionalSettings),
		Steps:                   types.Int64PointerValue(field.Steps),
		ShowAnchor:              types.BoolPointerValue(field.ShowAnchor),
		RestrictType:            types.StringPointerValue(field.RestrictType),
		ComponentGroupWhitelist: utils.ConvertToStringSlice(field.ComponentGroupWhitelist),
		ComponentTagWhitelist:   utils.ConvertToInt64Slice(field.ComponentTagWhitelist),
		AssetLinkType:           types.BoolPointerValue(field.AssetLinkType),
		AllowAdvancedSearch:     types.BoolPointerValue(field.AllowAdvancedSearch),
		AllowCustomAttributes:   types.BoolPointerValue(field.AllowCustomAttributes),
		AllowExternalUrl:        types.BoolPointerValue(field.AllowExternalUrl),
		AllowMultiline:          types.BoolPointerValue(field.AllowMultiline),
		Decimals:                types.Int64PointerValue(field.Decimals),
		EmailLinkType:           types.BoolPointerValue(field.EmailLinkType),
		EntryAppearance:         types.StringPointerValue(field.EntryAppearance),
		ExcludeEmptyOption:      types.BoolPointerValue(field.ExcludeEmptyOption),
		ExcludeFromMerge:        types.BoolPointerValue(field.ExcludeFromMerge),
		ExcludeFromOverwrite:    types.BoolPointerValue(field.ExcludeFromOverwrite),
		ForceMerge:              types.BoolPointerValue(field.ForceMerge),
		InlineLabel:             types.BoolPointerValue(field.InlineLabel),
		IsReferenceType:         types.BoolPointerValue(field.IsReferenceType),
		MaxValue:                types.Int64PointerValue(field.MaxValue),
		MinValue:                types.Int64PointerValue(field.MinValue),
	}
}

func getComponentTypes() map[string]string {
	return map[string]string{
		"bloks":      "Blocks: a field to interleave other components in your current one",
		"text":       "Text: a text field",
		"textarea":   "Textarea: a text area",
		"markdown":   "Markdown: write markdown with a text area and additional formatting options",
		"richtext":   "Richtext: write richtext with a text area and additional formatting options",
		"number":     "Number: a number field",
		"datetime":   "Date/Time: a date- and time picker",
		"boolean":    "Boolean: a checkbox - true/false",
		"options":    "Multi-Options: a list of checkboxes",
		"option":     "Single-Option: a single dropdown",
		"asset":      "Asset: Single asset (images, videos, audio, and documents)",
		"multiasset": "Multi-Assets: (images, videos, audio, and documents)",
		"multilink":  "Link: an input field for internal linking to other stories",
		"section":    "Group: no input possibility - allows you to group fields in sections",
		"tab":        "Tab: no input possibility - allows you to group fields in tabs",
		"custom":     "Plugin: Extend the editor yourself with a color picker or similar - Check out: Creating a Storyblok field type plugin",
		"image":      "Image (old): a upload field for a single image with cropping possibilities",
		"file":       "File (old): a upload field for a single file",
		"table":      "Table: a table field",
	}
}

func getComponentIcons() []string {
	return []string{
		"block-@",
		"block-1-2block",
		"block-add",
		"block-arrow-pointer",
		"block-block",
		"block-buildin",
		"block-cart",
		"block-center-m",
		"block-comment",
		"block-doc",
		"block-dollar-sign",
		"block-email",
		"block-image",
		"block-keyboard",
		"block-locked",
		"block-map-pin",
		"block-mobile",
		"block-monitor",
		"block-paycard",
		"block-resize-fc",
		"block-share",
		"block-shield",
		"block-shield-2",
		"block-sticker",
		"block-suitcase",
		"block-table",
		"block-table-2",
		"block-tag",
		"block-text-c",
		"block-text-img-c",
		"block-text-img-l",
		"block-text-img-r",
		"block-text-img-r-l",
		"block-text-img-t-l",
		"block-text-img-t-r",
		"block-text-l",
		"block-text-r",
		"block-unlocked",
		"block-wallet",
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

func serializeConditionalSettings(conditionalSettings *[]sbmgmt.ConditionalSettings) []conditionalSettingsModel {
	if conditionalSettings == nil {
		return nil
	}

	serializedConditionalSettings := make([]conditionalSettingsModel, len(*conditionalSettings))

	for i, conditionalSetting := range *conditionalSettings {
		serializedConditionalSettings[i] = conditionalSettingsModel{
			RuleMatch:      types.StringPointerValue((*string)(conditionalSetting.RuleMatch)),
			Modifications:  serializeModifications(conditionalSetting.Modifications),
			RuleConditions: serializeRuleConditions(conditionalSetting.RuleConditions),
		}
	}

	return serializedConditionalSettings
}

func serializeModifications(modifications *[]sbmgmt.Modification) []modificationModel {
	if modifications == nil {
		return nil
	}

	serializedModifications := make([]modificationModel, len(*modifications))

	for i, modification := range *modifications {
		serializedModifications[i] = modificationModel{
			Required: types.BoolPointerValue(modification.Required),
			Display:  types.StringPointerValue((*string)(modification.Display)),
		}
	}

	return serializedModifications
}
func serializeRuleConditions(ruleConditions *[]sbmgmt.RuleCondition) []ruleConditionModel {

	if ruleConditions == nil {
		return nil
	}

	serializedRuleConditions := make([]ruleConditionModel, len(*ruleConditions))

	for i, ruleCondition := range *ruleConditions {
		serializedRuleConditions[i] = ruleConditionModel{
			Value:      types.StringPointerValue(ruleCondition.Value),
			Validation: types.StringPointerValue((*string)(ruleCondition.Validation)),
			ValidatedObject: validatedObjectModel{
				FieldKey: types.StringPointerValue(ruleCondition.ValidatedObject.FieldKey),
			},
		}
	}

	return serializedRuleConditions
}

func deserializeConditionalSettings(conditionalSettings []conditionalSettingsModel) *[]sbmgmt.ConditionalSettings {
	if conditionalSettings == nil {
		return nil
	}

	deserializedConditionalSettings := make([]sbmgmt.ConditionalSettings, len(conditionalSettings))

	for i, setting := range conditionalSettings {

		deserializedConditionalSettings[i] = sbmgmt.ConditionalSettings{
			RuleMatch:      (*sbmgmt.ConditionalSettingsRuleMatch)(setting.RuleMatch.ValueStringPointer()),
			Modifications:  deserializeConditionalSettingsModificationsModel(setting.Modifications),
			RuleConditions: deserializeRuleConditions(setting.RuleConditions),
		}
	}

	return &deserializedConditionalSettings
}

func deserializeConditionalSettingsModificationsModel(conditionalSettingsModifications []modificationModel) *[]sbmgmt.Modification {
	deserializedModifications := make([]sbmgmt.Modification, len(conditionalSettingsModifications))

	for i, modification := range conditionalSettingsModifications {
		deserializedModifications[i] = sbmgmt.Modification{
			Required: modification.Required.ValueBoolPointer(),
			Display:  (*sbmgmt.ModificationDisplay)(modification.Display.ValueStringPointer()),
		}
	}

	return &deserializedModifications
}

func deserializeRuleConditions(ruleConditions []ruleConditionModel) *[]sbmgmt.RuleCondition {
	deserializedRuleConditions := make([]sbmgmt.RuleCondition, len(ruleConditions))

	validatedObjectFieldAttrValue := sbmgmt.ValidatedObjectFieldAttr("value")
	validatedValidatedObjectType := sbmgmt.ValidatedObjectType("field")

	for i, ruleCondition := range ruleConditions {
		deserializedRuleConditions[i] = sbmgmt.RuleCondition{
			Validation: (*sbmgmt.RuleConditionValidation)(ruleCondition.Validation.ValueStringPointer()),
			Value:      ruleCondition.Value.ValueStringPointer(),
			ValidatedObject: &sbmgmt.ValidatedObject{
				FieldKey:  ruleCondition.ValidatedObject.FieldKey.ValueStringPointer(),
				FieldAttr: &validatedObjectFieldAttrValue,
				Type:      &validatedValidatedObjectType,
			},
		}
	}

	return &deserializedRuleConditions
}
