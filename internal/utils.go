package internal

import (
	"fmt"

	"github.com/elliotchance/pie/v2"
	"github.com/labd/storyblok-go-sdk/sbmgmt"
	orderedmap "github.com/wk8/go-ordered-map/v2"
)

func ref[T any](s T) *T {
	return &s
}

func createIdentifier(spaceId int64, componentId int64) string {
	return fmt.Sprintf("%d/%d", spaceId, componentId)
}

func parseIdentifier(identifier string) (spaceId int64, componentId int64) {
	fmt.Sscanf(identifier, "%d/%d", &spaceId, &componentId)
	return
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
