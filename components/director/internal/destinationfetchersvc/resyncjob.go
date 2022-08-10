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

// DestinationResyncer missing godoc
//go:generate mockery --name=DestinationResyncer --output=automock --outpkg=automock --case=underscore --disable-version-string
type DestinationResyncer interface {
	FetchDestinationsOnDemand(ctx context.Context, subaccountID string) error
}

// SubscribedTenantFetcher missing godoc
//go:generate mockery --name=SubscribedTenantFetcher --output=automock --outpkg=automock --case=underscore --disable-version-string
type SubscribedTenantFetcher interface {
	GetBySubscribedRuntimes(ctx context.Context) ([]*model.BusinessTenantMapping, error)
}

type ResyncJobConfig struct {
	ElectionCfg       cronjob.ElectionConfig
	JobSchedulePeriod time.Duration
	ParallelTenants   int64
}

func StartDestinationFetcherResyncJob(ctx context.Context, cfg ResyncJobConfig,
	tenantFetcher SubscribedTenantFetcher, destinationResyncer DestinationResyncer) error {

	resyncJob := cronjob.CronJob{
		Name: "DestinationFetcherResync",
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
					resyncTenantDestinations(jobCtx, destinationResyncer, tenantID)
				}(tenant.ExternalTenant)
			}
			wg.Wait()
		},
		SchedulePeriod: cfg.JobSchedulePeriod,
	}
	return cronjob.RunCronJob(ctx, cfg.ElectionCfg, resyncJob)
}

func resyncTenantDestinations(ctx context.Context, destinationResyncer DestinationResyncer, tenantID string) {
	err := destinationResyncer.FetchDestinationsOnDemand(ctx, tenantID)
	if err != nil {
		log.C(ctx).WithError(err).Errorf("Could not resync destinations for tenant %s", tenantID)
	} else {
		log.C(ctx).WithError(err).Debugf("Resynced destinations for tenant %s", tenantID)
	}
}
