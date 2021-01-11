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
	"github.com/kyma-incubator/compass/components/director/pkg/log"
	"github.com/pivotal-cf/brokerapi/v7/domain"
	"github.com/pivotal-cf/brokerapi/v7/domain/apiresponses"
	"net/http"
)

type ProvisionEndpoint struct {
	credentialsProvider credentials_provider.Client
}

func (b *ProvisionEndpoint) Provision(ctx context.Context, instanceID string, details domain.ProvisionDetails, asyncAllowed bool) (domain.ProvisionedServiceSpec, error) {
	log.C(ctx).Infof("Provision instance with instanceID: %s, serviceID: %s, planID: %s, parameters: %s context: %s asyncAllowed: %t", instanceID, details.ServiceID, details.PlanID, string(details.RawParameters), string(details.RawContext), asyncAllowed)

	spec := domain.ProvisionedServiceSpec{}

	ok, err := b.credentialsProvider.CreateCredentialsIssuer(instanceID, details)
	if err != nil {
		log.C(ctx).WithError(err).Errorf("failed to create credentials issuer with ID %s", instanceID)
		return spec, apiresponses.NewFailureResponse(fmt.Errorf("error occurred while executing provision operation"), http.StatusInternalServerError, "provision")
	}

	if !ok {
		log.C(ctx).Errorf("Failed to provision service instance with ID %s. Issuer %s already exists.", instanceID, instanceID)
		return spec, apiresponses.ErrInstanceAlreadyExists
	}

	log.C(ctx).Infof("Successfully provisioned service instance with ID %s", instanceID)

	return spec, nil
}
