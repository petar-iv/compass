/*
This file contains the go structs definitions for the SystemActivation object.
*/

package v1

import (
	"encoding/json"
	metav1 "github.tools.sap/unified-resource-manager/api/pkg/apis/meta/v1"
	"github.tools.sap/unified-resource-manager/api/pkg/apis/runtime"
	"github.tools.sap/unified-resource-manager/api/pkg/apis/runtime/schema"
	"time"
)

// Here we define the SystemActivation type
type SystemActivation struct {
	// This object contains the Group, Version, and Type (GVT)
	// of the SystemActivation, and its metadata
	runtime.ResourceObject `json:",inline"`

	// The spec of the SystemActivation
	Spec SystemActivationSpec `json:"spec,omitempty"`
	// The status of the SystemActivation
	Status SystemActivationStatus `json:"status,omitempty"`
}

type SystemActivationTenantRef struct {
	UUID    string `json:"uuid,omitempty"`
	Name    string `json:"name,omitempty"`
	Group   string `json:"group,omitempty"`
	Version string `json:"version,omitempty"`
	//Type    string `json:"type,omitempty"`
}

type SecretOrigin string

const (
	SecretTypeSecret      SecretOrigin = "Secret"
	SecretTypeDestination SecretOrigin = "Destination"
)

type SystemActivationSecret struct {
	Origin SecretOrigin `json:"origin,omitempty"`
	Name   string       `json:"name,omitempty"`
}

// SystemActivationSpec defines the desired state of the SystemActivation
type SystemActivationSpec struct {
	TenantRef SystemActivationTenantRef `json:"size,omitempty"`
	Secret    SystemActivationSecret    `json:"secret,omitempty"`
	URL       string                    `json:"url,omitempty"`
}

type IntegrationStatus struct {
	Acquired time.Time `json:"acquired,omitempty"`
}

type SecretRefStatus struct {
	Path string `json:"path,omitempty"`
	Name string `json:"name,omitempty"`
}

// SystemActivationStatus defines the observed state of SystemActivation
// The operator creates a secret during the provisioning process and imprints the secret's name on the status
type SystemActivationStatus struct {
	Conditions        []metav1.Condition `json:"conditions"`
	IntegrationStatus *IntegrationStatus `json:"integrationStatus,omitempty"`
	SecretRef         *SecretRefStatus   `json:"secretRef,omitempty"`
}

// SystemActivationList defines a list of SystemActivation objects
type SystemActivationList struct {
	metav1.TypeMeta
	Revision string             `json:"revision"`
	Items    []SystemActivation `json:"items"`
}

// NewSystemActivation initialize an empty SystemActivation
func NewSystemActivation() *SystemActivation {
	serviceInstance := &SystemActivation{}
	serviceInstance.SetGroupVersionType(NewSystemActivationGVT())
	return serviceInstance
}

// NewSystemActivationGVT initialize a SystemActivation Group-Version-Type
func NewSystemActivationGVT() schema.GroupVersionType {
	return schema.GroupVersionType{
		Group:   Group,
		Version: Version,
		Type:    SystemActivationType,
	}
}

// GetResourceType returns SystemActivation's type
func (h *SystemActivation) GetResourceType() schema.ResourceType {
	return h.TypeMeta.GetResourceType()
}

// DeepCopyResource clones a SystemActivation object and returns it
func (h *SystemActivation) DeepCopyResource() runtime.Resource {
	if h == nil {
		return nil
	}
	out := new(SystemActivation)
	h.DeepCopyInto(out)
	return out
}

// DeepCopyInto clones a SystemActivation object into a provided SystemActivation object
func (i *SystemActivation) DeepCopyInto(out *SystemActivation) {
	*out = *i
	out.TypeMeta = i.TypeMeta
	bytes, _ := json.Marshal(i)
	//TODO check error
	_ = json.Unmarshal(bytes, out)

}

func (m *SystemActivation) GetRevision() string {
	return m.Revision
}

func (m *SystemActivation) SetRevision(revision string) {
	m.Revision = revision
}

func (m *SystemActivationList) GetRevision() string {
	return m.Revision
}

func (m *SystemActivationList) SetRevision(revision string) {
	m.Revision = revision
}

// DeepCopyResource clones a SystemActivationList and each of the SystemActivation objects in the list
func (m *SystemActivationList) DeepCopyResource() runtime.Resource {
	var ml SystemActivationList

	// Copy GVT
	ml.TypeMeta.SetGroupVersionType(m.TypeMeta.GroupVersionType())
	ml.SetRevision(m.GetRevision())

	if m.Items != nil {
		in, out := &m.Items, &ml.Items
		*out = make([]SystemActivation, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}

	return &ml
}
