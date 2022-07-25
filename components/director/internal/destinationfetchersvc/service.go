package destinationfetchersvc

import (
	"context"
	"fmt"
	"strconv"

	"github.com/kyma-incubator/compass/components/director/internal/domain/application"
	"github.com/kyma-incubator/compass/components/director/internal/domain/bundle"
	domain "github.com/kyma-incubator/compass/components/director/internal/domain/destination"
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
	GetSubdomains(ctx context.Context) ([]domain.Subdomain, error)
}

type DestinationService struct {
	transact    persistence.Transactioner
	uuidSvc     UUIDService
	repo        DestinationRepo
	bundleRepo  bundle.BundleRepository
	appRepo     application.ApplicationRepository
	oauthConfig OAuth2Config
	apiConfig   APIConfig
}

type DestinationAPIClient interface {
	FetchSubbacountDestinationsPage(page string) (*DestinationResponse, error)
}

func NewDestinationService(transact persistence.Transactioner, uuidSvc UUIDService, destRepo DestinationRepo, bundleRepo bundle.BundleRepository, appRepo application.ApplicationRepository, oauthConfig OAuth2Config, apiConfig APIConfig) *DestinationService {
	return &DestinationService{
		transact:    transact,
		uuidSvc:     uuidSvc,
		repo:        destRepo,
		bundleRepo:  bundleRepo,
		appRepo:     appRepo,
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
	bundle, err := d.bundleRepo.GetBySystemAndCorrelationId(ctx, "app-0ysx", "https://commerce.c-9a3de5f.kyma-dev.shoot.canary.k8s-hana.ondemand.com", "correlationId")
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println(bundle)

	if err = tx.Commit(); err != nil {
		return errors.Wrap(err, "failed to commit database transaction")
	}

	// // TODO Should we check explicitly if UCL is subscribed in the given subaccount?
	// subdomain, err := d.repo.GetSubdomain(ctx, subaccountID)
	// if err != nil {
	// 	return err
	// }

	// // TODO Should return 400 Bad Request?
	// if subdomain == nil {
	// 	return errors.New(fmt.Sprintf("subdomain for subaccount with id '%s' doesn't exist", subaccountID))
	// }

	// client, err := NewClient(d.oauthConfig, d.apiConfig, subdomain.Value)
	// if err != nil {
	// 	return errors.Wrap(err, "failed to create destinations API client")
	// }

	// tx, err := d.transact.Begin()
	// ctx = persistence.SaveToContext(ctx, tx)
	// if err != nil {
	// 	return err
	// }
	// defer d.transact.RollbackUnlessCommitted(ctx, tx)

	// if err := d.walkthroughPages(client, func(destinations []Destination) error {
	// 	log.C(ctx).Infof("Found %d destinations in subaccount '%s'", len(destinations), subaccountID)
	// 	for _, destination := range destinations {
	// 		correlationID := correlationIDPrefix + destination.CommunicationScenarioId
	// 		bundles, err := d.repo.GetBundleForDestination(ctx, destination.XFSystemName, destination.URL, correlationID)

	// 		if err != nil {
	// 			return err
	// 		}
	// 		if len(bundles) == 0 {
	// 			log.C(ctx).Debugf("Bundles for system with name: '%s', url: '%s' and correlation id: '%s' not found", destination.XFSystemName, destination.URL, correlationID)
	// 			continue
	// 		}

	// 		for _, bundle := range bundles {
	// 			destinationDB := domain.Entity{
	// 				ID:             d.uuidSvc.Generate(),
	// 				Name:           destination.Name,
	// 				Type:           destination.Type,
	// 				URL:            destination.URL,
	// 				Authentication: destination.Authentication,
	// 				BundleID:       bundle.ID,
	// 				TenantID:       subdomain.TenantID,
	// 				Revision:       d.uuidSvc.Generate(),
	// 			}
	// 			if err := d.repo.Upsert(ctx); err != nil {
	// 				return errors.Wrapf(err, "failed to insert destination data '%+v' to DB: %w", destinationDB)
	// 			}
	// 		}
	// 	}
	// 	return nil
	// }); err != nil {
	// 	return err
	// }

	// if err = tx.Commit(); err != nil {
	// 	return errors.Wrap(err, "failed to commit database transaction")
	// }

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
