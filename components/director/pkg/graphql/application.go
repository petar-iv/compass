package graphql

import (
	"net/url"

	"github.com/kyma-incubator/compass/components/director/pkg/log"
	"github.com/kyma-incubator/compass/components/director/pkg/resource"
)

type Application struct {
	Name                  string             `json:"name"`
	ProviderName          *string            `json:"providerName"`
	IntegrationSystemID   *string            `json:"integrationSystemID"`
	ApplicationTemplateID *string            `json:"applicationTemplateID"`
	Description           *string            `json:"description"`
	Status                *ApplicationStatus `json:"status"`
	HealthCheckURL        *string            `json:"healthCheckURL"`
	BaseURL               *string            `json:"baseURL"`
	*BaseEntity
}

func (e *Application) GetType() resource.Type {
	return resource.Application
}

func (e *Application) Sentinel() {}

func (e *Application) Template() map[string]interface{} {
	// TODO add host
	var host string
	if e.BaseURL != nil && len(*e.BaseURL) != 0 {
		if baseURL, err := url.Parse(*e.BaseURL); err == nil {
			host = baseURL.Host
		} else {
			log.D().Errorf("Failed to parse URL for base URL %q of application with ID %s", e.BaseURL, e.ID)
		}
	}

	return map[string]interface{}{
		"ID":      e.ID,
		"Name":    e.Name,
		"BaseURL": e.BaseURL,
		"Host":    host,
	}
}

// Extended types used by external API

type ApplicationPageExt struct {
	ApplicationPage
	Data []*ApplicationExt `json:"data"`
}

type ApplicationExt struct {
	Application
	Labels                Labels                           `json:"labels"`
	Webhooks              []Webhook                        `json:"webhooks"`
	Auths                 []*SystemAuth                    `json:"auths"`
	Bundle                BundleExt                        `json:"bundle"`
	Bundles               BundlePageExt                    `json:"bundles"`
	EventingConfiguration ApplicationEventingConfiguration `json:"eventingConfiguration"`
}
