package open_resource_discovery

import (
	"context"
	"encoding/json"
	"github.com/kyma-incubator/compass/components/director/internal/model"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/kyma-incubator/compass/components/director/pkg/log"
	"github.com/pkg/errors"
)

// Client represents ORD documents client
//go:generate mockery --name=Client --output=automock --outpkg=automock --case=underscore
type Client interface {
	FetchOpenResourceDiscoveryDocuments(ctx context.Context, webhook *model.Webhook) (Documents, string, error)
}

type client struct {
	*http.Client
}

// NewClient creates new ORD Client via a provided http.Client
func NewClient(httpClient *http.Client) *client {
	return &client{
		Client: httpClient,
	}
}

// FetchOpenResourceDiscoveryDocuments fetches all the documents for a single ORD .well-known endpoint
func (c *client) FetchOpenResourceDiscoveryDocuments(ctx context.Context, webhook *model.Webhook) (Documents, string, error) {
	config, err := c.fetchConfig(ctx, *webhook.URL)
	if err != nil {
		return nil, "", err
	}

	if webhook.ProxyURL != nil && *webhook.ProxyURL != "" { // TODO: In productive implementation this should be at the very start of the FetchOpenResourceDiscoveryDocuments function
		if err := c.setProxy(ctx, *webhook.ProxyURL); err != nil {
			return Documents{}, "", err
		}
		defer func() {
			if err := c.removeProxy(); err != nil {
				log.C(ctx).Errorf("Error occurred while reverting proxy transport configuration: %v", err.Error())
			}
		}()
	}

	docs := make([]*Document, 0, 0)
	actualBaseURL := *webhook.URL // TODO: Workaround due to provider/described system mismatch...
	for _, docDetails := range config.OpenResourceDiscoveryV1.Documents {
		documentURL := *webhook.URL + docDetails.URL
		u, err := url.ParseRequestURI(docDetails.URL)
		if err == nil {
			documentURL = docDetails.URL
			actualBaseURL = u.Scheme + "://" + u.Hostname() + ":" + u.Port()
		}

		strategy, ok := docDetails.AccessStrategies.GetSupported()
		if !ok {
			log.C(ctx).Warnf("Unsupported access strategies for ORD Document %q", documentURL)
			continue
		}

		doc, err := c.fetchOpenDiscoveryDocumentWithAccessStrategy(ctx, webhook, documentURL, strategy)
		if err != nil {
			return nil, "", errors.Wrapf(err, "error fetching ORD document from: %s", documentURL)
		}

		docs = append(docs, doc)
	}

	return docs, actualBaseURL, nil
}

func (c *client) setProxy(ctx context.Context, proxyURL string) error {
	proxyUrl, err := url.Parse(proxyURL)
	if err != nil {
		log.C(ctx).WithError(err).Warnf("Got error parsing proxy url: %s", proxyUrl)
		return err
	}

	transport := c.Client.Transport.(*http.Transport)
	transport.Proxy = http.ProxyURL(proxyUrl)

	return nil
}

func (c *client) removeProxy() error {
	transport := c.Client.Transport.(*http.Transport)
	transport.Proxy = nil

	return nil
}

func (c *client) fetchOpenDiscoveryDocumentWithAccessStrategy(ctx context.Context, webhook *model.Webhook, documentURL string, accessStrategy AccessStrategyType) (*Document, error) {
	parsedURL, err := url.Parse(documentURL)
	if err != nil {
		return nil, err
	}

	httpRequest := &http.Request{URL: parsedURL, Header: map[string][]string{}}
	if accessStrategy == BasicAccessStrategy {
		httpRequest.SetBasicAuth(webhook.Auth.Credential.Basic.Username, webhook.Auth.Credential.Basic.Password)
	}

	log.C(ctx).Infof("Fetching ORD Document %q", documentURL)
	resp, err := c.Do(httpRequest)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.Errorf("error while fetching open resource discovery document %q: status code %d", documentURL, resp.StatusCode)
	}

	defer closeBody(ctx, resp.Body)

	resp.Body = http.MaxBytesReader(nil, resp.Body, 2097152)
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "error reading document body")
	}
	result := &Document{}
	if err := json.Unmarshal(bodyBytes, &result); err != nil {
		return nil, errors.Wrap(err, "error unmarshaling document")
	}
	return result, nil
}

func closeBody(ctx context.Context, body io.ReadCloser) {
	err := body.Close()
	if err != nil {
		log.C(ctx).WithError(err).Warnf("Got error on closing response body")
	}
}

func (c *client) fetchConfig(ctx context.Context, url string) (*WellKnownConfig, error) {
	resp, err := c.Get(url + WellKnownEndpoint)
	if err != nil {
		return nil, errors.Wrap(err, "error while fetching open resource discovery well-known configuration")
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.Errorf("error while fetching open resource discovery well-known configuration: status code %d", resp.StatusCode)
	}

	defer closeBody(ctx, resp.Body)

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "error reading response body")
	}

	config := WellKnownConfig{}
	if err := json.Unmarshal(bodyBytes, &config); err != nil {
		return nil, errors.Wrap(err, "error unmarshaling json body")
	}

	return &config, nil
}
