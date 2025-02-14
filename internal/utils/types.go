package utils

import (
	"github.com/elliotchance/pie/v2"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func ConvertToStringSlice(strSlicePtr *[]string) []types.String {
	if strSlicePtr == nil {
		return nil
	}

	strSlice := make([]types.String, len(*strSlicePtr))
	for i, str := range *strSlicePtr {
		strSlice[i] = types.StringValue(str)
	}

	return strSlice
}

func ConvertToInt64Slice(intSlicePtr *[]int) []types.Int64 {
	if intSlicePtr == nil {
		return nil
	}

	intSlice := make([]types.Int64, len(*intSlicePtr))
	for i, str := range *intSlicePtr {
		intSlice[i] = types.Int64Value(int64(str))
	}

	return intSlice
}

func ConvertToPointerStringSlice(slice []types.String) *[]string {
	if slice == nil {
		return nil
	}

	result := pie.Map(slice, func(s types.String) string {
		return s.ValueString()
	})

	return &result
}

func ConvertToPointerIntSlice(slice []types.Int64) *[]int {
	if slice == nil {
		return nil
	}

	result := pie.Map(slice, func(s types.Int64) int {
		return int(s.ValueInt64())
	})

	return &result
}
