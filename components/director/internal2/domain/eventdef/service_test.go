package eventdef_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/pkg/errors"

	"github.com/kyma-incubator/compass/components/director/pkg/resource"

	"github.com/kyma-incubator/compass/components/director/pkg/apperrors"

	"github.com/kyma-incubator/compass/components/director/pkg/pagination"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/kyma-incubator/compass/components/director/internal2/domain/eventdef"
	"github.com/kyma-incubator/compass/components/director/internal2/domain/eventdef/automock"
	"github.com/kyma-incubator/compass/components/director/internal2/domain/tenant"
	"github.com/kyma-incubator/compass/components/director/internal2/model"
	"github.com/stretchr/testify/assert"
)

func TestService_Get(t *testing.T) {
	// given
	testErr := errors.New("Test error")
	id := "foo"
	eventAPIDefinition := fixMinModelEventAPIDefinition(id, "placeholder")

	ctx := context.TODO()
	ctx = tenant.SaveToContext(ctx, tenantID, externalTenantID)

	testCases := []struct {
		Name               string
		RepositoryFn       func() *automock.EventAPIRepository
		Input              model.EventDefinitionInput
		InputID            string
		ExpectedDocument   *model.EventDefinition
		ExpectedErrMessage string
	}{
		{
			Name: "Success",
			RepositoryFn: func() *automock.EventAPIRepository {
				repo := &automock.EventAPIRepository{}
				repo.On("GetByID", ctx, tenantID, id).Return(eventAPIDefinition, nil).Once()
				return repo
			},
			InputID:            id,
			ExpectedDocument:   eventAPIDefinition,
			ExpectedErrMessage: "",
		},
		{
			Name: "Returns error when Event Definition retrieval failed",
			RepositoryFn: func() *automock.EventAPIRepository {
				repo := &automock.EventAPIRepository{}
				repo.On("GetByID", ctx, tenantID, id).Return(nil, testErr).Once()
				return repo
			},
			InputID:            id,
			ExpectedDocument:   eventAPIDefinition,
			ExpectedErrMessage: testErr.Error(),
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			repo := testCase.RepositoryFn()

			svc := eventdef.NewService(repo, nil, nil, nil)

			// when
			eventAPIDefinition, err := svc.Get(ctx, testCase.InputID)

			// then
			if testCase.ExpectedErrMessage == "" {
				require.NoError(t, err)
				assert.Equal(t, testCase.ExpectedDocument, eventAPIDefinition)
			} else {
				assert.Contains(t, err.Error(), testCase.ExpectedErrMessage)
			}

			repo.AssertExpectations(t)
		})
	}
	t.Run("Error when tenant not in context", func(t *testing.T) {
		svc := eventdef.NewService(nil, nil, nil, nil)
		// WHEN
		_, err := svc.Get(context.TODO(), "")
		// THEN
		require.Error(t, err)
		assert.Contains(t, err.Error(), "cannot read tenant from context")
	})
}

func TestService_GetForBundle(t *testing.T) {
	// given
	testErr := errors.New("Test error")
	id := "foo"
	bundleID := "test"
	eventAPIDefinition := fixMinModelEventAPIDefinition(id, "placeholder")

	ctx := context.TODO()
	ctx = tenant.SaveToContext(ctx, tenantID, externalTenantID)

	testCases := []struct {
		Name               string
		RepositoryFn       func() *automock.EventAPIRepository
		Input              model.EventDefinitionInput
		InputID            string
		BundleID           string
		ExpectedEventDef   *model.EventDefinition
		ExpectedErrMessage string
	}{
		{
			Name: "Success",
			RepositoryFn: func() *automock.EventAPIRepository {
				repo := &automock.EventAPIRepository{}
				repo.On("GetForBundle", ctx, tenantID, id, bundleID).Return(eventAPIDefinition, nil).Once()
				return repo
			},
			InputID:            id,
			BundleID:           bundleID,
			ExpectedEventDef:   eventAPIDefinition,
			ExpectedErrMessage: "",
		},
		{
			Name: "Returns error when Event Definition retrieval failed",
			RepositoryFn: func() *automock.EventAPIRepository {
				repo := &automock.EventAPIRepository{}
				repo.On("GetForBundle", ctx, tenantID, id, bundleID).Return(nil, testErr).Once()
				return repo
			},
			InputID:            id,
			BundleID:           bundleID,
			ExpectedEventDef:   eventAPIDefinition,
			ExpectedErrMessage: testErr.Error(),
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			repo := testCase.RepositoryFn()

			svc := eventdef.NewService(repo, nil, nil, nil)

			// when
			eventAPIDefinition, err := svc.GetForBundle(ctx, testCase.InputID, testCase.BundleID)

			// then
			if testCase.ExpectedErrMessage == "" {
				require.NoError(t, err)
				assert.Equal(t, testCase.ExpectedEventDef, eventAPIDefinition)
			} else {
				assert.Contains(t, err.Error(), testCase.ExpectedErrMessage)
			}

			repo.AssertExpectations(t)
		})
	}
	t.Run("Error when tenant not in context", func(t *testing.T) {
		svc := eventdef.NewService(nil, nil, nil, nil)
		// WHEN
		_, err := svc.GetForBundle(context.TODO(), "", "")
		// THEN
		require.Error(t, err)
		assert.Contains(t, err.Error(), "cannot read tenant from context")
	})
}

func TestService_ListForBundle(t *testing.T) {
	// given
	testErr := errors.New("Test error")

	id := "foo"
	bundleID := "bar"

	eventAPIDefinitions := []*model.EventDefinition{
		fixMinModelEventAPIDefinition(id, "placeholder"),
		fixMinModelEventAPIDefinition(id, "placeholder"),
		fixMinModelEventAPIDefinition(id, "placeholder"),
	}
	eventAPIDefinitionPage := &model.EventDefinitionPage{
		Data:       eventAPIDefinitions,
		TotalCount: len(eventAPIDefinitions),
		PageInfo: &pagination.Page{
			HasNextPage: false,
			EndCursor:   "end",
			StartCursor: "start",
		},
	}

	first := 2
	after := "test"

	ctx := context.TODO()
	ctx = tenant.SaveToContext(ctx, tenantID, externalTenantID)

	testCases := []struct {
		Name               string
		RepositoryFn       func() *automock.EventAPIRepository
		InputPageSize      int
		InputCursor        string
		ExpectedResult     *model.EventDefinitionPage
		ExpectedErrMessage string
	}{
		{
			Name: "Success",
			RepositoryFn: func() *automock.EventAPIRepository {
				repo := &automock.EventAPIRepository{}
				repo.On("ListForBundle", ctx, tenantID, bundleID, first, after).Return(eventAPIDefinitionPage, nil).Once()
				return repo
			},
			InputPageSize:      first,
			InputCursor:        after,
			ExpectedResult:     eventAPIDefinitionPage,
			ExpectedErrMessage: "",
		},
		{
			Name: "Return error when page size is less than 1",
			RepositoryFn: func() *automock.EventAPIRepository {
				repo := &automock.EventAPIRepository{}
				return repo
			},
			InputPageSize:      0,
			InputCursor:        after,
			ExpectedResult:     eventAPIDefinitionPage,
			ExpectedErrMessage: "page size must be between 1 and 200",
		},
		{
			Name: "Return error when page size is bigger than 200",
			RepositoryFn: func() *automock.EventAPIRepository {
				repo := &automock.EventAPIRepository{}
				return repo
			},
			InputPageSize:      201,
			InputCursor:        after,
			ExpectedResult:     eventAPIDefinitionPage,
			ExpectedErrMessage: "page size must be between 1 and 200",
		},
		{
			Name: "Returns error when Event Definition listing failed",
			RepositoryFn: func() *automock.EventAPIRepository {
				repo := &automock.EventAPIRepository{}
				repo.On("ListForBundle", ctx, tenantID, bundleID, first, after).Return(nil, testErr).Once()
				return repo
			},
			InputPageSize:      first,
			InputCursor:        after,
			ExpectedResult:     nil,
			ExpectedErrMessage: testErr.Error(),
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			repo := testCase.RepositoryFn()

			svc := eventdef.NewService(repo, nil, nil, nil)

			// when
			docs, err := svc.ListForBundle(ctx, bundleID, testCase.InputPageSize, testCase.InputCursor)

			// then
			if testCase.ExpectedErrMessage == "" {
				require.NoError(t, err)
				assert.Equal(t, testCase.ExpectedResult, docs)
			} else {
				assert.Contains(t, err.Error(), testCase.ExpectedErrMessage)
			}

			repo.AssertExpectations(t)
		})
	}
	t.Run("Error when tenant not in context", func(t *testing.T) {
		svc := eventdef.NewService(nil, nil, nil, nil)
		// WHEN
		_, err := svc.ListForBundle(context.TODO(), "", 5, "")
		// THEN
		require.Error(t, err)
		assert.Contains(t, err.Error(), "cannot read tenant from context")
	})
}

func TestService_CreateToBundle(t *testing.T) {
	// given
	testErr := errors.New("Test error")

	id := "foo"
	bundleID := "bndlid"
	name := "Foo"

	timestamp := time.Now()
	frID := "fr-id"
	frURL := "foo.bar"
	spec := "test"

	modelFr := fixModelFetchRequest(frID, frURL, timestamp)

	modelInput := model.EventDefinitionInput{
		Name: name,
		Spec: &model.EventSpecInput{
			FetchRequest: &model.FetchRequestInput{
				URL: frURL,
			},
		},
		Version: &model.VersionInput{},
	}

	modelEventAPIDefinition := &model.EventDefinition{
		ID:       id,
		Tenant:   tenantID,
		BundleID: bundleID,
		Name:     name,
		Spec:     &model.EventSpec{},
		Version:  &model.Version{},
	}

	modelEventAPIDefinitionWithSpec := &model.EventDefinition{
		ID:       id,
		BundleID: bundleID,
		Tenant:   tenantID,
		Name:     name,
		Spec:     &model.EventSpec{Data: &spec},
		Version:  &model.Version{},
	}

	ctx := context.TODO()
	ctx = tenant.SaveToContext(ctx, tenantID, externalTenantID)

	testCases := []struct {
		Name                  string
		RepositoryFn          func() *automock.EventAPIRepository
		FetchRequestRepoFn    func() *automock.FetchRequestRepository
		UIDServiceFn          func() *automock.UIDService
		FetchRequestServiceFn func() *automock.FetchRequestService
		Input                 model.EventDefinitionInput
		ExpectedErr           error
	}{
		{
			Name: "Success",
			RepositoryFn: func() *automock.EventAPIRepository {
				repo := &automock.EventAPIRepository{}
				repo.On("Create", ctx, modelEventAPIDefinition).Return(nil).Once()
				repo.On("Update", ctx, modelEventAPIDefinition).Return(nil).Once()
				return repo
			},
			FetchRequestRepoFn: func() *automock.FetchRequestRepository {
				repo := &automock.FetchRequestRepository{}
				repo.On("Create", ctx, fixModelFetchRequest(frID, frURL, timestamp)).Return(nil).Once()
				return repo
			},
			UIDServiceFn: func() *automock.UIDService {
				svc := &automock.UIDService{}
				svc.On("Generate").Return(id).Once()
				svc.On("Generate").Return(frID).Once()
				return svc
			},
			FetchRequestServiceFn: func() *automock.FetchRequestService {
				svc := &automock.FetchRequestService{}
				svc.On("HandleSpec", ctx, fixModelFetchRequest(frID, frURL, timestamp)).Return(nil)
				return svc
			},
			Input:       modelInput,
			ExpectedErr: nil,
		},
		{
			Name: "Success fetched EventAPI Spec",
			RepositoryFn: func() *automock.EventAPIRepository {
				repo := &automock.EventAPIRepository{}
				repo.On("Create", ctx, modelEventAPIDefinition).Return(nil).Once()
				repo.On("Update", ctx, modelEventAPIDefinitionWithSpec).Return(nil).Once()

				return repo
			},
			FetchRequestRepoFn: func() *automock.FetchRequestRepository {
				repo := &automock.FetchRequestRepository{}
				repo.On("Create", ctx, modelFr).Return(nil).Once()

				return repo
			},
			UIDServiceFn: func() *automock.UIDService {
				svc := &automock.UIDService{}
				svc.On("Generate").Return(id).Once()
				svc.On("Generate").Return(frID).Once()
				return svc
			},
			FetchRequestServiceFn: func() *automock.FetchRequestService {
				svc := &automock.FetchRequestService{}
				svc.On("HandleSpec", ctx, fixModelFetchRequest(frID, frURL, timestamp)).Return(&spec)
				return svc
			},
			Input:       modelInput,
			ExpectedErr: nil,
		},
		{
			Name: "Error - Event Definition Creation",
			RepositoryFn: func() *automock.EventAPIRepository {
				repo := &automock.EventAPIRepository{}
				repo.On("Create", ctx, modelEventAPIDefinition).Return(testErr).Once()
				return repo
			},
			FetchRequestRepoFn: func() *automock.FetchRequestRepository {
				repo := &automock.FetchRequestRepository{}
				return repo
			},
			UIDServiceFn: func() *automock.UIDService {
				svc := &automock.UIDService{}
				svc.On("Generate").Return(id).Once()
				return svc
			},
			FetchRequestServiceFn: func() *automock.FetchRequestService {
				svc := &automock.FetchRequestService{}
				return svc
			},
			Input:       modelInput,
			ExpectedErr: testErr,
		},
		{
			Name: "Error - Fetch Request Creation",
			RepositoryFn: func() *automock.EventAPIRepository {
				repo := &automock.EventAPIRepository{}
				repo.On("Create", ctx, modelEventAPIDefinition).Return(nil).Once()
				return repo
			},
			FetchRequestRepoFn: func() *automock.FetchRequestRepository {
				repo := &automock.FetchRequestRepository{}
				repo.On("Create", ctx, fixModelFetchRequest(frID, frURL, timestamp)).Return(testErr).Once()
				return repo
			},
			UIDServiceFn: func() *automock.UIDService {
				svc := &automock.UIDService{}
				svc.On("Generate").Return(id).Once()
				svc.On("Generate").Return(frID).Once()
				return svc
			},
			FetchRequestServiceFn: func() *automock.FetchRequestService {
				svc := &automock.FetchRequestService{}
				return svc
			},
			Input:       modelInput,
			ExpectedErr: testErr,
		},
		{
			Name: "Error - EventAPI Update",
			RepositoryFn: func() *automock.EventAPIRepository {
				repo := &automock.EventAPIRepository{}
				repo.On("Create", ctx, modelEventAPIDefinition).Return(nil).Once()
				repo.On("Update", ctx, modelEventAPIDefinition).Return(testErr).Once()

				return repo
			},
			FetchRequestRepoFn: func() *automock.FetchRequestRepository {
				repo := &automock.FetchRequestRepository{}
				repo.On("Create", ctx, modelFr).Return(nil).Once()
				return repo
			},
			UIDServiceFn: func() *automock.UIDService {
				svc := &automock.UIDService{}
				svc.On("Generate").Return(id).Once()
				svc.On("Generate").Return(frID).Once()
				return svc
			},
			FetchRequestServiceFn: func() *automock.FetchRequestService {
				svc := &automock.FetchRequestService{}
				svc.On("HandleSpec", ctx, modelFr).Return(nil)
				return svc
			},
			Input:       modelInput,
			ExpectedErr: testErr,
		},
		{
			Name: "Success when fetching EventAPI Spec failed",
			RepositoryFn: func() *automock.EventAPIRepository {
				repo := &automock.EventAPIRepository{}
				repo.On("Create", ctx, modelEventAPIDefinition).Return(nil).Once()
				repo.On("Update", ctx, modelEventAPIDefinition).Return(nil).Once()

				return repo
			},
			FetchRequestRepoFn: func() *automock.FetchRequestRepository {
				repo := &automock.FetchRequestRepository{}
				repo.On("Create", ctx, modelFr).Return(nil).Once()
				return repo
			},
			UIDServiceFn: func() *automock.UIDService {
				svc := &automock.UIDService{}
				svc.On("Generate").Return(id).Once()
				svc.On("Generate").Return(frID).Once()
				return svc
			},
			FetchRequestServiceFn: func() *automock.FetchRequestService {
				svc := &automock.FetchRequestService{}
				svc.On("HandleSpec", ctx, fixModelFetchRequest(frID, frURL, timestamp)).Return(nil)
				return svc
			},
			Input:       modelInput,
			ExpectedErr: nil,
		},
	}

	for _, testCase := range testCases {
		t.Run(fmt.Sprintf("%s", testCase.Name), func(t *testing.T) {
			// given
			repo := testCase.RepositoryFn()
			fetchRequestRepo := testCase.FetchRequestRepoFn()
			uidSvc := testCase.UIDServiceFn()
			fetchRequestService := testCase.FetchRequestServiceFn()

			svc := eventdef.NewService(repo, fetchRequestRepo, uidSvc, fetchRequestService)
			svc.SetTimestampGen(func() time.Time { return timestamp })

			// when
			result, err := svc.CreateInBundle(ctx, bundleID, testCase.Input)

			// then
			if testCase.ExpectedErr != nil {
				require.Error(t, err)
				assert.Contains(t, err.Error(), testCase.ExpectedErr.Error())
			} else {
				assert.IsType(t, "string", result)
			}

			repo.AssertExpectations(t)
			fetchRequestRepo.AssertExpectations(t)
			uidSvc.AssertExpectations(t)
		})
	}
	t.Run("Error when tenant not in context", func(t *testing.T) {
		svc := eventdef.NewService(nil, nil, nil, nil)
		// WHEN
		_, err := svc.CreateInBundle(context.TODO(), "", model.EventDefinitionInput{})
		// THEN
		require.Error(t, err)
		assert.Contains(t, err.Error(), "cannot read tenant from context")
	})
}

func TestService_Update(t *testing.T) {
	// given
	testErr := errors.New("Test error")

	id := "foo"
	timestamp := time.Now()
	frID := "fr-id"
	frURL := "foo.bar"

	modelInput := model.EventDefinitionInput{
		Name: "Foo",
		Spec: &model.EventSpecInput{
			FetchRequest: &model.FetchRequestInput{
				URL: frURL,
			},
		},
		Version: &model.VersionInput{},
	}

	inputEventAPIDefinitionModel := mock.MatchedBy(func(api *model.EventDefinition) bool {
		return api.Name == modelInput.Name
	})

	eventAPIDefinitionModel := &model.EventDefinition{
		Name:     "Bar",
		Tenant:   tenantID,
		BundleID: "id",
		Spec:     &model.EventSpec{},
		Version:  &model.Version{},
	}

	ctx := context.TODO()
	ctx = tenant.SaveToContext(ctx, tenantID, externalTenantID)

	modelFr := fixModelFetchRequest(frID, frURL, timestamp)

	testCases := []struct {
		Name                  string
		RepositoryFn          func() *automock.EventAPIRepository
		FetchRequestRepoFn    func() *automock.FetchRequestRepository
		UIDServiceFn          func() *automock.UIDService
		FetchRequestServiceFn func() *automock.FetchRequestService
		Input                 model.EventDefinitionInput
		InputID               string
		ExpectedErr           error
	}{
		{
			Name: "Success",
			RepositoryFn: func() *automock.EventAPIRepository {
				repo := &automock.EventAPIRepository{}
				repo.On("GetByID", ctx, tenantID, id).Return(eventAPIDefinitionModel, nil).Once()
				repo.On("Update", ctx, inputEventAPIDefinitionModel).Return(nil).Once()
				return repo
			},
			FetchRequestRepoFn: func() *automock.FetchRequestRepository {
				repo := &automock.FetchRequestRepository{}
				repo.On("DeleteByReferenceObjectID", ctx, tenantID, model.EventAPIFetchRequestReference, id).Return(nil).Once()
				repo.On("Create", ctx, fixModelFetchRequest(frID, frURL, timestamp)).Return(nil).Once()
				return repo
			},
			UIDServiceFn: func() *automock.UIDService {
				svc := &automock.UIDService{}
				svc.On("Generate").Return(frID).Once()
				return svc
			},
			FetchRequestServiceFn: func() *automock.FetchRequestService {
				svc := &automock.FetchRequestService{}
				svc.On("HandleSpec", ctx, modelFr).Return(nil)
				return svc
			},
			InputID:     "foo",
			Input:       modelInput,
			ExpectedErr: nil,
		},
		{
			Name: "Update Error",
			RepositoryFn: func() *automock.EventAPIRepository {
				repo := &automock.EventAPIRepository{}
				repo.On("GetByID", ctx, tenantID, id).Return(eventAPIDefinitionModel, nil).Once()
				repo.On("Update", ctx, inputEventAPIDefinitionModel).Return(testErr).Once()
				return repo
			},
			FetchRequestRepoFn: func() *automock.FetchRequestRepository {
				repo := &automock.FetchRequestRepository{}
				repo.On("DeleteByReferenceObjectID", ctx, tenantID, model.EventAPIFetchRequestReference, id).Return(nil).Once()
				repo.On("Create", ctx, fixModelFetchRequest(frID, frURL, timestamp)).Return(nil).Once()
				return repo
			},
			UIDServiceFn: func() *automock.UIDService {
				svc := &automock.UIDService{}
				svc.On("Generate").Return(frID).Once()
				return svc
			},
			FetchRequestServiceFn: func() *automock.FetchRequestService {
				svc := &automock.FetchRequestService{}
				svc.On("HandleSpec", ctx, modelFr).Return(nil)
				return svc
			},
			InputID:     "foo",
			Input:       modelInput,
			ExpectedErr: testErr,
		},
		{
			Name: "Delete FetchRequest by reference Error",
			RepositoryFn: func() *automock.EventAPIRepository {
				repo := &automock.EventAPIRepository{}
				repo.On("GetByID", ctx, tenantID, "foo").Return(eventAPIDefinitionModel, nil).Once()
				return repo
			},
			FetchRequestRepoFn: func() *automock.FetchRequestRepository {
				repo := &automock.FetchRequestRepository{}
				repo.On("DeleteByReferenceObjectID", ctx, tenantID, model.EventAPIFetchRequestReference, id).Return(testErr).Once()
				return repo
			},
			UIDServiceFn: func() *automock.UIDService {
				svc := &automock.UIDService{}
				return svc
			},
			FetchRequestServiceFn: func() *automock.FetchRequestService {
				svc := &automock.FetchRequestService{}
				return svc
			},
			InputID:     "foo",
			Input:       modelInput,
			ExpectedErr: testErr,
		},
		{
			Name: "Fetch Request Creation Error",
			RepositoryFn: func() *automock.EventAPIRepository {
				repo := &automock.EventAPIRepository{}
				repo.On("GetByID", ctx, tenantID, "foo").Return(eventAPIDefinitionModel, nil).Once()
				return repo
			},
			FetchRequestRepoFn: func() *automock.FetchRequestRepository {
				repo := &automock.FetchRequestRepository{}
				repo.On("DeleteByReferenceObjectID", ctx, tenantID, model.EventAPIFetchRequestReference, id).Return(nil).Once()
				repo.On("Create", ctx, modelFr).Return(testErr).Once()
				return repo
			},
			UIDServiceFn: func() *automock.UIDService {
				svc := &automock.UIDService{}
				svc.On("Generate").Return(frID).Once()
				return svc
			},
			FetchRequestServiceFn: func() *automock.FetchRequestService {
				svc := &automock.FetchRequestService{}
				return svc
			},
			InputID:     "foo",
			Input:       modelInput,
			ExpectedErr: testErr,
		},
		{
			Name: "Get Error",
			FetchRequestRepoFn: func() *automock.FetchRequestRepository {
				repo := &automock.FetchRequestRepository{}
				return repo
			},
			UIDServiceFn: func() *automock.UIDService {
				svc := &automock.UIDService{}
				return svc
			},
			RepositoryFn: func() *automock.EventAPIRepository {
				repo := &automock.EventAPIRepository{}
				repo.On("GetByID", ctx, tenantID, id).Return(nil, testErr).Once()
				return repo
			},
			FetchRequestServiceFn: func() *automock.FetchRequestService {
				svc := &automock.FetchRequestService{}
				return svc
			},
			InputID:     "foo",
			Input:       modelInput,
			ExpectedErr: testErr,
		},
		{
			Name: "Success when fetch request failed",
			RepositoryFn: func() *automock.EventAPIRepository {
				repo := &automock.EventAPIRepository{}
				repo.On("GetByID", ctx, tenantID, id).Return(eventAPIDefinitionModel, nil).Once()
				repo.On("Update", ctx, inputEventAPIDefinitionModel).Return(nil).Once()
				return repo
			},
			FetchRequestRepoFn: func() *automock.FetchRequestRepository {
				repo := &automock.FetchRequestRepository{}
				repo.On("DeleteByReferenceObjectID", ctx, tenantID, model.EventAPIFetchRequestReference, id).Return(nil).Once()
				repo.On("Create", ctx, modelFr).Return(nil).Once()

				return repo
			},
			UIDServiceFn: func() *automock.UIDService {
				svc := &automock.UIDService{}
				svc.On("Generate").Return(frID).Once()
				return svc
			},
			FetchRequestServiceFn: func() *automock.FetchRequestService {
				svc := &automock.FetchRequestService{}
				svc.On("HandleSpec", ctx, modelFr).Return(nil)
				return svc
			},
			InputID:     "foo",
			Input:       modelInput,
			ExpectedErr: nil,
		},
	}

	for _, testCase := range testCases {
		t.Run(fmt.Sprintf("%s", testCase.Name), func(t *testing.T) {
			// given
			repo := testCase.RepositoryFn()
			fetchRequestRepo := testCase.FetchRequestRepoFn()
			uidSvc := testCase.UIDServiceFn()
			fetchRequestSvc := testCase.FetchRequestServiceFn()

			svc := eventdef.NewService(repo, fetchRequestRepo, uidSvc, fetchRequestSvc)
			svc.SetTimestampGen(func() time.Time { return timestamp })

			// when
			err := svc.Update(ctx, testCase.InputID, testCase.Input)

			// then
			if testCase.ExpectedErr == nil {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
				assert.Contains(t, err.Error(), testCase.ExpectedErr.Error())
			}

			repo.AssertExpectations(t)
			fetchRequestRepo.AssertExpectations(t)
			uidSvc.AssertExpectations(t)
		})
	}
	t.Run("Error when tenant not in context", func(t *testing.T) {
		svc := eventdef.NewService(nil, nil, nil, nil)
		// WHEN
		err := svc.Update(context.TODO(), "", model.EventDefinitionInput{})
		// THEN
		require.Error(t, err)
		assert.Contains(t, err.Error(), "cannot read tenant from context")
	})
}

func TestService_Delete(t *testing.T) {
	// given
	testErr := errors.New("Test error")

	id := "foo"

	ctx := context.TODO()
	ctx = tenant.SaveToContext(ctx, tenantID, externalTenantID)

	testCases := []struct {
		Name         string
		RepositoryFn func() *automock.EventAPIRepository
		Input        model.EventDefinitionInput
		InputID      string
		ExpectedErr  error
	}{
		{
			Name: "Success",
			RepositoryFn: func() *automock.EventAPIRepository {
				repo := &automock.EventAPIRepository{}
				repo.On("Delete", ctx, tenantID, id).Return(nil).Once()
				return repo
			},
			InputID:     id,
			ExpectedErr: nil,
		},
		{
			Name: "Delete Error",
			RepositoryFn: func() *automock.EventAPIRepository {
				repo := &automock.EventAPIRepository{}
				repo.On("Delete", ctx, tenantID, id).Return(testErr).Once()
				return repo
			},
			InputID:     id,
			ExpectedErr: testErr,
		},
	}

	for _, testCase := range testCases {
		t.Run(fmt.Sprintf("%s", testCase.Name), func(t *testing.T) {
			// given
			repo := testCase.RepositoryFn()

			svc := eventdef.NewService(repo, nil, nil, nil)

			// when
			err := svc.Delete(ctx, testCase.InputID)

			// then
			if testCase.ExpectedErr == nil {
				require.NoError(t, err)
			} else {
				assert.Contains(t, err.Error(), testCase.ExpectedErr.Error())
			}

			repo.AssertExpectations(t)
		})
	}
	t.Run("Error when tenant not in context", func(t *testing.T) {
		svc := eventdef.NewService(nil, nil, nil, nil)
		// WHEN
		err := svc.Delete(context.TODO(), "")
		// THEN
		require.Error(t, err)
		assert.Contains(t, err.Error(), "cannot read tenant from context")
	})
}

func TestService_RefetchSpec(t *testing.T) {
	// given
	testErr := errors.New("Test error")

	apiID := "foo"

	ctx := context.TODO()
	ctx = tenant.SaveToContext(ctx, tenantID, externalTenantID)

	dataBytes := "data"
	modelAPISpec := &model.EventSpec{
		Data: &dataBytes,
	}

	modelAPIDefinition := &model.EventDefinition{
		Spec: modelAPISpec,
	}

	timestamp := time.Now()
	fr := &model.FetchRequest{
		Status: &model.FetchRequestStatus{
			Condition: model.FetchRequestStatusConditionInitial,
			Timestamp: timestamp,
		},
	}

	testCases := []struct {
		Name               string
		RepositoryFn       func() *automock.EventAPIRepository
		FetchRequestRepoFn func() *automock.FetchRequestRepository
		FetchRequestSvcFn  func() *automock.FetchRequestService
		ExpectedAPISpec    *model.EventSpec
		ExpectedErr        error
	}{
		{
			Name: "Success",
			RepositoryFn: func() *automock.EventAPIRepository {
				repo := &automock.EventAPIRepository{}
				repo.On("GetByID", ctx, tenantID, apiID).Return(modelAPIDefinition, nil).Once()
				repo.On("Update", ctx, modelAPIDefinition).Return(nil).Once()

				return repo
			},
			FetchRequestRepoFn: func() *automock.FetchRequestRepository {
				repo := &automock.FetchRequestRepository{}
				repo.On("GetByReferenceObjectID", ctx, tenantID, model.EventAPIFetchRequestReference, apiID).Return(nil, nil)
				return repo
			},
			FetchRequestSvcFn: func() *automock.FetchRequestService {
				svc := &automock.FetchRequestService{}
				return svc
			},
			ExpectedAPISpec: modelAPISpec,
			ExpectedErr:     nil,
		},
		{
			Name: "Success - fetched Event API Spec",
			RepositoryFn: func() *automock.EventAPIRepository {
				repo := &automock.EventAPIRepository{}
				repo.On("GetByID", ctx, tenantID, apiID).Return(modelAPIDefinition, nil).Once()
				repo.On("Update", ctx, modelAPIDefinition).Return(nil).Once()
				return repo
			},
			FetchRequestRepoFn: func() *automock.FetchRequestRepository {
				repo := &automock.FetchRequestRepository{}
				repo.On("GetByReferenceObjectID", ctx, tenantID, model.EventAPIFetchRequestReference, apiID).Return(fr, nil)
				return repo
			},
			FetchRequestSvcFn: func() *automock.FetchRequestService {
				svc := &automock.FetchRequestService{}
				svc.On("HandleSpec", ctx, fr).Return(&dataBytes)
				return svc
			},
			ExpectedAPISpec: modelAPISpec,
			ExpectedErr:     nil,
		},
		{
			Name: "Get from repository error",
			RepositoryFn: func() *automock.EventAPIRepository {
				repo := &automock.EventAPIRepository{}
				repo.On("GetByID", ctx, tenantID, apiID).Return(nil, testErr).Once()
				return repo
			},
			FetchRequestRepoFn: func() *automock.FetchRequestRepository {
				repo := &automock.FetchRequestRepository{}
				return repo
			},
			FetchRequestSvcFn: func() *automock.FetchRequestService {
				svc := &automock.FetchRequestService{}
				return svc
			},
			ExpectedAPISpec: nil,
			ExpectedErr:     testErr,
		},
		{
			Name: "Get fetch request error",
			RepositoryFn: func() *automock.EventAPIRepository {
				repo := &automock.EventAPIRepository{}
				repo.On("GetByID", ctx, tenantID, apiID).Return(modelAPIDefinition, nil).Once()
				return repo
			},
			FetchRequestRepoFn: func() *automock.FetchRequestRepository {
				repo := &automock.FetchRequestRepository{}
				repo.On("GetByReferenceObjectID", ctx, tenantID, model.EventAPIFetchRequestReference, apiID).Return(nil, testErr)
				return repo
			},
			FetchRequestSvcFn: func() *automock.FetchRequestService {
				svc := &automock.FetchRequestService{}
				return svc
			},
			ExpectedAPISpec: nil,
			ExpectedErr:     errors.Wrapf(testErr, "while getting FetchRequest by Event Definition ID %s", apiID),
		},
		{
			Name: "Error when updating Event Event Definition failed",
			RepositoryFn: func() *automock.EventAPIRepository {
				repo := &automock.EventAPIRepository{}
				repo.On("GetByID", ctx, tenantID, apiID).Return(modelAPIDefinition, nil).Once()
				repo.On("Update", ctx, modelAPIDefinition).Return(testErr).Once()
				return repo
			},
			FetchRequestRepoFn: func() *automock.FetchRequestRepository {
				repo := &automock.FetchRequestRepository{}
				repo.On("GetByReferenceObjectID", ctx, tenantID, model.EventAPIFetchRequestReference, apiID).Return(nil, nil)
				return repo
			},
			FetchRequestSvcFn: func() *automock.FetchRequestService {
				svc := &automock.FetchRequestService{}
				return svc
			},
			ExpectedAPISpec: nil,
			ExpectedErr:     errors.Wrap(testErr, "while updating event api with event api spec"),
		},
	}

	for _, testCase := range testCases {
		t.Run(fmt.Sprintf("%s", testCase.Name), func(t *testing.T) {
			// given
			repo := testCase.RepositoryFn()
			frRepo := testCase.FetchRequestRepoFn()
			frSvc := testCase.FetchRequestSvcFn()

			svc := eventdef.NewService(repo, frRepo, nil, frSvc)

			// when
			result, err := svc.RefetchAPISpec(ctx, apiID)

			// then
			assert.Equal(t, testCase.ExpectedAPISpec, result)

			if testCase.ExpectedErr != nil {
				assert.Equal(t, testCase.ExpectedErr.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}
			repo.AssertExpectations(t)
		})
	}
	t.Run("Error when tenant not in context", func(t *testing.T) {
		svc := eventdef.NewService(nil, nil, nil, nil)
		// WHEN
		_, err := svc.RefetchAPISpec(context.TODO(), "")
		// THEN
		require.Error(t, err)
		assert.Contains(t, err.Error(), "cannot read tenant from context")
	})
}

func TestService_GetFetchRequest(t *testing.T) {
	// given
	ctx := context.TODO()
	ctx = tenant.SaveToContext(ctx, tenantID, externalTenantID)

	testErr := errors.New("Test error")

	id := "foo"
	refID := "doc-id"
	frURL := "foo.bar"
	timestamp := time.Now()

	fetchRequestModel := fixModelFetchRequest(id, frURL, timestamp)

	testCases := []struct {
		Name                 string
		RepositoryFn         func() *automock.EventAPIRepository
		FetchRequestRepoFn   func() *automock.FetchRequestRepository
		ExpectedFetchRequest *model.FetchRequest
		ExpectedErrMessage   string
	}{
		{
			Name: "Success",
			RepositoryFn: func() *automock.EventAPIRepository {
				repo := &automock.EventAPIRepository{}
				repo.On("Exists", ctx, tenantID, refID).Return(true, nil).Once()
				return repo
			},
			FetchRequestRepoFn: func() *automock.FetchRequestRepository {
				repo := &automock.FetchRequestRepository{}
				repo.On("GetByReferenceObjectID", ctx, tenantID, model.EventAPIFetchRequestReference, refID).Return(fetchRequestModel, nil).Once()
				return repo
			},
			ExpectedFetchRequest: fetchRequestModel,
			ExpectedErrMessage:   "",
		},
		{
			Name: "Success - Not Found",
			RepositoryFn: func() *automock.EventAPIRepository {
				repo := &automock.EventAPIRepository{}
				repo.On("Exists", ctx, tenantID, refID).Return(true, nil).Once()
				return repo
			},
			FetchRequestRepoFn: func() *automock.FetchRequestRepository {
				repo := &automock.FetchRequestRepository{}
				repo.On("GetByReferenceObjectID", ctx, tenantID, model.EventAPIFetchRequestReference, refID).Return(nil, apperrors.NewNotFoundError(resource.EventDefinition, "")).Once()
				return repo
			},
			ExpectedFetchRequest: nil,
			ExpectedErrMessage:   "",
		},
		{
			Name: "Error - Get FetchRequest",
			RepositoryFn: func() *automock.EventAPIRepository {
				repo := &automock.EventAPIRepository{}
				repo.On("Exists", ctx, tenantID, refID).Return(true, nil).Once()

				return repo
			},
			FetchRequestRepoFn: func() *automock.FetchRequestRepository {
				repo := &automock.FetchRequestRepository{}
				repo.On("GetByReferenceObjectID", ctx, tenantID, model.EventAPIFetchRequestReference, refID).Return(nil, testErr).Once()
				return repo
			},
			ExpectedFetchRequest: nil,
			ExpectedErrMessage:   testErr.Error(),
		},
		{
			Name: "Error - Event Definition doesn't exist",
			RepositoryFn: func() *automock.EventAPIRepository {
				repo := &automock.EventAPIRepository{}
				repo.On("Exists", ctx, tenantID, refID).Return(false, nil).Once()

				return repo
			},
			FetchRequestRepoFn: func() *automock.FetchRequestRepository {
				repo := &automock.FetchRequestRepository{}
				return repo
			},
			ExpectedErrMessage:   "Event Definition with ID doc-id doesn't exist",
			ExpectedFetchRequest: nil,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			repo := testCase.RepositoryFn()
			fetchRequestRepo := testCase.FetchRequestRepoFn()
			svc := eventdef.NewService(repo, fetchRequestRepo, nil, nil)

			// when
			l, err := svc.GetFetchRequest(ctx, refID)

			// then
			if testCase.ExpectedErrMessage == "" {
				require.NoError(t, err)
				assert.Equal(t, l, testCase.ExpectedFetchRequest)
			} else {
				require.Error(t, err)
				assert.Contains(t, err.Error(), testCase.ExpectedErrMessage)
			}

			repo.AssertExpectations(t)
			fetchRequestRepo.AssertExpectations(t)
		})
	}

	t.Run("Returns error on loading tenant", func(t *testing.T) {
		svc := eventdef.NewService(nil, nil, nil, nil)
		// when
		_, err := svc.GetFetchRequest(context.TODO(), "dd")
		assert.True(t, apperrors.IsCannotReadTenant(err))
	})
}