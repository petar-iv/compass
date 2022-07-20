package destinationfetchersvc

import (
	"fmt"

	"github.com/kyma-incubator/compass/components/director/pkg/persistence"
)

type DestinationService struct {
	transact                  persistence.Transactioner
	destinationStorageService DestinationStorageService
}

type DestinationStorageService interface {
}

func NewDestinationService(transact persistence.Transactioner, storageSvc DestinationStorageService) *DestinationService {
	return &DestinationService{
		transact:                  transact,
		destinationStorageService: storageSvc,
	}
}

func (d DestinationService) SyncDestinations() {
	fmt.Println("#########")
}
