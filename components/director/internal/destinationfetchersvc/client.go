package destinationfetchersvc

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/kyma-incubator/compass/components/director/internal/model"
	"github.com/kyma-incubator/compass/components/director/pkg/apperrors"
	"github.com/kyma-incubator/compass/components/director/pkg/config"
	"github.com/kyma-incubator/compass/components/director/pkg/oauth"
	"github.com/kyma-incubator/compass/components/director/pkg/resource"
	"github.com/pkg/errors"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
)

type OAuth2Config struct {
	ClientID           string `envconfig:"APP_CLIENT_ID"`
	OAuthTokenEndpoint string `envconfig:"APP_OAUTH_TOKEN_ENDPOINT"`
	SkipSSLValidation  bool   `envconfig:"APP_OAUTH_SKIP_SSL_VALIDATION,default=false"`
	X509Config         oauth.X509Config
}

type APIConfig struct {
	//TODO optional?
	EndpointGetSubbacountDestinations string        `envconfig:"APP_ENDPOINT_GET_SUBACCOUNT_DESTINATIONS"`
	EndpointFindDestination           string        `envconfig:"APP_ENDPOINT_FIND_DESTINATION"`
	Timeout                           time.Duration `envconfig:"APP_DESTINATIONS_TIMEOUT"`
	PageSize                          int           `envconfig:"APP_DESTINATIONS_PAGE_SIZE"`
	PagingPageParam                   string        `envconfig:"APP_DESTINATIONS_PAGE_PARAM"`
	PagingSizeParam                   string        `envconfig:"APP_DESTINATIONS_PAGE_SIZE_PARAM"`
	PagingCountParam                  string        `envconfig:"APP_DESTINATIONS_PAGE_COUNT_PARAM"`
	PagingCountHeader                 string        `envconfig:"APP_DESTINATIONS_PAGE_COUNT_HEADER"`
}

type Client struct {
	httpClient *http.Client
	apiConfig  APIConfig
}

type DestinationResponse struct {
	destinations []model.DestinationInput
	pageCount    string
}

func NewClient(instanceConfig config.InstanceConfig, apiConfig APIConfig, subdomain string) (*Client, error) {
	ctx := context.Background()

	u, err := url.Parse(instanceConfig.TokenURL)
	if err != nil {
		return nil, errors.Errorf("failed to parse auth url '%s': %v", instanceConfig.TokenURL, err)
	}
	parts := strings.Split(u.Hostname(), ".")
	originalSubdomain := parts[0]

	// TODO use env variable for token path
	tokenURL := strings.Replace(instanceConfig.TokenURL, originalSubdomain, subdomain, 1) + "/oauth/token"
	cfg := clientcredentials.Config{
		ClientID:  instanceConfig.ClientID,
		TokenURL:  tokenURL,
		AuthStyle: oauth2.AuthStyleInParams,
	}
	cert, err := tls.X509KeyPair([]byte(instanceConfig.Cert), []byte(instanceConfig.Key))
	if err != nil {
		return nil, errors.Errorf("failed to create destinations client x509 pair: %v", err)
	}

	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			Certificates: []tls.Certificate{cert},
		},
	}

	mtlClient := &http.Client{
		Transport: transport,
		Timeout:   apiConfig.Timeout,
	}

	ctx = context.WithValue(ctx, oauth2.HTTPClient, mtlClient)

	httpClient := cfg.Client(ctx)
	httpClient.Timeout = apiConfig.Timeout

	return &Client{
		httpClient: httpClient,
		apiConfig:  apiConfig,
	}, nil
}

func (c *Client) FetchSubbacountDestinationsPage(page string) (*DestinationResponse, error) {
	req, err := c.buildRequest(c.apiConfig.EndpointGetSubbacountDestinations, page)
	if err != nil {
		return nil, err
	}

	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "failed to execute HTTP request")
	}

	var destinations []model.DestinationInput
	if err := json.NewDecoder(res.Body).Decode(&destinations); err != nil {
		return nil, errors.Wrap(err, "failed to decode response body")
	}

	if res.StatusCode != http.StatusOK {
		return nil, errors.Errorf("received status code %d when trying to fetch destinations", res.StatusCode)
	}

	pageCount := res.Header.Get(c.apiConfig.PagingCountHeader)
	if pageCount == "" {
		return nil, errors.Wrapf(err, "failed to extract header '%s' from destinations response", c.apiConfig.PagingCountParam)
	}

	return &DestinationResponse{
		destinations: destinations,
		pageCount:    pageCount,
	}, nil
}

func (c *Client) buildRequest(url string, page string) (*http.Request, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to build request")
	}

	query := req.URL.Query()
	query.Add(c.apiConfig.PagingCountParam, "true")
	query.Add(c.apiConfig.PagingPageParam, page)
	query.Add(c.apiConfig.PagingSizeParam, strconv.Itoa(c.apiConfig.PageSize))
	req.URL.RawQuery = query.Encode()
	return req, nil
}

func (c *Client) fetchDestinationSensitiveData(destinationName string) ([]byte, error) {
	req, err := http.NewRequest(http.MethodGet, c.apiConfig.EndpointFindDestination+destinationName, nil)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to build request")
	}

	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "failed to execute HTTP request")
	}

	if res.StatusCode == http.StatusNotFound {
		return nil, apperrors.NewNotFoundErrorWithMessage(resource.Destination, destinationName, "")
	}

	if res.StatusCode != http.StatusOK {
		return nil, errors.Errorf("received status code %d when trying to get destination info for %s",
			res.StatusCode, destinationName)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read body of response")
	}

	return body, nil
}
