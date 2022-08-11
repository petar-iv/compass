package destinationfetchersvc_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/kyma-incubator/compass/components/director/internal/destinationfetchersvc"
	"github.com/kyma-incubator/compass/components/director/internal/destinationfetchersvc/automock"
	"github.com/kyma-incubator/compass/components/director/internal/model"
	"github.com/pkg/errors"

	"github.com/kyma-incubator/compass/components/director/pkg/apperrors"
	"github.com/kyma-incubator/compass/components/director/pkg/config"
	persistenceAutomock "github.com/kyma-incubator/compass/components/director/pkg/persistence/automock"
	"github.com/kyma-incubator/compass/components/director/pkg/persistence/txtest"
	"github.com/kyma-incubator/compass/components/director/pkg/resource"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const (
	tenantID          = "f09ba084-0e82-49ab-ab2e-b7ecc988312d"
	userContextHeader = "user_context"
	runtimeID         = "d09ba084-0e82-49ab-ab2e-b7ecc988312d"
	tenantLabelKey    = "subaccount"
	regionLabelKey    = "region"
	region            = "region1"

	subdomainLabelValue = "127" // will be replaced in 127.0.0.1 when fetching token for destination service
)

var (
	labelTenantID = "f09ba084-0e82-49ab-ab2e-b7ecc988312f"
	testErr       = errors.New("test error")
)

func TestService_SyncTenantDestinations(t *testing.T) {
	//GIVEN
	destinationServer := newDestinationServer(t)
	destinationServer.server.Start()

	txGen := txtest.NewTransactionContextGenerator(testErr)

	cert, key := generateTestCertAndKey(t, "test")
	instanceConfig := config.InstanceConfig{
		ClientID:     tenantID,
		ClientSecret: "secret",
		URL:          destinationServer.server.URL,
		TokenURL:     destinationServer.server.URL + "/oauth/token",
		Cert:         string(cert),
		Key:          string(key),
	}

	destAPIConfig := destinationfetchersvc.APIConfig{
		GoroutineLimit:                2,
		RetryInterval:                 0,
		RetryAttempts:                 2,
		EndpointGetTenantDestinations: "/subaccountDestinations",
		EndpointFindDestination:       "/destinations",
		Timeout:                       time.Second * 10,
		PageSize:                      1,
		PagingPageParam:               "$page",
		PagingSizeParam:               "$pageSize",
		PagingCountParam:              "$pageCount",
		PagingCountHeader:             "Page-Count",
	}

	destConfig := config.DestinationsConfig{
		RegionToInstanceConfig: map[string]config.InstanceConfig{
			region: instanceConfig,
		},
		OauthTokenPath: "/oauth-path",
	}

	testCases := []struct {
		Name                string
		LabelRepo           func() *automock.LabelRepo
		DestRepo            func() *automock.DestinationRepo
		Transactioner       func() (*persistenceAutomock.PersistenceTx, *persistenceAutomock.Transactioner)
		BundleRepo          func() *automock.BundleRepo
		UUIDService         func() *automock.UUIDService
		ExpectedErrorOutput string
	}{
		{
			Name: "Sync tenant destinations",
			Transactioner: func() (*persistenceAutomock.PersistenceTx, *persistenceAutomock.Transactioner) {
				return txGen.ThatSucceedsMultipleTimes(4)
			},
			LabelRepo:   successfulLabelRegionAndSubdomainRequest,
			BundleRepo:  successfulBundleRepo("bundleID"),
			DestRepo:    successfulDestinationRepo("bundleID"),
			UUIDService: successfulUUIDService,
		},
		{
			Name: "When getting bundles fails should continue to process destinations",
			Transactioner: func() (*persistenceAutomock.PersistenceTx, *persistenceAutomock.Transactioner) {
				return txGen.ThatSucceedsMultipleTimes(4)
			},
			LabelRepo:   successfulLabelRegionAndSubdomainRequest,
			BundleRepo:  failingBundleRepo,
			DestRepo:    unusedDestinationsRepo,
			UUIDService: unusedUUIDService,
		},
		{
			Name: "When getting bundles fails should continue to process destinations",
			Transactioner: func() (*persistenceAutomock.PersistenceTx, *persistenceAutomock.Transactioner) {
				return txGen.ThatSucceedsMultipleTimes(4)
			},
			LabelRepo:   successfulLabelRegionAndSubdomainRequest,
			BundleRepo:  bundleRepoWithNoBundles,
			DestRepo:    unusedDestinationsRepo,
			UUIDService: unusedUUIDService,
		},
		{
			Name: "When destination upsert fails should continue to process destinations",
			Transactioner: func() (*persistenceAutomock.PersistenceTx, *persistenceAutomock.Transactioner) {
				return txGen.ThatSucceedsMultipleTimes(4)
			},
			LabelRepo:   successfulLabelRegionAndSubdomainRequest,
			BundleRepo:  successfulBundleRepo("bundleID"),
			DestRepo:    failingDestinationRepo,
			UUIDService: successfulUUIDService,
		},
		{
			Name:                "Failed to begin transaction to database",
			Transactioner:       txGen.ThatFailsOnBegin,
			LabelRepo:           unusedLabelRepo,
			BundleRepo:          unusedBundleRepo,
			DestRepo:            unusedDestinationsRepo,
			UUIDService:         unusedUUIDService,
			ExpectedErrorOutput: testErr.Error(),
		},
		{
			Name:          "Failed to find subdomain label",
			Transactioner: txGen.ThatSucceeds,
			LabelRepo: func() *automock.LabelRepo {
				repo := &automock.LabelRepo{}
				repo.On("GetSubdomainLabelForSubscribedRuntime", mock.Anything, tenantID).
					Return(nil, apperrors.NewNotFoundError(resource.Label, "id"))
				return repo
			},
			BundleRepo:          unusedBundleRepo,
			DestRepo:            unusedDestinationsRepo,
			UUIDService:         unusedUUIDService,
			ExpectedErrorOutput: fmt.Sprintf("tenant %s not found", tenantID),
		},
		{
			Name:          "Error while getting subdomain label",
			Transactioner: txGen.ThatSucceeds,
			LabelRepo: func() *automock.LabelRepo {
				repo := &automock.LabelRepo{}
				repo.On("GetSubdomainLabelForSubscribedRuntime", mock.Anything, tenantID).
					Return(nil, testErr)
				return repo
			},
			BundleRepo:          unusedBundleRepo,
			DestRepo:            unusedDestinationsRepo,
			UUIDService:         unusedUUIDService,
			ExpectedErrorOutput: testErr.Error(),
		},
		{
			Name:                "Failed to commit transaction",
			Transactioner:       txGen.ThatFailsOnCommit,
			LabelRepo:           successfulLabelSubdomainRequest,
			BundleRepo:          unusedBundleRepo,
			DestRepo:            unusedDestinationsRepo,
			UUIDService:         unusedUUIDService,
			ExpectedErrorOutput: testErr.Error(),
		},
		{
			Name: "Failed getting region",
			Transactioner: func() (*persistenceAutomock.PersistenceTx, *persistenceAutomock.Transactioner) {
				return txGen.ThatSucceedsMultipleTimes(2)
			},
			LabelRepo:           failedLabelRegionAndSuccessfulSubdomainRequest,
			BundleRepo:          unusedBundleRepo,
			DestRepo:            unusedDestinationsRepo,
			UUIDService:         unusedUUIDService,
			ExpectedErrorOutput: testErr.Error(),
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			_, tx := testCase.Transactioner()
			destRepo := testCase.DestRepo()
			labelRepo := testCase.LabelRepo()
			bundleRepo := testCase.BundleRepo()
			uuidService := testCase.UUIDService()
			defer mock.AssertExpectationsForObjects(t, tx, destRepo, labelRepo, uuidService, bundleRepo)

			destSvc := destinationfetchersvc.DestinationService{
				Transactioner:      tx,
				UUIDSvc:            uuidService,
				Repo:               destRepo,
				BundleRepo:         bundleRepo,
				LabelRepo:          labelRepo,
				DestinationsConfig: destConfig,
				APIConfig:          destAPIConfig,
			}

			ctx := context.Background()
			// WHEN
			err := destSvc.SyncTenantDestinations(ctx, tenantID)

			// THEN
			if len(testCase.ExpectedErrorOutput) > 0 {
				assert.ErrorContains(t, err, testCase.ExpectedErrorOutput)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func unusedLabelRepo() *automock.LabelRepo              { return &automock.LabelRepo{} }
func unusedDestinationsRepo() *automock.DestinationRepo { return &automock.DestinationRepo{} }
func unusedBundleRepo() *automock.BundleRepo            { return &automock.BundleRepo{} }
func unusedUUIDService() *automock.UUIDService          { return &automock.UUIDService{} }

func successfulUUIDService() *automock.UUIDService {
	uuidService := &automock.UUIDService{}
	uuidService.On("Generate").Return("9b26a428-d526-469c-a5ef-2856f3ce0430")
	return uuidService
}
func successfulLabelSubdomainRequest() *automock.LabelRepo {
	repo := &automock.LabelRepo{}
	label := model.NewLabelForRuntime(runtimeID, labelTenantID, tenantLabelKey, subdomainLabelValue)
	label.Tenant = &labelTenantID
	repo.On("GetSubdomainLabelForSubscribedRuntime", mock.Anything, tenantID).Return(label, nil)
	return repo
}

func failedLabelRegionAndSuccessfulSubdomainRequest() *automock.LabelRepo {
	repo := &automock.LabelRepo{}
	label := model.NewLabelForRuntime(runtimeID, labelTenantID, tenantLabelKey, subdomainLabelValue)
	label.Tenant = &labelTenantID
	repo.On("GetSubdomainLabelForSubscribedRuntime", mock.Anything, tenantID).Return(label, nil)
	repo.On("GetByKey", mock.Anything, labelTenantID, model.TenantLabelableObject, labelTenantID, regionLabelKey).
		Return(nil, testErr)
	return repo
}

func successfulLabelRegionAndSubdomainRequest() *automock.LabelRepo {
	repo := &automock.LabelRepo{}
	label := model.NewLabelForRuntime(runtimeID, labelTenantID, tenantLabelKey, subdomainLabelValue)
	label.Tenant = &labelTenantID
	repo.On("GetSubdomainLabelForSubscribedRuntime", mock.Anything, tenantID).Return(label, nil)
	label = model.NewLabelForRuntime(runtimeID, labelTenantID, regionLabelKey, region)
	label.Tenant = &labelTenantID
	repo.On("GetByKey", mock.Anything, labelTenantID, model.TenantLabelableObject, labelTenantID, regionLabelKey).
		Return(label, nil)
	return repo
}

func successfulBundleRepo(bundleID string) func() *automock.BundleRepo {
	return func() *automock.BundleRepo {
		bundleRepo := unusedBundleRepo()
		bundleRepo.On("GetBySystemAndCorrelationId",
			mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(
			[]*model.Bundle{{
				BaseEntity: &model.BaseEntity{
					ID: bundleID,
				},
			}}, nil)
		return bundleRepo
	}
}

func failingBundleRepo() *automock.BundleRepo {
	bundleRepo := unusedBundleRepo()
	bundleRepo.On("GetBySystemAndCorrelationId",
		mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil, testErr)
	return bundleRepo
}

func bundleRepoWithNoBundles() *automock.BundleRepo {
	bundleRepo := unusedBundleRepo()
	bundleRepo.On("GetBySystemAndCorrelationId",
		mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return([]*model.Bundle{}, nil)
	return bundleRepo
}

func successfulDestinationRepo(bundleID string) func() *automock.DestinationRepo {
	return func() *automock.DestinationRepo {
		destinationRepo := unusedDestinationsRepo()
		destinationRepo.On("Upsert",
			mock.Anything, mock.Anything, mock.Anything, mock.Anything, bundleID, mock.Anything).Return(nil)
		return destinationRepo
	}
}

func failingDestinationRepo() *automock.DestinationRepo {
	destinationRepo := unusedDestinationsRepo()
	destinationRepo.On("Upsert",
		mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(testErr)
	return destinationRepo
}
