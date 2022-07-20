package destinationfetchersvc

import (
	"context"
	"fmt"

	"github.com/kyma-incubator/compass/components/director/pkg/persistence"
)

type DestinationRepo interface {
	Upsert(ctx context.Context) error
	Delete(ctx context.Context) error
}

type DestinationService struct {
	transact persistence.Transactioner
	repo     DestinationRepo
}

func NewDestinationService(transact persistence.Transactioner, destinationRepo DestinationRepo) *DestinationService {
	return &DestinationService{
		transact: transact,
		repo:     destinationRepo,
	}
}

func (d DestinationService) SyncDestinations() error {
	ctx := context.Background()

	tx, err := d.transact.Begin()
	if err != nil {
		return err
	}
	defer d.transact.RollbackUnlessCommitted(ctx, tx)
	ctx = persistence.SaveToContext(ctx, tx)

	d.repo.Delete(ctx)

	if err = tx.Commit(); err != nil {
		return err
	}
	return nil
}
