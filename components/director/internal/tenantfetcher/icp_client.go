package tenantfetcher

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"

	"github.com/kyma-incubator/compass/components/director/internal/model"
	"github.com/kyma-incubator/compass/components/director/pkg/tenant"

	"github.com/kyma-incubator/compass/components/director/pkg/log"
	"github.com/pkg/errors"
	"github.com/tidwall/gjson"
)

const PageSize = 80

type ICPConfig struct {
	TokenURL              string `envconfig:"optional,APP_ICP_TOKEN_URL"`
	EntitlementsURL       string `envconfig:"optional,APP_ICP_ENTITLEMENTS_URL"`
	GetTokenAuthorization string `envconfig:"optional,APP_ICP_GET_TOKEN_AUTH"`
}

type Entitlement struct {
	CustomerName string `json:"EndCustomerName"`
	CustomerID   string `json:"EndCustomerID"`
}

type client struct {
	httpClient *http.Client
	config     ICPConfig
	pageSize   int
}

func NewICPClient(httpClient *http.Client, config ICPConfig) *client {
	return &client{
		httpClient: httpClient,
		config:     config,
		pageSize:   PageSize,
	}
}

func (c *client) GetCustomers(ctx context.Context, sinceDate string) ([]model.BusinessTenantMappingInput, error) {
	token, err := c.getAccessToken(ctx)
	if err != nil {
		return nil, err
	}
	skip := 0
	allCustomers := make(map[string]string)
	for {
		customers, err := c.getCustomersPage(ctx, token, c.pageSize, skip, sinceDate)
		if err != nil {
			log.C(ctx).Error(err)
			return nil, err
		}
		if len(customers) == 0 {
			break
		}
		for id, name := range customers {
			allCustomers[id] = name
		}
		skip += c.pageSize
	}

	var tenantsToStore []model.BusinessTenantMappingInput
	for id, name := range allCustomers {
		tenantsToStore = append(tenantsToStore, model.BusinessTenantMappingInput{
			Name:           name,
			ExternalTenant: id,
			Type:           tenant.TypeToStr(tenant.Customer),
			Provider:       "ICP",
		})
	}
	return tenantsToStore, nil
}

func (c *client) getAccessToken(ctx context.Context) (string, error) {
	body := []byte("grant_type=client_credentials")

	req, err := http.NewRequest(http.MethodPost, c.config.TokenURL, bytes.NewBuffer(body))
	if err != nil {
		return "", errors.Wrap(err, "while making token request")
	}
	req.Header.Set("Authorization", "Basic "+c.config.GetTokenAuthorization)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", errors.Wrap(err, "while getting token")
	}

	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.C(ctx).WithError(err).Errorf("An error has occurred while closing response body: %v", err)
		}
	}()
	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", errors.Wrap(err, "while reading token response body")
	}

	token := gjson.Get(string(responseBody), "access_token")

	return token.String(), nil
}

func (c *client) getCustomersPage(ctx context.Context, token string, pageSize int, skip int, sinceDate string) (map[string]string, error) {
	requestURl, err := c.buildRequestURL(c.config.EntitlementsURL, pageSize, skip, sinceDate)
	if err != nil {
		return nil, err
	}
	log.C(ctx).Infof("request url: %s", requestURl)
	log.C(ctx).Infof("token: %s", token)
	req, err := http.NewRequest(http.MethodGet, requestURl, nil)
	if err != nil {
		return nil, errors.Wrap(err, "while get entitlements request")
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "while getting token")
	}

	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.C(ctx).WithError(err).Errorf("An error has occurred while closing response body: %v", err)
		}
	}()
	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "while reading token response body")
	}

	entitlementsAsString := gjson.Get(string(responseBody), "d.results")

	var entitlements []Entitlement
	err = json.Unmarshal([]byte(entitlementsAsString.String()), &entitlements)
	if err != nil {
		return nil, err
	}

	customers := make(map[string]string)
	for _, entitlement := range entitlements {
		customers[entitlement.CustomerID] = entitlement.CustomerName
	}

	log.C(ctx).Info("Successfully fetched customers info from entitlements")
	return customers, nil
}

func (c *client) buildRequestURL(endpoint string, pageSize int, skip int, sinceDate string) (string, error) {
	queryParams := make(map[string]string)
	queryParams["$top"] = strconv.Itoa(pageSize)
	queryParams["$skip"] = strconv.Itoa(skip)
	queryParams["$orderby"] = "DateOfLastChange desc"
	queryParams["$select"] = "EndCustomerName,EndCustomerID"
	queryParams["$filter"] = fmt.Sprintf("DateOfLastChange gt datetimeoffset'%s'", sinceDate)

	u, err := url.Parse(endpoint)
	if err != nil {
		return "", err
	}

	q, err := url.ParseQuery(u.RawQuery)
	if err != nil {
		return "", err
	}

	for qKey, qValue := range queryParams {
		q.Add(qKey, qValue)
	}

	u.RawQuery = q.Encode()

	return u.String(), nil
}
