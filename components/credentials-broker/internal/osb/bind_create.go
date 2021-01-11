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

type BindEndpoint struct {
	credentialsProvider credentials_provider.Client
}

func (b *BindEndpoint) Bind(ctx context.Context, instanceID, bindingID string, details domain.BindDetails, asyncAllowed bool) (domain.Binding, error) {
	log.C(ctx).Infof("Bind instanceID: %s bindingID: %s parameters: %s context: %s asyncAllowed: %t", instanceID, bindingID, string(details.RawParameters), string(details.RawContext), asyncAllowed)

	binding := domain.Binding{}

	credentials, err := b.credentialsProvider.CreateCredentials(instanceID, bindingID, details)
	if err != nil {
		log.C(ctx).WithError(err).Errorf("Failed to create credentials %s from issuer with ID %s", bindingID, instanceID)
		return binding, apiresponses.NewFailureResponse(fmt.Errorf("error occurred while executing bind operation"), http.StatusInternalServerError, "bind")
	}

	log.C(ctx).Infof("Successfully created binding %s", bindingID)

	return domain.Binding{
		Credentials: credentials, // TODO: Add ORD Service URL in the binding
	}, nil
}
