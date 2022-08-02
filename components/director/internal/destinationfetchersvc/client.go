package destinationfetchersvc

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/kyma-incubator/compass/components/director/internal/model"
	"github.com/kyma-incubator/compass/components/director/pkg/apperrors"
	"github.com/kyma-incubator/compass/components/director/pkg/config"
	"github.com/kyma-incubator/compass/components/director/pkg/resource"
	"github.com/pkg/errors"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
)

type APIConfig struct {
	//TODO optional, default?
	GoroutineLimit                    int64         `envconfig:"APP_DESTINATIONS_SENSITIVE_GOROUTINE_LIMIT"`
	RetryInterval                     int           `envconfig:"APP_DESTINATIONS_RETRY_INTERVAL"`
	RetryLimit                        int           `envconfig:"APP_DESTINATIONS_RETRY_LIMIT"`
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
	apiURL     string
}

type DestinationResponse struct {
	destinations []model.DestinationInput
	pageCount    string
}

func NewClient(instanceConfig config.InstanceConfig, apiConfig APIConfig, tokenPath, subdomain string) (*Client, error) {
	ctx := context.Background()

	u, err := url.Parse(instanceConfig.TokenURL)
	if err != nil {
		return nil, errors.Errorf("failed to parse auth url '%s': %v", instanceConfig.TokenURL, err)
	}
	parts := strings.Split(u.Hostname(), ".")
	originalSubdomain := parts[0]

	tokenURL := strings.Replace(instanceConfig.TokenURL, originalSubdomain, subdomain, 1) + tokenPath
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
		apiURL:     instanceConfig.URL,
	}, nil
}

func (c *Client) FetchSubbacountDestinationsPage(page string) (*DestinationResponse, error) {
	url := c.apiURL + c.apiConfig.EndpointGetSubbacountDestinations
	req, err := c.buildRequest(url, page)
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

func (c *Client) FetchDestinationSensitiveData(destinationName string) ([]byte, error) {
	url := fmt.Sprintf("%s%s/%s", c.apiURL, c.apiConfig.EndpointFindDestination, destinationName)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to build request")
	}

	res, err := c.sendRequestWithRetry(req)
	if err != nil {
		return nil, err
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

func (c *Client) sendRequestWithRetry(req *http.Request) (*http.Response, error) {
	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "failed to execute HTTP request")
	}

	i := 0
	for i < c.apiConfig.RetryLimit && res.StatusCode == http.StatusInternalServerError {
		time.Sleep(time.Millisecond * time.Duration(c.apiConfig.RetryInterval))
		i++
		res, err = c.httpClient.Do(req)
		if err != nil {
			return nil, errors.Wrap(err, "failed to execute HTTP request")
		}
	}

	return res, nil
}
