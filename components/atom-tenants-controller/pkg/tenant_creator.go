package pkg

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"

	rmlogger "github.tools.sap/unified-resource-manager/api/pkg/apis/logger"
)

type TenantCreator interface {
	StoreTenants(ctx context.Context, payload RequestPayload) error
}

type service struct {
	client           *http.Client
	tenantFetcherURL string
}

func NewCreator(client *http.Client, tenantFetcherURL string) *service {
	return &service{
		client:           client,
		tenantFetcherURL: tenantFetcherURL,
	}
}

func (s *service) StoreTenants(ctx context.Context, payload RequestPayload) error {
	log := rmlogger.FromContext(ctx, "tenants-aggregator")

	jsonData, err := json.Marshal(payload)
	req, err := http.NewRequest(http.MethodPost, s.tenantFetcherURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return err
	}

	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Info("Got error on closing request body")
		}
	}()

	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	log.Info(string(responseBody))

	return nil
}
