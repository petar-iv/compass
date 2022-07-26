package destinationfetchersvc

import (
	"bytes"
	"context"
	"fmt"
	"strconv"

	domain "github.com/kyma-incubator/compass/components/director/internal/domain/destination"
	"github.com/kyma-incubator/compass/components/director/internal/model"
	"github.com/kyma-incubator/compass/components/director/pkg/apperrors"
	"github.com/kyma-incubator/compass/components/director/pkg/log"
	"github.com/kyma-incubator/compass/components/director/pkg/persistence"

	"github.com/pkg/errors"
)

const (
	correlationIDPrefix = "sap.s4:communicationScenario:"
)

type UUIDService interface {
	Generate() string
}

type DestinationRepo interface {
	Upsert(ctx context.Context) error
	Delete(ctx context.Context, revision string) error
	GetSubdomain(ctx context.Context, subaccountId string) (*domain.Subdomain, error)
}

type LabelRepo interface {
	ListSubdomainLabelsForRuntimes(ctx context.Context) ([]*model.Label, error)
	GetSubdomainLabelForRuntime(ctx context.Context, subaccountId string) (*model.Label, error)
}

type BundleRepo interface {
	GetBySystemAndCorrelationId(ctx context.Context, tenantId, systemName, systemURL, correlationId string) ([]*model.Bundle, error)
}

type DestinationService struct {
	transact    persistence.Transactioner
	uuidSvc     UUIDService
	repo        DestinationRepo
	bundleRepo  BundleRepo
	labelRepo   LabelRepo
	oauthConfig OAuth2Config
	apiConfig   APIConfig
}

type DestinationAPIClient interface {
	FetchSubbacountDestinationsPage(page string) (*DestinationResponse, error)
}

func NewDestinationService(transact persistence.Transactioner, uuidSvc UUIDService, destRepo DestinationRepo, bundleRepo BundleRepo, labelRepo LabelRepo, oauthConfig OAuth2Config, apiConfig APIConfig) *DestinationService {
	return &DestinationService{
		transact:    transact,
		uuidSvc:     uuidSvc,
		repo:        destRepo,
		bundleRepo:  bundleRepo,
		labelRepo:   labelRepo,
		oauthConfig: oauthConfig,
		apiConfig:   apiConfig,
	}
}

func (d DestinationService) SyncSubaccountDestinations(ctx context.Context, subaccountID string) error {
	tx, err := d.transact.Begin()
	ctx = persistence.SaveToContext(ctx, tx)
	if err != nil {
		return err
	}
	defer d.transact.RollbackUnlessCommitted(ctx, tx)

	label, err := d.labelRepo.GetSubdomainLabelForRuntime(ctx, subaccountID)
	if err != nil {
		return err
	}

	client, err := NewClient(d.oauthConfig, d.apiConfig, label.Value.(string))
	if err != nil {
		return errors.Wrap(err, "failed to create destinations API client")
	}

	if err := d.walkthroughPages(client, func(destinations []Destination) error {
		log.C(ctx).Infof("Found %d destinations in subaccount '%s'", len(destinations), subaccountID)
		for _, destination := range destinations {
			correlationID := correlationIDPrefix + destination.CommunicationScenarioId
			bundles, err := d.bundleRepo.GetBySystemAndCorrelationId(ctx, *label.Tenant, destination.XFSystemName, destination.URL, correlationID)

			if apperrors.IsNotFoundError(err) {
				log.C(ctx).Infof("No bundle found for system '%s', url '%s', correlation id '%s'. Will skip this destination ...", destination.XFSystemName, destination.URL, correlationID)
				continue
			}

			if err != nil {
				log.C(ctx).Errorf("Failed to fetch bundle for system '%s', url '%s', correlation id '%s': %v", destination.XFSystemName, destination.URL, correlationID, err)
				continue
			}

			for _, bundle := range bundles {
				destinationDB := domain.Entity{
					ID:             d.uuidSvc.Generate(),
					Name:           destination.Name,
					Type:           destination.Type,
					URL:            destination.URL,
					Authentication: destination.Authentication,
					BundleID:       bundle.ID,
					TenantID:       *label.Tenant,
					Revision:       d.uuidSvc.Generate(),
				}

\				if err := d.repo.Upsert(ctx); err != nil {
					return errors.Wrapf(err, "failed to insert destination data '%+v' to DB: %w", destinationDB)
				}
			}
		}
		return nil
	}); err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return errors.Wrap(err, "failed to commit database transaction")
	}

	return nil
}

type processFunc func([]Destination) error

func (d DestinationService) walkthroughPages(client *Client, process processFunc) error {
	hasMorePages := true

	for page := 1; hasMorePages; page++ {
		pageString := strconv.Itoa(page)
		resp, err := client.FetchSubbacountDestinationsPage(pageString)
		if err != nil {
			return errors.Wrap(err, "failed to fetch destinations page")
		}

		if err := process(resp.destinations); err != nil {
			return errors.Wrap(err, "failed to process destinations page")
		}

		hasMorePages = pageString != resp.pageCount
	}

	return nil
}

func (d DestinationService) FetchDestinationsSensitiveData(ctx context.Context, subaccountID string, destinationNames []string) ([]byte, error) {
	subdomain := "i305674-4"
	log.C(ctx).Infof("Fetching data for subdomain: %s \n", subdomain)

	client, err := NewClient(d.oauthConfig, d.apiConfig, subdomain)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create destinations API client")
	}

	results := make([][]byte, len(destinationNames))
	for i, destination := range destinationNames {
		log.C(ctx).Infof("Fetching data for destination: %s \n", destination)
		results[i], err = client.fetchDestinationSensitiveData(destination)
		if err != nil {
			return nil, err
		}
	}

	combinedInfoJSON := bytes.Join(results, []byte(","))
	combinedInfoJSON = append(combinedInfoJSON, ']', '}')

	return append([]byte("{ \"destinations\": ["), combinedInfoJSON...), nil
}
