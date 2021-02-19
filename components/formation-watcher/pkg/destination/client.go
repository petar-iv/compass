package destination

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/kyma-incubator/compass/components/formation-watcher/pkg/log"
)

type Client struct {
	httpClient *http.Client
	baseURL    string
}

func NewClient(baseURL string, httpClient *http.Client) *Client {
	return &Client{
		httpClient: httpClient,
		baseURL:    baseURL,
	}
}

func (c *Client) GetAllDestinations(ctx context.Context) ([]Destination, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s/%s", c.baseURL, "subaccountDestinations"), nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Unexpected status code returned by destination service %d", resp.StatusCode)
	}
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var result []Destination
	if err := json.Unmarshal(bytes, &result); err != nil {
		return nil, err
	}

	return result, nil
}

func (c *Client) CreateDestination(ctx context.Context, dest Destination) error {
	bytesBody, err := json.Marshal(dest)
	if err != nil {
		return err
	}
	body := bytes.NewBuffer(bytesBody)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, fmt.Sprintf("%s/%s", c.baseURL, "subaccountDestinations"), body)
	if err != nil {
		return err
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusCreated {
		respBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		log.C(ctx).Errorf("Request failed with %s", string(respBytes))
		return fmt.Errorf("Unexpected status code returned by destination service %d", resp.StatusCode)
	}

	return nil
}

func (c *Client) DeleteDestination(ctx context.Context, name string) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, fmt.Sprintf("%s/%s/%s", c.baseURL, "subaccountDestinations", name), nil)
	if err != nil {
		return err
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Unexpected status code returned by destination service %d", resp.StatusCode)
	}

	return nil
}
