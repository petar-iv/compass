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

	"github.com/pivotal-cf/brokerapi/v7/domain"
)

type CatalogEndpoint struct{}

func (b *CatalogEndpoint) Services(ctx context.Context) ([]domain.Service, error) {
	return []domain.Service{
		{
			ID:          "68edbadf-20f6-40b5-b21b-a37073a1d314",
			Name:        "cmp",
			Description: "Central Management Plane",
			Bindable:    true,
			Metadata: &domain.ServiceMetadata{
				DisplayName: "Central Management Plane",
			},
			Plans: []domain.ServicePlan{
				{
					ID:          "d949557d-8805-4b2f-85da-cc12dd0c6b95",
					Name:        "ord-access",
					Description: "Access to Open Resource Discovery Service",
					Metadata:    &domain.ServicePlanMetadata{},
				},
			},
		},
	}, nil
}
