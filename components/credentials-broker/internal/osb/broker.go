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

import credentials_provider "credentials-broker/internal/credentialsprovider"

func NewCredentialsBroker(client credentials_provider.Client) *SystemBroker {
	return &SystemBroker{
		CatalogEndpoint: &CatalogEndpoint{},
		ProvisionEndpoint: &ProvisionEndpoint{
			credentialsProvider: client,
		},
		DeprovisionEndpoint: &DeprovisionEndpoint{
			credentialsProvider: client,
		},
		UpdateInstanceEndpoint:        &UpdateInstanceEndpoint{},
		GetInstanceEndpoint:           &GetInstanceEndpoint{},
		InstanceLastOperationEndpoint: &InstanceLastOperationEndpoint{},
		BindEndpoint: &BindEndpoint{
			credentialsProvider: client,
		},
		UnbindEndpoint: &UnbindEndpoint{
			credentialsProvider: client,
		},
		GetBindingEndpoint:        &GetBindingEndpoint{},
		BindLastOperationEndpoint: &BindLastOperationEndpoint{},
	}
}

type SystemBroker struct {
	*CatalogEndpoint
	*ProvisionEndpoint
	*DeprovisionEndpoint
	*UpdateInstanceEndpoint
	*GetInstanceEndpoint
	*InstanceLastOperationEndpoint
	*BindEndpoint
	*UnbindEndpoint
	*GetBindingEndpoint
	*BindLastOperationEndpoint
}
