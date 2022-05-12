package model

import (
	"encoding/json"

	rmmeta "github.tools.sap/unified-resource-manager/api/pkg/apis/meta/v1"
	"github.tools.sap/unified-resource-manager/api/pkg/apis/runtime"
	"github.tools.sap/unified-resource-manager/api/pkg/apis/runtime/schema"
)

const (
	Group   = "accounts.resource.api.sap"
	Version = "v1"
)

// FolderSpec defines the desired state of Folder
type FolderSpec struct {
	DisplayName string `json:"displayName,omitempty"`
	Description string `json:"description,omitempty"`
}

// FolderStatus defines the observed state of Folder
type FolderStatus struct {
	Conditions []rmmeta.Condition `json:"conditions,omitempty"`
}

// Folder is the Schema for the folders API.
// Folders allow the grouping of resources that fit a common business or technical need and to set boundaries between different resource groups in a multi-level hierarchy within an organization base.
// Folders are used mostly for structuring the account entities and resources.
// Folders can contain other folders, and each folder can have only a single parent. And each resource group can have only a single folder as its parent.
type Folder struct {
	runtime.ResourceObject `json:",inline"`

	Spec   FolderSpec   `json:"spec,omitempty"`
	Status FolderStatus `json:"status,omitempty"`
}

// FolderList contains a list of Folder
type FolderList struct {
	rmmeta.TypeMeta
	Revision string   `json:"revision"`
	Items    []Folder `json:"items"`
}

func (folder *Folder) GetResourceType() schema.ResourceType {
	return folder.ResourceObject.TypeMeta.GetResourceType()
}

func (folder *Folder) DeepCopy() *Folder {
	out := NewFolder()

	bytes, err := json.Marshal(folder)
	if err != nil {
		return nil
	}

	err = json.Unmarshal(bytes, out)
	if err != nil {
		return nil
	}

	return out
}

func (folder *Folder) DeepCopyResource() runtime.Resource {
	if c := folder.DeepCopy(); c != nil {
		return c
	}
	return nil
}

func (l *FolderList) GetRevision() string {
	return l.Revision
}

func (l *FolderList) SetRevision(revision string) {
	l.Revision = revision
}

func (l *FolderList) GetResourceType() schema.ResourceType {
	return l.TypeMeta.GetResourceType()
}

func (l *FolderList) DeepCopyResource() runtime.Resource {
	var out FolderList

	out.TypeMeta.SetGroupVersionType(l.TypeMeta.GroupVersionType())
	out.SetRevision(l.GetRevision())

	for _, item := range l.Items {
		if c := item.DeepCopy(); c != nil {
			out.Items = append(out.Items, *c)
		}

	}

	return &out
}

func NewFolderGVT() schema.GroupVersionType {
	return schema.GroupVersionType{
		Group:   Group,
		Version: Version,
		Type:    "Folder",
	}
}

func NewFolder() *Folder {
	folder := Folder{}
	folder.SetGroupVersionType(NewFolderGVT())
	return &folder
}
