/*
 * Copyright 2020 The Compass Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package osb

import (
	"context"
	credentials_provider "credentials-broker/internal/credentialsprovider"
	"fmt"
	"net/http"

	"github.com/kyma-incubator/compass/components/director/pkg/log"
	"github.com/pivotal-cf/brokerapi/v7/domain"
	"github.com/pivotal-cf/brokerapi/v7/domain/apiresponses"
)

type UnbindEndpoint struct {
	credentialsProvider credentials_provider.Client
}

func (b *UnbindEndpoint) Unbind(ctx context.Context, instanceID, bindingID string, details domain.UnbindDetails, asyncAllowed bool) (domain.UnbindSpec, error) {
	log.C(ctx).Infof("Unbind instanceID: %s bindingID: %s details: %+v asyncAllowed: %v", instanceID, bindingID, details, asyncAllowed)

	spec := domain.UnbindSpec{}

	ok, err := b.credentialsProvider.DeleteCredentials(instanceID, bindingID)
	if err != nil {
		log.C(ctx).WithError(err).Errorf("failed to delete binding with ID %s", bindingID)
		return spec, apiresponses.NewFailureResponse(fmt.Errorf("error occurred while executing unbind operation"), http.StatusInternalServerError, "unbind")
	}

	if !ok {
		log.C(ctx).Errorf("Failed to delete binding with ID %s. Credentials with ID %s not found.", bindingID, bindingID)
		return spec, apiresponses.ErrBindingDoesNotExist
	}

	log.C(ctx).Infof("Successfully deleted binding with ID %s", bindingID)

	return spec, nil
}
