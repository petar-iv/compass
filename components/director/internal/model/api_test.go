package model_test

import (
	"testing"

	"github.com/kyma-incubator/compass/components/director/internal/domain/api"

	"github.com/kyma-incubator/compass/components/director/internal/model"
	"github.com/stretchr/testify/assert"
)

func TestAPIDefinitionInput_ToAPIDefinitionWithBundleID(t *testing.T) {
	// GIVEN
	id := "foo"
	appID := "baz"
	desc := "Sample"
	name := "sample"
	targetURL := "https://foo.bar"
	group := "sampleGroup"

	testCases := []struct {
		Name     string
		Input    *model.APIDefinitionInput
		Expected *model.APIDefinition
	}{
		{
			Name: "All properties given",
			Input: &model.APIDefinitionInput{
				Name:        name,
				Description: &desc,
				TargetURLs:  api.ConvertTargetURLToJSONArray(targetURL),
				Group:       &group,
			},
			Expected: &model.APIDefinition{
				ApplicationID: appID,
				Name:          name,
				Description:   &desc,
				TargetURLs:    api.ConvertTargetURLToJSONArray(targetURL),
				Group:         &group,
				BaseEntity: &model.BaseEntity{
					ID:    id,
					Ready: true,
				},
			},
		},
		{
			Name:     "Nil",
			Input:    nil,
			Expected: nil,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			// WHEN
			result := testCase.Input.ToAPIDefinitionWithinBundle(id, appID, 0)

			// THEN
			assert.Equal(t, testCase.Expected, result)
		})
	}
}
