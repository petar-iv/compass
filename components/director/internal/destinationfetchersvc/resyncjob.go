package destinationfetchersvc

import (
	"context"
	"github.com/kyma-incubator/compass/components/director/internal/model"
	"github.com/kyma-incubator/compass/components/director/pkg/cronjob"
	"github.com/kyma-incubator/compass/components/director/pkg/log"
	"golang.org/x/sync/semaphore"
	"time"
)

type DestinationResyncer interface {
	FetchDestinationsOnDemand(ctx context.Context, subaccountID string) error
}

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
		Fn: func(ctx context.Context) {
			subscribedTenants, err := tenantFetcher.GetBySubscribedRuntimes(ctx)
			if err != nil {
				log.C(ctx).WithError(err).Errorf("Could not fetch subscribed tenants for destination resync")
				return
			}
			sem := semaphore.NewWeighted(cfg.ParallelTenants)
			for _, tenant := range subscribedTenants {
				resyncDestinations(ctx, destinationResyncer, tenant.ExternalTenant, sem)
			}
		},
		SchedulePeriod: cfg.JobSchedulePeriod,
	}
	return cronjob.RunCronJob(ctx, cfg.ElectionCfg, resyncJob)
}

func resyncDestinations(
	ctx context.Context, destinationResyncer DestinationResyncer, tenantID string, sem *semaphore.Weighted) {

	sem.Acquire(ctx, 1)
	defer sem.Release(1)

	err := destinationResyncer.FetchDestinationsOnDemand(ctx, tenantID)
	if err != nil {
		log.C(ctx).WithError(err).Errorf("Could not resync destinations for tenant %s", tenantID)
	} else {
		log.C(ctx).WithError(err).Debugf("Resynced destinations for tenant %s", tenantID)
	}
}
