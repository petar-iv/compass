package tenantfetchersvc

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/kyma-incubator/compass/components/director/internal/features"
	"github.com/kyma-incubator/compass/components/director/internal/model"
	"github.com/kyma-incubator/compass/components/director/internal/tenantfetcher"
	"github.com/kyma-incubator/compass/components/director/pkg/graphql"
	"github.com/kyma-incubator/compass/components/director/pkg/log"
	"github.com/kyma-incubator/compass/components/director/pkg/oauth"
	"github.com/kyma-incubator/compass/components/director/pkg/persistence"
	"github.com/kyma-incubator/compass/components/director/pkg/tenant"
	"github.com/tidwall/gjson"
)

const (
	// InternalServerError message
	InternalServerError = "Internal Server Error"
	compassURL          = "https://github.com/kyma-incubator/compass"
	TenantProviderName  = "Atom"
)

// TenantFetcher is used to fectch tenants for creation;
//go:generate mockery --name=TenantFetcher --output=automock --outpkg=automock --case=underscore
type TenantFetcher interface {
	FetchTenantOnDemand(ctx context.Context, tenantID string) error
}

// TenantSubscriber is used to apply subscription changes for tenants;
//go:generate mockery --name=TenantSubscriber --output=automock --outpkg=automock --case=underscore
type TenantSubscriber interface {
	Subscribe(ctx context.Context, tenantSubscriptionRequest *TenantSubscriptionRequest) error
	Unsubscribe(ctx context.Context, tenantSubscriptionRequest *TenantSubscriptionRequest) error
}

// HandlerConfig is the configuration required by the tenant handler.
// It includes configurable parameters for incoming requests, including different tenant IDs json properties, and path parameters.
type HandlerConfig struct {
	TenantOnDemandHandlerEndpoint string `envconfig:"APP_TENANT_ON_DEMAND_HANDLER_ENDPOINT,default=/v1/fetch/{tenantId}"`
	RegionalHandlerEndpoint       string `envconfig:"APP_REGIONAL_HANDLER_ENDPOINT,default=/v1/regional/{region}/callback/{tenantId}"`
	AtomTenantsEndpoint           string `envconfig:"APP_ATOM_TENANTS_ENDPOINT,default=/v1/atom"`
	DependenciesEndpoint          string `envconfig:"APP_DEPENDENCIES_ENDPOINT,default=/v1/dependencies"`
	TenantPathParam               string `envconfig:"APP_TENANT_PATH_PARAM,default=tenantId"`
	RegionPathParam               string `envconfig:"APP_REGION_PATH_PARAM,default=region"`

	DirectorGraphQLEndpoint     string        `envconfig:"APP_DIRECTOR_GRAPHQL_ENDPOINT"`
	ClientTimeout               time.Duration `envconfig:"default=60s"`
	HTTPClientSkipSslValidation bool          `envconfig:"APP_HTTP_CLIENT_SKIP_SSL_VALIDATION,default=false"`

	TenantProviderConfig
	features.Config

	Database persistence.DatabaseConfig
}

// TenantProviderConfig includes the configuration for tenant providers - the tenant ID json property names, the subdomain property name, and the tenant provider name.
type TenantProviderConfig struct {
	TenantIDProperty               string `envconfig:"APP_TENANT_PROVIDER_TENANT_ID_PROPERTY,default=tenantId"`
	SubaccountTenantIDProperty     string `envconfig:"APP_TENANT_PROVIDER_SUBACCOUNT_TENANT_ID_PROPERTY,default=subaccountTenantId"`
	CustomerIDProperty             string `envconfig:"APP_TENANT_PROVIDER_CUSTOMER_ID_PROPERTY,default=customerId"`
	SubdomainProperty              string `envconfig:"APP_TENANT_PROVIDER_SUBDOMAIN_PROPERTY,default=subdomain"`
	TenantProvider                 string `envconfig:"APP_TENANT_PROVIDER,default=external-provider"`
	SubscriptionProviderIDProperty string `envconfig:"APP_TENANT_PROVIDER_SUBSCRIPTION_PROVIDER_ID_PROPERTY,default=subscriptionProviderId"`
}

// EventsConfig contains configuration for Events API requests
type EventsConfig struct {
	OAuthConfig        tenantfetcher.OAuth2Config
	APIConfig          tenantfetcher.APIConfig
	AuthMode           oauth.AuthMode `envconfig:"APP_OAUTH_AUTH_MODE,default=standard"`
	QueryConfig        tenantfetcher.QueryConfig
	TenantFieldMapping tenantfetcher.TenantFieldMapping
}

type handler struct {
	fetcher         TenantFetcher
	gqlClient       DirectorGraphQLClient
	tenantConverter TenantConverter
	subscriber      TenantSubscriber
	config          HandlerConfig
}

// NewTenantsHTTPHandler returns a new HTTP handler, responsible for creation and deletion of regional and non-regional tenants.
func NewTenantsHTTPHandler(subscriber TenantSubscriber, config HandlerConfig) *handler {
	return &handler{
		subscriber: subscriber,
		config:     config,
	}
}

// NewTenantFetcherHTTPHandler returns a new HTTP handler, responsible for creation of on-demand tenants.
func NewTenantFetcherHTTPHandler(fetcher TenantFetcher, config HandlerConfig) *handler {
	return &handler{
		fetcher: fetcher,
		config:  config,
	}
}

// NewTenantWriterHandler missing godoc.
func NewTenantWriterHandler(gqlClient DirectorGraphQLClient, tenantConverter TenantConverter) *handler {
	return &handler{
		gqlClient:       gqlClient,
		tenantConverter: tenantConverter,
	}
}

// FetchTenantOnDemand fetches External tenants registry events for a provided subaccount and creates a subaccount tenant
func (h *handler) FetchTenantOnDemand(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()

	vars := mux.Vars(request)
	tenantID, ok := vars[h.config.TenantPathParam]
	if !ok {
		log.C(ctx).WithError(errors.New("tenant path parameter is missing from request")).Error()
		http.Error(writer, "Tenant path parameter is missing from request", http.StatusBadRequest)
		return
	}

	log.C(ctx).Infof("Fetching create event for tenant with ID %s", tenantID)

	err := h.fetcher.FetchTenantOnDemand(ctx, tenantID)
	if err != nil {
		log.C(ctx).WithError(err).Errorf("Error while processing request for creation of tenant %s: %v", tenantID, err)
		http.Error(writer, InternalServerError, http.StatusInternalServerError)
		return
	}
	writeCreatedResponse(writer, ctx, tenantID)
}

func (h *handler) StoreAtomTenants(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()

	defer func() {
		if err := request.Body.Close(); err != nil {
			log.C(ctx).Error("Got error on closing request body", err)
		}
	}()

	b, err := io.ReadAll(request.Body)
	if err != nil {
		log.C(ctx).Error(err)
	}

	var payload RequestPayload
	if err = json.Unmarshal(b, &payload); err != nil {
		log.C(ctx).Error(err)
	}

	tenantsToCreate := getTenantsToBeCreated(payload)

	maxChunkSize := 100
	tenantsToCreateGQL := h.tenantConverter.MultipleInputToGraphQLInput(tenantsToCreate)
	err = executeInChunks(ctx, tenantsToCreateGQL, func(ctx context.Context, chunk []graphql.BusinessTenantMappingInput) error {
		return h.gqlClient.WriteTenants(ctx, chunk)
	}, maxChunkSize)

	if err != nil {
		log.C(ctx).Error(err)
	}

	writer.Header().Set("Content-Type", "text/plain")
	writer.WriteHeader(http.StatusOK)
}

func writeCreatedResponse(writer http.ResponseWriter, ctx context.Context, tenantID string) {
	writer.Header().Set("Content-Type", "text/plain")
	writer.WriteHeader(http.StatusOK)
	if _, err := writer.Write([]byte(compassURL)); err != nil {
		log.C(ctx).WithError(err).Errorf("Failed to write response body for request for creation of tenant %s: %v", tenantID, err)
	}
}

// SubscribeTenant handles subscription for tenant. If tenant does not exist, will create it first.
func (h *handler) SubscribeTenant(writer http.ResponseWriter, request *http.Request) {
	h.applySubscriptionChange(writer, request, h.subscriber.Subscribe)
}

// UnSubscribeTenant handles unsubscription for tenant which will remove the tenant id label from the runtime
func (h *handler) UnSubscribeTenant(writer http.ResponseWriter, request *http.Request) {
	h.applySubscriptionChange(writer, request, h.subscriber.Unsubscribe)
}

// Dependencies handler returns all external services where once created in Compass, the tenant should be created as well.
func (h *handler) Dependencies(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	if _, err := writer.Write([]byte("{}")); err != nil {
		log.C(request.Context()).WithError(err).Errorf("Failed to write response body for dependencies request")
		return
	}
}

func (h *handler) applySubscriptionChange(writer http.ResponseWriter, request *http.Request, subscriptionFunc subscriptionFunc) {
	ctx := request.Context()

	vars := mux.Vars(request)
	region, ok := vars[h.config.RegionPathParam]
	if !ok {
		log.C(ctx).Error("Region path parameter is missing from request")
		http.Error(writer, "Region path parameter is missing from request", http.StatusBadRequest)
		return
	}

	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		log.C(ctx).WithError(err).Errorf("Failed to read tenant information from request body: %v", err)
		http.Error(writer, InternalServerError, http.StatusInternalServerError)
		return
	}

	subscriptionRequest, err := h.getSubscriptionRequest(body, region)
	if err != nil {
		log.C(ctx).WithError(err).Errorf("Failed to extract tenant information from request body: %v", err)
		http.Error(writer, fmt.Sprintf("Failed to extract tenant information from request body: %v", err), http.StatusBadRequest)
		return
	}

	mainTenantID := subscriptionRequest.MainTenantID()
	if err := subscriptionFunc(ctx, subscriptionRequest); err != nil {
		log.C(ctx).WithError(err).Errorf("Failed to apply subscription change for tenant %s: %v", mainTenantID, err)
		http.Error(writer, InternalServerError, http.StatusInternalServerError)
		return
	}

	respondSuccess(ctx, writer, mainTenantID)
}

func (h *handler) getSubscriptionRequest(body []byte, region string) (*TenantSubscriptionRequest, error) {
	properties, err := getProperties(body, map[string]bool{
		h.config.TenantIDProperty:               true,
		h.config.SubaccountTenantIDProperty:     false,
		h.config.SubdomainProperty:              true,
		h.config.CustomerIDProperty:             false,
		h.config.SubscriptionProviderIDProperty: true,
	})
	if err != nil {
		return nil, err
	}

	req := &TenantSubscriptionRequest{
		AccountTenantID:        properties[h.config.TenantIDProperty],
		SubaccountTenantID:     properties[h.config.SubaccountTenantIDProperty],
		CustomerTenantID:       properties[h.config.CustomerIDProperty],
		Subdomain:              properties[h.config.SubdomainProperty],
		SubscriptionProviderID: properties[h.config.SubscriptionProviderIDProperty],
		Region:                 region,
	}

	if req.AccountTenantID == req.SubaccountTenantID {
		req.SubaccountTenantID = ""
	}

	if req.AccountTenantID == req.CustomerTenantID {
		req.CustomerTenantID = ""
	}

	return req, nil
}

func getProperties(body []byte, props map[string]bool) (map[string]string, error) {
	resultProps := map[string]string{}
	for propName, mandatory := range props {
		result := gjson.GetBytes(body, propName).String()
		if mandatory && len(result) == 0 {
			return nil, fmt.Errorf("mandatory property %q is missing from request body", propName)
		}
		resultProps[propName] = result
	}

	return resultProps, nil
}

func respondSuccess(ctx context.Context, writer http.ResponseWriter, mainTenantID string) {
	writer.Header().Set("Content-Type", "text/plain")
	writer.WriteHeader(http.StatusOK)
	if _, err := writer.Write([]byte(compassURL)); err != nil {
		log.C(ctx).WithError(err).Errorf("Failed to write response body for tenant request creation for tenant %s: %v", mainTenantID, err)
	}
}

func getTenantsToBeCreated(payload RequestPayload) []model.BusinessTenantMappingInput {
	var toBeCreated []model.BusinessTenantMappingInput
	if len(payload.Customer) > 0 {
		toBeCreated = append(toBeCreated, model.BusinessTenantMappingInput{
			Name:           payload.Customer,
			ExternalTenant: payload.Customer,
			Type:           tenant.TypeToStr(tenant.Customer),
			Provider:       TenantProviderName,
		})
	}
	if payload.Organization != nil {
		toBeCreated = append(toBeCreated, model.BusinessTenantMappingInput{
			Name:           payload.Organization.Name,
			ExternalTenant: payload.Organization.Path,
			Parent:         payload.Customer,
			Type:           tenant.TypeToStr(tenant.Organization),
			Provider:       TenantProviderName,
		})
	}
	for _, folder := range payload.Folders {
		lastFolder := toBeCreated[len(toBeCreated)-1]
		toBeCreated = append(toBeCreated, model.BusinessTenantMappingInput{
			Name:           folder.Name,
			ExternalTenant: folder.Path,
			Parent:         lastFolder.ExternalTenant,
			Type:           tenant.TypeToStr(tenant.Folder),
			Provider:       TenantProviderName,
		})
	}
	if payload.ResourceGroup != nil {
		parent := toBeCreated[len(toBeCreated)-1]
		toBeCreated = append(toBeCreated, model.BusinessTenantMappingInput{
			Name:           payload.ResourceGroup.Name,
			ExternalTenant: payload.ResourceGroup.Path,
			Parent:         parent.ExternalTenant,
			Type:           tenant.TypeToStr(tenant.ResourceGroup),
			Provider:       TenantProviderName,
		})
	}
	return toBeCreated
}

func executeInChunks(ctx context.Context, tenants []graphql.BusinessTenantMappingInput, f func(ctx context.Context, chunk []graphql.BusinessTenantMappingInput) error, maxChunkSize int) error {
	for {
		if len(tenants) == 0 {
			return nil
		}
		chunkSize := int(math.Min(float64(len(tenants)), float64(maxChunkSize)))
		tenantsChunk := tenants[:chunkSize]
		if err := f(ctx, tenantsChunk); err != nil {
			return err
		}
		tenants = tenants[chunkSize:]
	}
}

type Tenant struct {
	Name string `json:"name"`
	Path string `json:"path"`
}

type RequestPayload struct {
	Customer      string    `json:"customer"`
	Organization  *Tenant   `json:"organization"`
	Folders       []*Tenant `json:"folders,omitempty"`
	ResourceGroup *Tenant   `json:"resource_group,omitempty"`
}
