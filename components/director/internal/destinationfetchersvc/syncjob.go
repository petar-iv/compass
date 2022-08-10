package destinationfetchersvc

import (
	"context"
	"github.com/kyma-incubator/compass/components/director/internal/model"
	"github.com/kyma-incubator/compass/components/director/pkg/cronjob"
	"github.com/kyma-incubator/compass/components/director/pkg/log"
	"golang.org/x/sync/semaphore"
	"sync"
	"time"
)

// DestinationSyncer missing godoc
//go:generate mockery --name=DestinationSyncer --output=automock --outpkg=automock --case=underscore --disable-version-string
type DestinationSyncer interface {
	SyncTenantDestinations(ctx context.Context, tenantID string) error
}

// SubscribedTenantFetcher missing godoc
//go:generate mockery --name=SubscribedTenantFetcher --output=automock --outpkg=automock --case=underscore --disable-version-string
type SubscribedTenantFetcher interface {
	GetBySubscribedRuntimes(ctx context.Context) ([]*model.BusinessTenantMapping, error)
}

type SyncJobConfig struct {
	ElectionCfg       cronjob.ElectionConfig
	JobSchedulePeriod time.Duration
	ParallelTenants   int64
}

func StartDestinationFetcherSyncJob(ctx context.Context, cfg SyncJobConfig,
	tenantFetcher SubscribedTenantFetcher, destinationSyncer DestinationSyncer) error {

	resyncJob := cronjob.CronJob{
		Name: "DestinationFetcherSync",
		Fn: func(jobCtx context.Context) {
			subscribedTenants, err := tenantFetcher.GetBySubscribedRuntimes(jobCtx)
			if err != nil {
				log.C(jobCtx).WithError(err).Errorf("Could not fetch subscribed tenants for destination resync")
				return
			}
			sem := semaphore.NewWeighted(cfg.ParallelTenants)
			wg := &sync.WaitGroup{}
			for _, tenant := range subscribedTenants {
				wg.Add(1)
				go func(tenantID string) {
					defer wg.Done()
					if err := sem.Acquire(jobCtx, 1); err != nil {
						log.C(jobCtx).WithError(err).Errorf("Could not acquire semaphor")
						return
					}
					defer sem.Release(1)
					resyncTenantDestinations(jobCtx, destinationSyncer, tenantID)
				}(tenant.ExternalTenant)
			}
			wg.Wait()
		},
		SchedulePeriod: cfg.JobSchedulePeriod,
	}
	return cronjob.RunCronJob(ctx, cfg.ElectionCfg, resyncJob)
}

func resyncTenantDestinations(ctx context.Context, destinationSyncer DestinationSyncer, tenantID string) {
	err := destinationSyncer.SyncTenantDestinations(ctx, tenantID)
	if err != nil {
		log.C(ctx).WithError(err).Errorf("Could not resync destinations for tenant %s", tenantID)
	} else {
		log.C(ctx).WithError(err).Debugf("Resynced destinations for tenant %s", tenantID)
	}
}
