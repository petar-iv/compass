package destinationfetchersvc_test

import (
	"context"
	"fmt"
	"testing"

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
	tenantID           = "f09ba084-0e82-49ab-ab2e-b7ecc988312d"
	userContextHeader  = "user_context"
	runtimeID          = "d09ba084-0e82-49ab-ab2e-b7ecc988312d"
	subaccountLabelKey = "subaccount"
	regionLabelKey     = "region"
	region             = "region1"
)

var (
	labelTenantID = "f09ba084-0e82-49ab-ab2e-b7ecc988312f"
	testErr       = errors.New("test error")
)

func TestService_SyncSubaccountDestinations(t *testing.T) {
	//GIVEN

	txGen := txtest.NewTransactionContextGenerator(testErr)

	cert, key := generateTestCertAndKey(t, "test")
	instanceConfig := config.InstanceConfig{
		ClientID:     tenantID,
		ClientSecret: "secret",
		URL:          "https://destination-configuration.com",
		TokenURL:     "https://test.auth.com",
		Cert:         string(cert),
		Key:          string(key),
	}

	destAPIConfig := destinationfetchersvc.APIConfig{
		GoroutineLimit: 10,
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
		TenantRepo          func() *automock.TenantRepo
		BundleRepo          func() *automock.BundleRepo
		UUIDService         func() *automock.UUIDService
		ExpectedErrorOutput string
	}{
		{
			Name:                "Failed to begin transaction to database",
			Transactioner:       txGen.ThatFailsOnBegin,
			LabelRepo:           unusedLabelRepo,
			TenantRepo:          unusedTenantRepo,
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
			TenantRepo:          unusedTenantRepo,
			BundleRepo:          unusedBundleRepo,
			DestRepo:            unusedDestinationsRepo,
			UUIDService:         unusedUUIDService,
			ExpectedErrorOutput: fmt.Sprintf("subaccount %s not found", tenantID),
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
			TenantRepo:          unusedTenantRepo,
			BundleRepo:          unusedBundleRepo,
			DestRepo:            unusedDestinationsRepo,
			UUIDService:         unusedUUIDService,
			ExpectedErrorOutput: testErr.Error(),
		},
		{
			Name:                "Failed to commit transaction",
			Transactioner:       txGen.ThatFailsOnCommit,
			LabelRepo:           successfullLabelSubdomainRequest,
			TenantRepo:          unusedTenantRepo,
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
			LabelRepo:           failedLabelRegionAndSuccesfullSubdomainRequest,
			TenantRepo:          unusedTenantRepo,
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
			tenantRepo := testCase.TenantRepo()
			bundleRepo := testCase.BundleRepo()
			uuidService := testCase.UUIDService()
			defer mock.AssertExpectationsForObjects(t, tx, destRepo, labelRepo, uuidService, tenantRepo, bundleRepo)

			destSvc := destinationfetchersvc.DestinationService{
				Transactioner:      tx,
				UUIDSvc:            uuidService,
				Repo:               destRepo,
				BundleRepo:         bundleRepo,
				LabelRepo:          labelRepo,
				TenantRepo:         tenantRepo,
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
func unusedTenantRepo() *automock.TenantRepo            { return &automock.TenantRepo{} }
func unusedDestinationsRepo() *automock.DestinationRepo { return &automock.DestinationRepo{} }
func unusedBundleRepo() *automock.BundleRepo            { return &automock.BundleRepo{} }
func unusedUUIDService() *automock.UUIDService          { return &automock.UUIDService{} }
func usedUUIDService() *automock.UUIDService {
	uuidService := &automock.UUIDService{}
	uuidService.On("Generate").Return("9b26a428-d526-469c-a5ef-2856f3ce0430")
	return uuidService
}
func successfullLabelSubdomainRequest() *automock.LabelRepo {
	labelValue := "ta-subdomain"
	repo := &automock.LabelRepo{}
	label := model.NewLabelForRuntime(runtimeID, labelTenantID, subaccountLabelKey, labelValue)
	label.Tenant = &labelTenantID
	repo.On("GetSubdomainLabelForSubscribedRuntime", mock.Anything, tenantID).
		Return(label, nil)
	return repo
}

func failedLabelRegionAndSuccesfullSubdomainRequest() *automock.LabelRepo {
	labelValue := "ta-subdomain"
	repo := &automock.LabelRepo{}
	label := model.NewLabelForRuntime(runtimeID, labelTenantID, subaccountLabelKey, labelValue)
	label.Tenant = &labelTenantID
	repo.On("GetSubdomainLabelForSubscribedRuntime", mock.Anything, tenantID).
		Return(label, nil)
	repo.On("GetByKey", mock.Anything, labelTenantID, model.TenantLabelableObject, labelTenantID, regionLabelKey).
		Return(nil, testErr)

	return repo
}

func successfullLabelRegionAndSubdomainRequest() *automock.LabelRepo {
	labelValue := "ta-subdomain"
	repo := &automock.LabelRepo{}
	label := model.NewLabelForRuntime(runtimeID, labelTenantID, subaccountLabelKey, labelValue)
	label.Tenant = &labelTenantID
	repo.On("GetSubdomainLabelForSubscribedRuntime", mock.Anything, tenantID).
		Return(label, nil)
	label = model.NewLabelForRuntime(runtimeID, labelTenantID, regionLabelKey, region)
	label.Tenant = &labelTenantID
	repo.On("GetByKey", mock.Anything, labelTenantID, model.TenantLabelableObject, labelTenantID, regionLabelKey).
		Return(label, nil)

	return repo
}
