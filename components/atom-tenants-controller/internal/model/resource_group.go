package model

import (
	"encoding/json"

	rmmeta "github.tools.sap/unified-resource-manager/api/pkg/apis/meta/v1"
	"github.tools.sap/unified-resource-manager/api/pkg/apis/runtime"
	"github.tools.sap/unified-resource-manager/api/pkg/apis/runtime/schema"
)

// ResourceGroupSpec defines the desired state of ResourceGroup
type ResourceGroupSpec struct {
	DisplayName string `json:"displayName,omitempty"`
	Description string `json:"description,omitempty"`
}

// ResourceGroupStatus defines the observed state of ResourceGroup
type ResourceGroupStatus struct {
	Conditions []rmmeta.Condition `json:"conditions,omitempty"`
}

// ResourceGroup is the Schema for the resource-groups API.
// Resource groups are logical containers for managing several resources that share the same lifecycle.
// Resources can be added to and removed from a resource group, and can be moved between resource groups at any time; however, a resource can exist in only one resource group.
// A resource can connect to resources in other resource groups, for example, when two related resources do not share the same lifecycle.
// Resource groups can be scoped for access control of administrative actions.
// The resource group is the common workspace for interacting with all providers in the SAP Resource Manager deployment and management service, such as CRUD operations, quota assignments, billing and usage data, access controls, custom labels and tags, and more.
// Resource groups are not bound to any specific security model or platform region.
type ResourceGroup struct {
	runtime.ResourceObject `json:",inline"`

	Spec   ResourceGroupSpec   `json:"spec,omitempty"`
	Status ResourceGroupStatus `json:"status,omitempty"`
}

// ResourceGroupList contains a list of ResourceGroup
type ResourceGroupList struct {
	rmmeta.TypeMeta
	Revision string          `json:"revision"`
	Items    []ResourceGroup `json:"items"`
}

func (rg *ResourceGroup) GetResourceType() schema.ResourceType {
	return rg.ResourceObject.TypeMeta.GetResourceType()
}

func (rg *ResourceGroup) DeepCopy() *ResourceGroup {
	out := NewResourceGroup()

	bytes, err := json.Marshal(rg)
	if err != nil {
		return nil
	}

	err = json.Unmarshal(bytes, out)
	if err != nil {
		return nil
	}

	return out
}

func (rg *ResourceGroup) DeepCopyResource() runtime.Resource {
	if c := rg.DeepCopy(); c != nil {
		return c
	}
	return nil
}

func (rgl *ResourceGroupList) GetRevision() string {
	return rgl.Revision
}

func (rgl *ResourceGroupList) SetRevision(revision string) {
	rgl.Revision = revision
}

func (rgl *ResourceGroupList) GetResourceType() schema.ResourceType {
	return rgl.TypeMeta.GetResourceType()
}

func (rgl *ResourceGroupList) DeepCopyResource() runtime.Resource {
	var out ResourceGroupList

	// Copy GVT
	out.TypeMeta.SetGroupVersionType(rgl.TypeMeta.GroupVersionType())
	out.SetRevision(rgl.GetRevision())

	// Copy all the items - Resources
	for _, item := range rgl.Items {
		if c := item.DeepCopy(); c != nil {
			out.Items = append(out.Items, *c)
		}

	}

	return &out
}

func NewResourceGroupGVT() schema.GroupVersionType {
	return schema.GroupVersionType{
		Group:   Group,
		Version: Version,
		Type:    "ResourceGroup",
	}
}

func NewResourceGroup() *ResourceGroup {
	resourceGroup := ResourceGroup{}
	resourceGroup.SetGroupVersionType(NewResourceGroupGVT())
	return &resourceGroup
}
