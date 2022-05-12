package model

import (
	"encoding/json"

	rmmeta "github.tools.sap/unified-resource-manager/api/pkg/apis/meta/v1"
	"github.tools.sap/unified-resource-manager/api/pkg/apis/runtime"
	"github.tools.sap/unified-resource-manager/api/pkg/apis/runtime/schema"
)

// OrganizationSpec defines the desired state of Organization
type OrganizationSpec struct {

	// The display name for this Organization
	DisplayName string `json:"displayName,omitempty"`
	// An optional description for this Organization
	Description string `json:"description,omitempty"`
}

// OrganizationStatus defines the observed state of Organization
type OrganizationStatus struct {
	Conditions []rmmeta.Condition `json:"conditions,omitempty"`
}

// Organization is the Schema for the organizations API
type Organization struct {
	runtime.ResourceObject `json:",inline"`

	Spec   OrganizationSpec   `json:"spec,omitempty"`
	Status OrganizationStatus `json:"status,omitempty"`
}

// OrganizationList contains a list of Organization
type OrganizationList struct {
	rmmeta.TypeMeta
	Revision string         `json:"revision"`
	Items    []Organization `json:"items"`
}

func (org *Organization) GetResourceType() schema.ResourceType {
	return org.ResourceObject.TypeMeta.GetResourceType()
}

func (org *Organization) DeepCopy() *Organization {
	out := NewOrganization()

	bytes, err := json.Marshal(org)
	if err != nil {
		return nil
	}

	err = json.Unmarshal(bytes, out)
	if err != nil {
		return nil
	}

	return out
}

func (org *Organization) DeepCopyResource() runtime.Resource {
	if c := org.DeepCopy(); c != nil {
		return c
	}
	return nil
}

func (orgList *OrganizationList) GetRevision() string {
	return orgList.Revision
}

func (orgList *OrganizationList) SetRevision(revision string) {
	orgList.Revision = revision
}

func (orgList *OrganizationList) GetResourceType() schema.ResourceType {
	return orgList.TypeMeta.GetResourceType()
}

func (orgList *OrganizationList) DeepCopyResource() runtime.Resource {
	var out OrganizationList

	// Copy GVT
	out.TypeMeta.SetGroupVersionType(orgList.TypeMeta.GroupVersionType())
	out.SetRevision(orgList.GetRevision())

	// Copy all the items - Resources
	for _, item := range orgList.Items {
		if c := item.DeepCopy(); c != nil {
			out.Items = append(out.Items, *c)
		}

	}

	return &out
}

func NewOrganizationList() *OrganizationList {
	list := OrganizationList{}
	list.Items = []Organization{}
	list.SetGroupVersionType(NewOrganizationGVT())
	return &list
}

func NewOrganizationGVT() schema.GroupVersionType {
	return schema.GroupVersionType{
		Group:   Group,
		Version: Version,
		Type:    "Organization",
	}
}

func NewOrganization() *Organization {
	organization := Organization{}
	organization.SetGroupVersionType(NewOrganizationGVT())
	return &organization
}
