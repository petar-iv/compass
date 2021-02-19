package model_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/kyma-incubator/compass/components/director/internal2/model"
	"github.com/stretchr/testify/assert"
)

func TestVendorInput_ToVendor(t *testing.T) {
	// given
	id := "foo"
	appID := "bar"
	vendorType := "Sample"
	name := "sample"
	tenant := "tenant"
	labels := json.RawMessage("{}")

	testCases := []struct {
		Name     string
		Input    *model.VendorInput
		Expected *model.Vendor
	}{
		{
			Name: "All properties given",
			Input: &model.VendorInput{
				OrdID:  id,
				Title:  name,
				Type:   vendorType,
				Labels: labels,
			},
			Expected: &model.Vendor{
				OrdID:         id,
				TenantID:      tenant,
				ApplicationID: appID,
				Title:         name,
				Type:          vendorType,
				Labels:        labels,
			},
		},
		{
			Name:     "Nil",
			Input:    nil,
			Expected: nil,
		},
	}

	for _, testCase := range testCases {
		t.Run(fmt.Sprintf("%s", testCase.Name), func(t *testing.T) {

			// when
			result := testCase.Input.ToVendor(tenant, appID)

			// then
			assert.Equal(t, testCase.Expected, result)
		})
	}
}
