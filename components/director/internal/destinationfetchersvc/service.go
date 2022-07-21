package destinationfetchersvc

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/kyma-incubator/compass/components/director/internal/model"
	"github.com/kyma-incubator/compass/components/director/pkg/persistence"
)

type DestinationRepo interface {
	Upsert(ctx context.Context) error
	Delete(ctx context.Context) error
}

type DestinationService struct {
	transact    persistence.Transactioner
	repo        DestinationRepo
	oauthConfig OAuth2Config
	apiConfig   APIConfig
}

type DestinationAPIClient interface {
	FetchSubbacountDestinationsPage(page string) (*DestinationResponse, error)
}

func NewDestinationService(transact persistence.Transactioner, destinationRepo DestinationRepo, oauthConfig OAuth2Config, apiConfig APIConfig) *DestinationService {
	return &DestinationService{
		transact:    transact,
		repo:        destinationRepo,
		oauthConfig: oauthConfig,
		apiConfig:   apiConfig,
	}
}

func (d DestinationService) SyncSubaccountDestinations(subaccountID string) error {
	//TODO get subdomain from subaccountID
	subdomain := "i331217-provider"

	client, err := NewClient(d.oauthConfig, d.apiConfig, subdomain)
	if err != nil {
		return err
	}

	if err := d.walkthroughPages(client, func(destinations []model.Destination) error {
		log.Printf("found %d destinations in subaccount %s", len(destinations), subaccountID)
		for _, destination := range destinations {
			fmt.Println(destination)
		}
		return nil
	}); err != nil {
		return err
	}

	// ctx := context.Background()

	// tx, err := d.transact.Begin()
	// if err != nil {
	// 	return err
	// }
	// defer d.transact.RollbackUnlessCommitted(ctx, tx)
	// ctx = persistence.SaveToContext(ctx, tx)

	// d.repo.Delete(ctx)

	// if err = tx.Commit(); err != nil {
	// 	return err
	// }
	return nil
}

type processFunc func([]model.Destination) error

func (d DestinationService) walkthroughPages(client *Client, process processFunc) error {
	hasMorePages := true

	for page := 1; hasMorePages; page++ {
		pageString := strconv.Itoa(page)
		resp, err := client.FetchSubbacountDestinationsPage(pageString)
		fmt.Println(err)
		if err != nil {
			return err
		}

		if err := process(resp.destinations); err != nil {
			return fmt.Errorf("failed to process destinations page: %w", err)
		}

		hasMorePages = pageString != resp.pageCount
	}

	return nil
}
