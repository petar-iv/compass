package credentialsprovider

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/pivotal-cf/brokerapi/v7/domain"
	"github.com/tidwall/gjson"
	"io/ioutil"
	"net/http"
	"net/url"
)

// HTTPClient should be implemented by clients which can send HTTP requests
type HTTPClient interface {
	Do(*http.Request) (*http.Response, error)
}

// Credentials provider client interface
type Client interface {
	CreateCredentialsIssuer(id string, provisionDetails domain.ProvisionDetails) (bool, error)
	DeleteCredentialsIssuer(id string) (bool, error)
	CreateCredentials(issuerID string, id string, bindingDetails domain.BindDetails) (*Credentials, error)
	DeleteCredentials(issuerID string, id string) (bool, error)
}

type Credentials struct {
	URL          string `json:"url"`
	ClientID     string `json:"clientid"`
	ClientSecret string `json:"clientsecret"`
}

type credentialsProviderClient struct {
	config *Config
	HTTPClient
}

func NewClient(config *Config, securedClient *http.Client) Client {
	return &credentialsProviderClient{
		config:     config,
		HTTPClient: securedClient,
	}
}

// CreateCredentialsIssuer creates a new credentials issuer for the specified instance ID
func (c *credentialsProviderClient) CreateCredentialsIssuer(instanceID string, provisionDetails domain.ProvisionDetails) (bool, error) {
	tenantID := gjson.GetBytes(provisionDetails.RawContext, c.config.TenantIDContextKey)
	if !tenantID.Exists() {
		return false, fmt.Errorf("failed to get tenant id from osb context")
	}

	queryParams := map[string]string{
		c.config.IssuerNameQueryParam: instanceID,
		c.config.TenantIDQueryParam:   tenantID.Str,
	}

	resp, err := c.sendRequest(http.MethodPost, c.config.CredentialsIssuerPath, []byte(fmt.Sprintf(c.config.IssuerBody, instanceID)), queryParams)
	if err != nil {
		return false, fmt.Errorf("failed to create credentials issuer: %v", err)
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusCreated:
		return true, nil
	case http.StatusConflict:
		return false, nil
	default:
		return false, handleUnexpectedResponse(resp)
	}
}

// DeleteCredentialsIssuer deletes a credentials issuer for the specified instance ID
func (c *credentialsProviderClient) DeleteCredentialsIssuer(instanceID string) (bool, error) {
	return c.sendDeleteRequest(c.config.CredentialsIssuerPath + "/" + url.PathEscape(instanceID))
}

func (c *credentialsProviderClient) CreateCredentials(instanceID, bindingID string, details domain.BindDetails) (*Credentials, error) {
	body, err := buildCreateCredentialsBody(details)
	if err != nil {
		return nil, fmt.Errorf("failed to build create credentials body: %v", err)
	}

	endpoint := fmt.Sprintf(c.config.CredentialsPath, url.PathEscape(instanceID), url.PathEscape(bindingID))
	resp, err := c.sendRequest(http.MethodPut, endpoint, body, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to execute create credentials request: %v", err)
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		var cred Credentials
		decoder := json.NewDecoder(resp.Body)
		if err = decoder.Decode(&cred); err != nil {
			return nil, fmt.Errorf("failed to process create credentials response body: %v", err)
		}
		return &cred, nil
	default:
		return nil, handleUnexpectedResponse(resp)
	}
}

// DeleteBinding deletes the binding for the specified instance ID and binding ID
func (c *credentialsProviderClient) DeleteCredentials(instanceID, bindingID string) (bool, error) {
	endpoint := fmt.Sprintf(c.config.CredentialsPath, url.PathEscape(instanceID), url.PathEscape(bindingID))
	return c.sendDeleteRequest(endpoint)
}

func (c *credentialsProviderClient) sendDeleteRequest(endpoint string) (bool, error) {
	resp, err := c.sendRequest(http.MethodDelete, endpoint, nil, nil)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		return true, nil
	case http.StatusNotFound:
		return false, nil
	default:
		return false, handleUnexpectedResponse(resp)
	}
}

func (c *credentialsProviderClient) sendRequest(method, endpoint string, body []byte, queryParams map[string]string) (*http.Response, error) {
	req, err := http.NewRequest(method, c.config.URL+endpoint, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	if len(queryParams) > 0 {
		query := req.URL.Query()
		for k, v := range queryParams {
			query.Add(k, v)
		}
		req.URL.RawQuery = query.Encode()
	}
	return c.Do(req)
}

func buildCreateCredentialsBody(details domain.BindDetails) ([]byte, error) {
	appGUID := extractAppGuid(details)
	body := struct {
		AppGUID       string          `json:"app_guid"`
		RawParameters json.RawMessage `json:"parameters,omitempty"`
	}{
		AppGUID:       appGUID,
		RawParameters: details.GetRawParameters(),
	}

	return json.Marshal(body)
}

func extractAppGuid(details domain.BindDetails) string {
	if len(details.AppGUID) != 0 {
		return details.AppGUID
	}

	if details.BindResource != nil && len(details.BindResource.AppGuid) != 0 {
		return details.BindResource.AppGuid
	}

	// For service keys app_guid is expected to be empty.
	return ""
}

func handleUnexpectedResponse(resp *http.Response) error {
	responseBody, _ := ioutil.ReadAll(resp.Body)
	return fmt.Errorf("credentials provider responded with status code %d and body %s", resp.StatusCode, string(responseBody))
}
