package utils

import (
	"fmt"
)

// CreateIdentifier creates a composite identifier from a space ID and an ID.
func CreateIdentifier(spaceId int64, id int64) string {
	return fmt.Sprintf("%d/%d", spaceId, id)
}

func ParseIdentifier(identifier string) (spaceId int64, id int64) {
	_, _ = fmt.Sscanf(identifier, "%d/%d", &spaceId, &id)
	return
}
