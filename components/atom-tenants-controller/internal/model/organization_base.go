package model

import (
	"encoding/json"

	rmmeta "github.tools.sap/unified-resource-manager/api/pkg/apis/meta/v1"
	"github.tools.sap/unified-resource-manager/api/pkg/apis/runtime"
	"github.tools.sap/unified-resource-manager/api/pkg/apis/runtime/schema"
)

// OrganizationBaseSpec defines the desired state of OrganizationBase
type OrganizationBaseSpec struct {
	DisplayName   string   `json:"displayName,omitempty"`
	Description   string   `json:"description,omitempty"`
	InitialAdmins []string `json:"initialAdmins"`
}

// OrganizationBaseStatus defines the observed state of OrganizationBase
type OrganizationBaseStatus struct {
	// OrganizationName is the name of the Organization owned by this OrganizationBase
	OrganizationName string `json:"organizationName,omitempty"`
	// AdminRoleBindingName is the name of the RoleBinding that binds the initial admins of the OrganizationBase to the Role accounts.admin
	AdminRoleBindingName string `json:"adminRoleBindingName,omitempty"`

	// AdminRoleBindingName is the name of the RoleBinding that binds the account viewers of the OrganizationBase to the Role accounts.viewer
	ViewerRoleBindingName string `json:"viewerRoleBindingName,omitempty"`

	Conditions []rmmeta.Condition `json:"conditions,omitempty"`
}

func (obs *OrganizationBaseStatus) SetAdminRoleBindingName(name string) {
	obs.AdminRoleBindingName = name
}

func (obs *OrganizationBaseStatus) SetViewerRoleBindingName(name string) {
	obs.ViewerRoleBindingName = name
}

// OrganizationBase is the Schema for the organizations API.
// Organizations are the root entity of the multi-level hierarchy for an SAP Intelligent Enterprise account, and comprise folders and resource groups.
// OrganizationBase represents the contractual entity of an SAP customer or partner, and includes all eligible products that the customer's organization base is entitled to use across SAP.
// The organization base is also the context used for connecting across all SAP components and business processes, for example, for billing SAP customers and partners.
type OrganizationBase struct {
	//metav1.TypeMeta   `json:",inline"`
	//metav1.ObjectMeta `json:"metadata,omitempty"`
	runtime.ResourceObject `json:",inline"`

	Spec   OrganizationBaseSpec   `json:"spec,omitempty"`
	Status OrganizationBaseStatus `json:"status,omitempty"`
}

func (ob *OrganizationBase) GetResourceType() schema.ResourceType {
	return ob.ResourceObject.TypeMeta.GetResourceType()
}

func (ob *OrganizationBase) DeepCopy() *OrganizationBase {
	out := NewOrganizationBase()

	bytes, err := json.Marshal(ob)
	if err != nil {
		return nil
	}

	err = json.Unmarshal(bytes, out)
	if err != nil {
		return nil
	}

	return out
}

func (ob *OrganizationBase) DeepCopyResource() runtime.Resource {
	if c := ob.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// OrganizationBaseList contains a list of OrganizationBase
type OrganizationBaseList struct {
	rmmeta.TypeMeta
	Revision string             `json:"revision"`
	Items    []OrganizationBase `json:"items"`
}

func (obl *OrganizationBaseList) GetRevision() string {
	return obl.Revision
}

func (obl *OrganizationBaseList) SetRevision(revision string) {
	obl.Revision = revision
}

func (obl *OrganizationBaseList) GetResourceType() schema.ResourceType {
	return obl.TypeMeta.GetResourceType()
}

func (obl *OrganizationBaseList) DeepCopyResource() runtime.Resource {
	var out OrganizationBaseList
	// Copy GVT
	out.TypeMeta.SetGroupVersionType(obl.TypeMeta.GroupVersionType())
	out.SetRevision(obl.GetRevision())

	// Copy all the items - Resources
	for _, item := range obl.Items {
		if c := item.DeepCopy(); c != nil {
			out.Items = append(out.Items, *c)
		}
	}
	return &out
}

func (obs *OrganizationBaseStatus) DeepCopy() *OrganizationBaseStatus {
	out := NewOrganizationBaseStatus()

	bytes, err := json.Marshal(obs)
	if err != nil {
		return nil
	}

	err = json.Unmarshal(bytes, out)
	if err != nil {
		return nil
	}

	return out
}

func NewOrganizationBaseGVT() schema.GroupVersionType {
	return schema.GroupVersionType{
		Group:   Group,
		Version: Version,
		Type:    "OrganizationBase",
	}
}

func NewOrganizationBase() *OrganizationBase {
	organizationBase := OrganizationBase{}
	organizationBase.SetGroupVersionType(NewOrganizationBaseGVT())
	return &organizationBase
}

func NewOrganizationBaseStatus() *OrganizationBaseStatus {
	organizationBaseStatus := OrganizationBaseStatus{}
	return &organizationBaseStatus
}
