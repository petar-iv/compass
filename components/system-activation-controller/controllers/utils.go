package controllers

import (
	v1 "github.com/kyma-incubator/compass/components/system-activation-controller/api/v1"
	"github.tools.sap/unified-resource-manager/api/pkg/apis/errors"
	"github.tools.sap/unified-resource-manager/api/pkg/apis/runtime"
	"net/http"
)

// IsObjectBeingDeleted returns True if a DeletionTime is imprinted on the object's Metadata
func IsObjectBeingDeleted(sampleDB *v1.SystemActivation) bool {
	return sampleDB.GetDeletionTime() != nil && !sampleDB.GetDeletionTime().IsZero()
}

// IsNotFound return True if the error code is StatusNotFound (404)
func IsNotFound(err error) bool {
	status := ErrorStatus(err)
	return status == http.StatusNotFound
}

// IgnoreNotFound returns nil if IsNotFound return True, and the error otherwise
func IgnoreNotFound(err error) error {
	if IsNotFound(err) {
		return nil
	}
	return err
}

// IsResourceAlreadyExists return True if the error status code is  StatusConflict (409)
func IsResourceAlreadyExists(err error) bool {
	status := ErrorStatus(err)
	return status == http.StatusConflict
}

// ErrorStatus extracts the error status code
func ErrorStatus(err error) int {
	if err == nil {
		return 0
	}
	if e, ok := err.(errors.ErrorDetails); ok {
		return int(e.ErrorStatus.Code)
	}
	return 0
}

func FormatMessageResource(message string, resourceKey runtime.ResourceKey) string {
	return FormatMessagePathName(message, resourceKey.GetPath(), resourceKey.GetName())
}

func FormatMessagePathName(message string, path string, name string) string {
	return message + "\t path: " + path + "\tname: " + name
}
