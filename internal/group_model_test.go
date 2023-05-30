package internal

import (
	"testing"
	"time"

	"github.com/gofrs/uuid"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/labd/storyblok-go-sdk/sbmgmt"
	"github.com/stretchr/testify/assert"
)

func TestComponentGroupResourceModel_ToRemoteInput(t *testing.T) {
	model := &componentGroupResourceModel{
		Name: types.StringValue("TestGroup"),
	}

	expectedInput := sbmgmt.ComponentGroupInput{
		Name: "TestGroup",
	}

	actualInput := model.toRemoteInput()

	assert.Equal(t, expectedInput, actualInput, "Converted input does not match expected input")
}

func TestComponentGroupResourceModel_FromRemote(t *testing.T) {
	spaceID := int64(123)
	groupID := int64(456)
	name := "TestGroup"

	componentGroup := &sbmgmt.ComponentGroup{
		Id:        groupID,
		Name:      name,
		Uuid:      must(uuid.FromString("ebd1af2e-875f-47e5-8886-4d3baea94d99")),
		CreatedAt: must(time.Parse(time.RFC3339, "2023-05-26T12:34:56Z")),
		UpdatedAt: must(time.Parse(time.RFC3339, "2023-05-27T10:11:12Z")),
	}

	model := &componentGroupResourceModel{}
	err := model.fromRemote(spaceID, componentGroup)

	assert.NoError(t, err, "Error occurred during conversion")

	expectedModel := &componentGroupResourceModel{
		ID:        types.StringValue(createIdentifier(spaceID, groupID)),
		GroupID:   types.Int64Value(groupID),
		CreatedAt: types.StringValue("2023-05-26 12:34:56 +0000 UTC"),
		UpdatedAt: types.StringValue("2023-05-27 10:11:12 +0000 UTC"),
		UUID:      types.StringValue("ebd1af2e-875f-47e5-8886-4d3baea94d99"),
		Name:      types.StringValue(name),
	}

	assert.Equal(t, expectedModel, model, "Converted model does not match expected model")
}
