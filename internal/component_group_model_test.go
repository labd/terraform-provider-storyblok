package internal

import (
	"testing"

	"github.com/gofrs/uuid"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/labd/storyblok-go-sdk/sbmgmt"
	"github.com/stretchr/testify/assert"

	"github.com/labd/terraform-provider-storyblok/internal/utils"
)

func TestComponentGroupResourceModel_ToRemoteInput(t *testing.T) {
	model := &componentGroupResourceModel{
		Name: types.StringValue("TestGroup"),
	}

	expectedInput := sbmgmt.ComponentGroupCreateInput{
		ComponentGroup: sbmgmt.ComponentGroupBase{
			Name: "TestGroup",
		},
	}

	actualInput := model.toCreateInput()

	assert.Equal(t, expectedInput, actualInput, "Converted input does not match expected input")
}

func TestComponentGroupResourceModel_FromRemote(t *testing.T) {
	spaceID := int64(123)
	groupID := int64(456)
	name := "TestGroup"

	componentGroup := &sbmgmt.ComponentGroup{
		Id:   groupID,
		Name: name,
		Uuid: utils.Must(uuid.FromString("ebd1af2e-875f-47e5-8886-4d3baea94d99")),
	}

	model := &componentGroupResourceModel{}
	err := model.fromRemote(spaceID, componentGroup)

	assert.NoError(t, err, "Error occurred during conversion")

	expectedModel := &componentGroupResourceModel{
		ID:      types.StringValue(utils.CreateIdentifier(spaceID, groupID)),
		GroupID: types.Int64Value(groupID),
		UUID:    types.StringValue("ebd1af2e-875f-47e5-8886-4d3baea94d99"),
		Name:    types.StringValue(name),
	}

	assert.Equal(t, expectedModel, model, "Converted model does not match expected model")
}
