package destinationfetchersvc

import (
	"context"
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
	GetSubscribedTenantIDs(ctx context.Context) ([]string, error)
}

type SyncJobConfig struct {
	ElectionCfg       cronjob.ElectionConfig
	JobSchedulePeriod time.Duration
	ParallelTenants   int64
}

func StartDestinationFetcherSyncJob(ctx context.Context, cfg SyncJobConfig, destinationSyncer DestinationSyncer) error {
	resyncJob := cronjob.CronJob{
		Name: "DestinationFetcherSync",
		Fn: func(jobCtx context.Context) {
			subscribedTenants, err := destinationSyncer.GetSubscribedTenantIDs(jobCtx)
			if err != nil {
				log.C(jobCtx).WithError(err).Errorf("Could not fetch subscribed tenants for destination resync")
				return
			}
			if len(subscribedTenants) == 0 {
				log.C(jobCtx).Info("No subscribed tenants found. Skipping sync job")
				return
			}
			sem := semaphore.NewWeighted(cfg.ParallelTenants)
			wg := &sync.WaitGroup{}
			for _, tenantID := range subscribedTenants {
				wg.Add(1)
				go func(tenantID string) {
					defer wg.Done()
					if err := sem.Acquire(jobCtx, 1); err != nil {
						log.C(jobCtx).WithError(err).Errorf("Could not acquire semaphor")
						return
					}
					defer sem.Release(1)
					syncTenantDestinations(jobCtx, destinationSyncer, tenantID)
				}(tenantID)
			}
			wg.Wait()
		},
		SchedulePeriod: cfg.JobSchedulePeriod,
	}
	return cronjob.RunCronJob(ctx, cfg.ElectionCfg, resyncJob)
}

func syncTenantDestinations(ctx context.Context, destinationSyncer DestinationSyncer, tenantID string) {
	err := destinationSyncer.SyncTenantDestinations(ctx, tenantID)
	if err != nil {
		log.C(ctx).WithError(err).Errorf("Could not resync destinations for tenant %s", tenantID)
	} else {
		log.C(ctx).WithError(err).Debugf("Resynced destinations for tenant %s", tenantID)
	}
}