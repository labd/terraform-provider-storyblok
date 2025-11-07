package utils

import (
	"bytes"
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	"github.com/elliotchance/pie/v2"
	"github.com/gofrs/uuid"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/labd/storyblok-go-sdk/sbmgmt"
	orderedmap "github.com/wk8/go-ordered-map/v2"
)

func SortComponentFields(input map[string]sbmgmt.FieldInput) *orderedmap.OrderedMap[string, sbmgmt.FieldInput] {
	type Pair struct {
		Key   string
		Value *sbmgmt.FieldInput
	}

	values := make([]Pair, 0, len(input))
	for key := range input {
		value := input[key]
		values = append(values, Pair{
			Key:   key,
			Value: &value,
		})
	}

	sorted := pie.SortUsing(values, func(a, b Pair) bool {
		return a.Value.Pos < b.Value.Pos
	})

	result := orderedmap.New[string, sbmgmt.FieldInput]()
	for _, item := range sorted {
		result.Set(item.Key, *item.Value)
	}

	return result
}

func AsUUIDPointer(s basetypes.StringValue) *uuid.UUID {
	var componentGroupUuid *uuid.UUID
	if !s.IsNull() {
		value := uuid.Must(uuid.FromString(s.ValueString()))

		componentGroupUuid = &value
	}
	return componentGroupUuid
}

func FromUUID(v *uuid.UUID) types.String {
	if v == nil || v.IsNil() {
		return types.StringPointerValue(nil)
	}
	return types.StringValue(v.String())
}

func FromStringPointer(v *string) types.String {
	if v == nil {
		return types.StringPointerValue(nil)
	}
	return types.StringValue(*v)
}

func Must[T any](v T, err error) T {
	if err != nil {
		panic(err)
	}
	return v
}

func HCLTemplate(data string, params map[string]any) string {
	var out bytes.Buffer
	tmpl := template.Must(template.New("hcl").Parse(data))
	err := tmpl.Execute(&out, params)
	if err != nil {
		panic(err)
	}
	return out.String()
}

func CleanHeaders(headers http.Header, keep ...string) http.Header {
	for key := range headers {
		if !contains(keep, key) {
			headers.Del(key)
		}
	}
	return headers
}

func Int64ToStringInterfacePointer(options types.Int64) *interface{} {
	if options.IsNull() || options.IsUnknown() {
		return nil
	}
	str := strconv.FormatInt(options.ValueInt64(), 10)
	var v interface{} = str
	return &v
}

func InterfacePointerToInt64(options *interface{}) (types.Int64, error) {
	if options == nil {
		return types.Int64PointerValue(nil), nil
	}
	switch v := (*options).(type) {
	case int64:
		return types.Int64Value(v), nil
	case int:
		return types.Int64Value(int64(v)), nil
	case float64:
		return types.Int64Value(int64(v)), nil
	case string:
		parsed, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return types.Int64PointerValue(nil), fmt.Errorf("cannot convert string to int64: %v", err)
		}
		return types.Int64Value(parsed), nil
	default:
		return types.Int64PointerValue(nil), fmt.Errorf("unsupported type: %T", v)
	}
}

func contains(keep []string, key string) bool {
	for _, k := range keep {
		if k == key {
			return true
		}
	}
	return false
}
