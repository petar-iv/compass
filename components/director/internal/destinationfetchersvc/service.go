package destinationfetchersvc

import (
	"bytes"
	"context"
	"fmt"
	"strconv"

	"github.com/kyma-incubator/compass/components/director/internal/model"
	"github.com/kyma-incubator/compass/components/director/pkg/apperrors"
	"github.com/kyma-incubator/compass/components/director/pkg/log"
	"github.com/kyma-incubator/compass/components/director/pkg/persistence"
	"github.com/kyma-incubator/compass/components/director/pkg/resource"

	"github.com/pkg/errors"
)

const (
	correlationIDPrefix = "sap.s4:communicationScenario:"
	regionLabelKey      = "region"
)

type UUIDService interface {
	Generate() string
}

type DestinationRepo interface {
	Upsert(ctx context.Context, in model.DestinationInput, id, tenantID, bundleID, revision string) error
	Delete(ctx context.Context, revision string) error
}

type LabelRepo interface {
	GetSubdomainLabelForSubscribedRuntime(ctx context.Context, subaccountId string) (*model.Label, error)
	GetByKey(ctx context.Context, tenant string, objectType model.LabelableObject, objectID, key string) (*model.Label, error)
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
	label, err := d.getSubscribedSubdomainLabel(ctx, subaccountID)
	if err != nil {
		return err
	}

	region, err := d.getRegionLabel(ctx, *label.Tenant)
	if err != nil {
		return err
	}

	fmt.Println(region.Value.(string))
	// oauthConfig := oautconfigs[region]

	subdomain := label.Value.(string)
	client, err := NewClient(d.oauthConfig, d.apiConfig, subdomain)
	if err != nil {
		log.C(ctx).WithError(err).Errorf("Failed to create DestinationInput API client: %v", err)
		return err
	}

	if err := d.walkthroughPages(client, func(destinations []model.DestinationInput) error {
		log.C(ctx).Infof("Found %d destinations in subaccount '%s'", len(destinations), subaccountID)

		tx, err := d.transact.Begin()
		if err != nil {
			log.C(ctx).WithError(err).Errorf("Failed to begin db transaction: %v")
			return err
		}
		ctx = persistence.SaveToContext(ctx, tx)
		defer d.transact.RollbackUnlessCommitted(ctx, tx)

		for _, destination := range destinations {
			correlationID := correlationIDPrefix + destination.CommunicationScenarioId
			bundles, err := d.bundleRepo.GetBySystemAndCorrelationId(ctx, *label.Tenant, destination.XFSystemName, destination.URL, correlationID)

			if len(bundles) == 0 {
				log.C(ctx).Infof("No bundles found for system '%s', url '%s', correlation id '%s'", destination.XFSystemName, destination.URL, correlationID)
				continue
			}

			if err != nil {
				log.C(ctx).WithError(err).Errorf("Failed to fetch bundle for system '%s', url '%s', correlation id '%s', tenant id '%s': %v", destination.XFSystemName, destination.URL, correlationID, *label.Tenant, err)
				continue
			}

			for _, bundle := range bundles {
				id := d.uuidSvc.Generate()
				revision := d.uuidSvc.Generate()
				if err := d.repo.Upsert(ctx, destination, id, *label.Tenant, bundle.ID, revision); err != nil {
					log.C(ctx).WithError(err).Errorf("Failed to insert destination with name '%s' for bunlde '%s' and tenant '%s' to DB: %v", destination.Name, bundle.ID, *label.Tenant, err)
					continue
				}
			}
		}

		if err = tx.Commit(); err != nil {
			log.C(ctx).WithError(err).Errorf("Failed to commit database transaction %v", err)
			return err
		}
		return nil
	}); err != nil {
		log.C(ctx).WithError(err).Errorf("Failed to sync destinations for subaccount '%s': %v", subaccountID, err)
		return err
	}

	return nil
}

type processFunc func([]model.DestinationInput) error

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
	//subdomain := "i305674-4"

	fmt.Println(subaccountID)
	tx, err := d.transact.Begin()
	ctx = persistence.SaveToContext(ctx, tx)
	if err != nil {
		return nil, err
	}

	defer d.transact.RollbackUnlessCommitted(ctx, tx)
	log.C(ctx).Infof("Getting subdomain name\n")
	label, err := d.labelRepo.GetSubdomainLabelForSubscribedRuntime(ctx, subaccountID)
	if err != nil {
		return nil, err
	}

	subdomain := label.Value.(string)
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

func (d DestinationService) getSubscribedSubdomainLabel(ctx context.Context, subaccountID string) (*model.Label, error) {
	tx, err := d.transact.Begin()
	if err != nil {
		log.C(ctx).WithError(err).Errorf("Failed to begin db transaction: %v")
		return nil, err
	}
	ctx = persistence.SaveToContext(ctx, tx)
	defer d.transact.RollbackUnlessCommitted(ctx, tx)

	label, err := d.labelRepo.GetSubdomainLabelForSubscribedRuntime(ctx, subaccountID)
	if err != nil {
		if apperrors.IsNotFoundError(err) {
			log.C(ctx).Errorf("No subscribed subdomain found for subbaccount '%s'", subaccountID)
			return nil, apperrors.NewNotFoundErrorWithMessage(resource.Label, subaccountID, fmt.Sprintf("subaccount %s not found", subaccountID))
		}
		log.C(ctx).WithError(err).Errorf("Failed to get subdomain for subaccount '%s' from db: %v", subaccountID, err)
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		log.C(ctx).WithError(err).Errorf("Failed to commit database transaction %v", err)
		return nil, err
	}

	return label, nil
}

func (d DestinationService) getRegionLabel(ctx context.Context, tenantID string) (*model.Label, error) {
	tx, err := d.transact.Begin()
	if err != nil {
		log.C(ctx).WithError(err).Errorf("Failed to begin db transaction: %v")
		return nil, err
	}
	ctx = persistence.SaveToContext(ctx, tx)
	defer d.transact.RollbackUnlessCommitted(ctx, tx)

	region, err := d.labelRepo.GetByKey(ctx, tenantID, model.TenantLabelableObject, tenantID, regionLabelKey)
	if err != nil {
		log.C(ctx).WithError(err).Errorf("Failed to fetch region for tenant '%s': %v", tenantID, err)
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		log.C(ctx).WithError(err).Errorf("Failed to commit database transaction %v", err)
		return nil, err
	}

	return region, nil
}
