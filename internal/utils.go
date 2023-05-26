package internal

import "fmt"

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
