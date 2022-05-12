package fetcher_client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/pkg/errors"
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
		return errors.Wrap(err, "while making HTTP request")
	}

	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Info("Got error on closing request body")
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return errors.New(fmt.Sprintf("received unexpecte status code %d while making request to tenant fetcher", resp.StatusCode))
	}
	return nil
}
