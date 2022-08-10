package destinationfetchersvc_test

import (
	"errors"
	"github.com/kyma-incubator/compass/components/director/internal/destinationfetchersvc"
	"github.com/kyma-incubator/compass/components/director/internal/destinationfetchersvc/automock"
	"github.com/kyma-incubator/compass/components/director/internal/model"
	"github.com/kyma-incubator/compass/components/director/pkg/cronjob"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/net/context"
	"testing"
	"time"
)

func TestDestinationResyncJob(t *testing.T) {
	const testTimeout = time.Second * 5

	var (
		tenantsToResync = []*model.BusinessTenantMapping{
			{
				ExternalTenant: "t1",
			},
			{
				ExternalTenant: "t2",
			},
			{
				ExternalTenant: "t3",
			},
		}
		cfg = destinationfetchersvc.ResyncJobConfig{
			ParallelTenants:   2,
			JobSchedulePeriod: time.Minute,
			ElectionCfg: cronjob.ElectionConfig{
				ElectionEnabled: false,
			},
		}
		expectedError = errors.New("expected")

		cancelCtxAfterAllDoneReceived = func(done <-chan struct{}, doneCount int, cancel context.CancelFunc) {
			go func() {
				defer cancel()
				for i := 0; i < doneCount; i++ {
					select {
					case <-done:
					case <-time.After(testTimeout):
						t.Errorf("Test timed out - not all tenants re-synced")
						return
					}
				}
			}()
		}
	)

	t.Run("Should re-sync all tenants", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())

		done := make(chan struct{}, len(tenantsToResync))
		addDone := func(args mock.Arguments) {
			done <- struct{}{}
		}
		cancelCtxAfterAllDoneReceived(done, len(tenantsToResync), cancel)

		subscribedTenantFetcher := &automock.SubscribedTenantFetcher{}
		subscribedTenantFetcher.Mock.On("GetBySubscribedRuntimes", mock.Anything).
			Return(tenantsToResync, nil)

		destinationResyncer := &automock.DestinationResyncer{}
		destinationResyncer.Mock.On("FetchDestinationsOnDemand",
			mock.Anything, tenantsToResync[0].ExternalTenant).Return(nil).Run(addDone)
		destinationResyncer.Mock.On("FetchDestinationsOnDemand",
			mock.Anything, tenantsToResync[1].ExternalTenant).Return(nil).Run(addDone)
		destinationResyncer.Mock.On("FetchDestinationsOnDemand",
			mock.Anything, tenantsToResync[2].ExternalTenant).Return(nil).Run(addDone)

		err := destinationfetchersvc.StartDestinationFetcherResyncJob(ctx, cfg, subscribedTenantFetcher, destinationResyncer)
		assert.Nil(t, err)
		assert.Equal(t, 1, len(subscribedTenantFetcher.Calls))
		assert.Equal(t, len(tenantsToResync), len(destinationResyncer.Calls))
	})

	t.Run("Should not fail on one tenant re-sync failure", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())

		done := make(chan struct{}, len(tenantsToResync))
		addDone := func(args mock.Arguments) {
			done <- struct{}{}
		}
		cancelCtxAfterAllDoneReceived(done, len(tenantsToResync), cancel)

		subscribedTenantFetcher := &automock.SubscribedTenantFetcher{}
		subscribedTenantFetcher.Mock.On("GetBySubscribedRuntimes", mock.Anything).
			Return(tenantsToResync, nil)

		destinationResyncer := &automock.DestinationResyncer{}
		destinationResyncer.Mock.On("FetchDestinationsOnDemand",
			mock.Anything, tenantsToResync[0].ExternalTenant).Return(nil).Run(addDone)
		destinationResyncer.Mock.On("FetchDestinationsOnDemand",
			mock.Anything, tenantsToResync[1].ExternalTenant).Return(expectedError).Run(addDone)
		destinationResyncer.Mock.On("FetchDestinationsOnDemand",
			mock.Anything, tenantsToResync[2].ExternalTenant).Return(nil).Run(addDone)

		err := destinationfetchersvc.StartDestinationFetcherResyncJob(ctx, cfg, subscribedTenantFetcher, destinationResyncer)
		assert.Nil(t, err)
		assert.Equal(t, 1, len(subscribedTenantFetcher.Calls))
		assert.Equal(t, len(tenantsToResync), len(destinationResyncer.Calls))
	})

	t.Run("Should not re-sync if subscribed tenants could not be fetched", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())

		subscribedTenantFetcher := &automock.SubscribedTenantFetcher{}
		subscribedTenantFetcher.Mock.On("GetBySubscribedRuntimes", mock.Anything).
			Return(nil, expectedError).Run(func(args mock.Arguments) { cancel() })

		destinationResyncer := &automock.DestinationResyncer{}

		err := destinationfetchersvc.StartDestinationFetcherResyncJob(ctx, cfg, subscribedTenantFetcher, destinationResyncer)
		assert.Nil(t, err)
		assert.Equal(t, 1, len(subscribedTenantFetcher.Calls))
		assert.Equal(t, 0, len(destinationResyncer.Calls))
	})

	t.Run("Should not re-sync if there are no subscribed tenants", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())

		subscribedTenantFetcher := &automock.SubscribedTenantFetcher{}
		subscribedTenantFetcher.Mock.On("GetBySubscribedRuntimes", mock.Anything).
			Return(nil, nil).Run(func(args mock.Arguments) { cancel() })

		destinationResyncer := &automock.DestinationResyncer{}

		err := destinationfetchersvc.StartDestinationFetcherResyncJob(ctx, cfg, subscribedTenantFetcher, destinationResyncer)
		assert.Nil(t, err)
		assert.Equal(t, 1, len(subscribedTenantFetcher.Calls))
		assert.Equal(t, 0, len(destinationResyncer.Calls))
	})

}
