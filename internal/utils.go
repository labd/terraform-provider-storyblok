package internal

import (
	"fmt"

	"github.com/elliotchance/pie/v2"
	"github.com/gofrs/uuid"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/labd/storyblok-go-sdk/sbmgmt"
	orderedmap "github.com/wk8/go-ordered-map/v2"
)

func ref[T any](s T) *T {
	return &s
}

func createIdentifier(spaceId int64, id int64) string {
	return fmt.Sprintf("%d/%d", spaceId, id)
}

func parseIdentifier(identifier string) (spaceId int64, id int64) {
	fmt.Sscanf(identifier, "%d/%d", &spaceId, &id)
	return
}

func getClient(data any) sbmgmt.ClientWithResponsesInterface {
	c, ok := data.(sbmgmt.ClientWithResponsesInterface)
	if !ok {
		panic("invalid client type")
	}
	return c
}

func sortComponentFields(input map[string]sbmgmt.FieldInput) *orderedmap.OrderedMap[string, sbmgmt.FieldInput] {
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

func convertToStringSlice(strSlicePtr *[]string) []types.String {
	if strSlicePtr == nil {
		return nil
	}

	strSlice := make([]types.String, len(*strSlicePtr))
	for i, str := range *strSlicePtr {
		strSlice[i] = types.StringValue(str)
	}

	return strSlice
}

func convertToPointerStringSlice(slice []types.String) *[]string {
	if slice == nil {
		return nil
	}

	result := pie.Map(slice, func(s types.String) string {
		return s.ValueString()
	})

	return &result
}

func convertToPointerIntSlice(slice []types.Int64) *[]int {
	if slice == nil {
		return nil
	}

	result := pie.Map(slice, func(s types.Int64) int {
		return int(s.ValueInt64())
	})

	return &result
}

func asUUID(s types.String) uuid.UUID {
	if s.IsNull() {
		return uuid.Nil
	}
	val, err := uuid.FromString(s.ValueString())
	if err != nil {
		panic(err)
	}
	return val
}

func asUUIDPointer(s basetypes.StringValue) *uuid.UUID {
	var componentGroupUuid *uuid.UUID
	if !s.IsNull() {
		value := uuid.Must(uuid.FromString(s.ValueString()))

		componentGroupUuid = &value
	}
	return componentGroupUuid
}

func fromUUID(v *uuid.UUID) types.String {
	if v == nil || v.IsNil() {
		return types.StringPointerValue(nil)
	}
	return types.StringValue(v.String())
}

func fromStringPointer(v *string) types.String {
	if v == nil {
		return types.StringPointerValue(nil)
	}
	return types.StringValue(*v)
}

func must[T any](v T, err error) T {
	if err != nil {
		panic(err)
	}
	return v
}
