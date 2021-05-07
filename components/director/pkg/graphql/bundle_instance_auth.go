package graphql

import (
	"encoding/json"

	"github.com/kyma-incubator/compass/components/director/pkg/resource"
)

type BundleInstanceAuth struct {
	*BaseEntity
	// Context of BundleInstanceAuth - such as Runtime ID, namespace
	Context *JSON `json:"context"`
	// User input while requesting Bundle Instance Auth
	InputParams *JSON `json:"inputParams"`
	// It may be empty if status is PENDING.
	// Populated with `bundle.defaultAuth` value if `bundle.defaultAuth` is defined. If not, Compass notifies Application/Integration System about the Auth request.
	Auth   *Auth                     `json:"auth"`
	Status *BundleInstanceAuthStatus `json:"status"`
}

func (*BundleInstanceAuth) GetType() resource.Type {
	return resource.BundleInstanceAuth
}

func (e *BundleInstanceAuth) Sentinel() {}

func (e *BundleInstanceAuth) Template() map[string]interface{}{
	var inputParams map[string]interface{}
	if e.InputParams != nil {
		json.Unmarshal([]byte(*e.InputParams), &inputParams)
	}
	return map[string]interface{}{
		"ID": e.ID,
		"inputParams": inputParams,
	}
}