package runtime_test

import (
	"context"
	"database/sql/driver"
	"fmt"
	"regexp"
	"testing"

	"github.com/google/uuid"
	"github.com/kyma-incubator/compass/components/director/pkg/pagination"
	"github.com/kyma-incubator/compass/components/director/pkg/str"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/kyma-incubator/compass/components/director/internal/domain/runtime"
	"github.com/kyma-incubator/compass/components/director/internal/domain/runtime/automock"
	"github.com/kyma-incubator/compass/components/director/internal/labelfilter"
	"github.com/kyma-incubator/compass/components/director/internal/model"
	"github.com/kyma-incubator/compass/components/director/internal/repo/testdb"
	"github.com/kyma-incubator/compass/components/director/pkg/persistence"
	"github.com/stretchr/testify/require"
)

func TestPgRepository_GetByID(t *testing.T) {
	rtModel := fixDetailedModelRuntime(t, "foo", "Foo", "Lorem ipsum")
	rtEntity := fixDetailedEntityRuntime(t, "foo", "Foo", "Lorem ipsum")

	suite := testdb.RepoGetTestSuite{
		Name: "Get Runtime By ID",
		SqlQueryDetails: []testdb.SqlQueryDetails{
			{
				Query:    regexp.QuoteMeta(`SELECT id, name, description, status_condition, status_timestamp, creation_timestamp FROM public.runtimes WHERE id = $1 AND (id IN (SELECT id FROM tenant_runtimes WHERE tenant_id = $2))`),
				Args:     []driver.Value{runtimeID, tenantID},
				IsSelect: true,
				ValidRowsProvider: func() []*sqlmock.Rows {
					return []*sqlmock.Rows{sqlmock.NewRows(fixColumns).AddRow(rtModel.ID, rtModel.Name, rtModel.Description, rtModel.Status.Condition, rtModel.Status.Timestamp, rtModel.CreationTimestamp)}
				},
				InvalidRowsProvider: func() []*sqlmock.Rows {
					return []*sqlmock.Rows{sqlmock.NewRows(fixColumns)}
				},
			},
		},
		ConverterMockProvider: func() testdb.Mock {
			return &automock.EntityConverter{}
		},
		RepoConstructorFunc:       runtime.NewRepository,
		ExpectedModelEntity:       rtModel,
		ExpectedDBEntity:          rtEntity,
		MethodArgs:                []interface{}{tenantID, runtimeID},
		DisableConverterErrorTest: true,
	}

	suite.Run(t)
}

func TestPgRepository_GetByFiltersAndID(t *testing.T) {
	rtModel := fixDetailedModelRuntime(t, "foo", "Foo", "Lorem ipsum")
	rtEntity := fixDetailedEntityRuntime(t, "foo", "Foo", "Lorem ipsum")

	suite := testdb.RepoGetTestSuite{
		Name: "Get Runtime By Filters and ID",
		SqlQueryDetails: []testdb.SqlQueryDetails{
			{
				Query: regexp.QuoteMeta(`SELECT id, name, description, status_condition, status_timestamp, creation_timestamp FROM public.runtimes WHERE id = $1 
												AND id IN (SELECT "runtime_id" FROM public.labels WHERE "runtime_id" IS NOT NULL AND (id IN (SELECT id FROM runtime_labels_tenants WHERE tenant_id = $2)) AND "key" = $3 AND "value" ?| array[$4]) 
												AND (id IN (SELECT id FROM tenant_runtimes WHERE tenant_id = $5))`),
				Args:     []driver.Value{runtimeID, tenantID, model.ScenariosKey, "scenario", tenantID},
				IsSelect: true,
				ValidRowsProvider: func() []*sqlmock.Rows {
					return []*sqlmock.Rows{sqlmock.NewRows(fixColumns).AddRow(rtModel.ID, rtModel.Name, rtModel.Description, rtModel.Status.Condition, rtModel.Status.Timestamp, rtModel.CreationTimestamp)}
				},
				InvalidRowsProvider: func() []*sqlmock.Rows {
					return []*sqlmock.Rows{sqlmock.NewRows(fixColumns)}
				},
			},
		},
		ConverterMockProvider: func() testdb.Mock {
			return &automock.EntityConverter{}
		},
		RepoConstructorFunc:       runtime.NewRepository,
		ExpectedModelEntity:       rtModel,
		ExpectedDBEntity:          rtEntity,
		MethodName:                "GetByFiltersAndID",
		MethodArgs:                []interface{}{tenantID, runtimeID, []*labelfilter.LabelFilter{labelfilter.NewForKeyWithQuery(model.ScenariosKey, `$[*] ? ( @ == "scenario" )`)}},
		DisableConverterErrorTest: true,
	}

	suite.Run(t)
}

func TestPgRepository_GetByFiltersGlobal_ShouldReturnRuntimeModelForRuntimeEntity(t *testing.T) {
	// given
	rtModel := fixDetailedModelRuntime(t, "foo", "Foo", "Lorem ipsum")
	rtEntity := fixDetailedEntityRuntime(t, "foo", "Foo", "Lorem ipsum")

	mockConverter := &automock.EntityConverter{}
	mockConverter.On("FromEntity", rtEntity).Return(rtModel, nil).Once()

	sqlxDB, sqlMock := testdb.MockDatabase(t)
	defer sqlMock.AssertExpectations(t)

	rows := sqlmock.NewRows([]string{"id", "name", "description", "status_condition", "status_timestamp", "creation_timestamp"}).
		AddRow(rtModel.ID, rtModel.Name, rtModel.Description, rtModel.Status.Condition, rtModel.Status.Timestamp, rtModel.CreationTimestamp)

	sqlMock.ExpectQuery(`^SELECT (.+) FROM public.runtimes WHERE id IN \(SELECT "runtime_id" FROM public\.labels WHERE "runtime_id" IS NOT NULL AND "key" = \$1\)$`).
		WithArgs("someKey").
		WillReturnRows(rows)

	ctx := persistence.SaveToContext(context.TODO(), sqlxDB)

	pgRepository := runtime.NewRepository(mockConverter)

	// when
	filters := []*labelfilter.LabelFilter{labelfilter.NewForKey("someKey")}
	modelRuntime, err := pgRepository.GetByFiltersGlobal(ctx, filters)

	//then
	require.NoError(t, err)
	require.Equal(t, rtModel, modelRuntime)
	mockConverter.AssertExpectations(t)
}

func TestPgRepository_GetOldestForFilters(t *testing.T) {
	rtModel := fixDetailedModelRuntime(t, "foo", "Foo", "Lorem ipsum")
	rtEntity := fixDetailedEntityRuntime(t, "foo", "Foo", "Lorem ipsum")

	suite := testdb.RepoGetTestSuite{
		Name: "Get Oldest Runtime By Filters",
		SqlQueryDetails: []testdb.SqlQueryDetails{
			{
				Query: regexp.QuoteMeta(`SELECT id, name, description, status_condition, status_timestamp, creation_timestamp FROM public.runtimes WHERE  
												id IN (SELECT "runtime_id" FROM public.labels WHERE "runtime_id" IS NOT NULL AND (id IN (SELECT id FROM runtime_labels_tenants WHERE tenant_id = $1)) AND "key" = $2 AND "value" ?| array[$3]) 
												AND (id IN (SELECT id FROM tenant_runtimes WHERE tenant_id = $4)) ORDER BY creation_timestamp ASC`),
				Args:     []driver.Value{tenantID, model.ScenariosKey, "scenario", tenantID},
				IsSelect: true,
				ValidRowsProvider: func() []*sqlmock.Rows {
					return []*sqlmock.Rows{sqlmock.NewRows(fixColumns).AddRow(rtModel.ID, rtModel.Name, rtModel.Description, rtModel.Status.Condition, rtModel.Status.Timestamp, rtModel.CreationTimestamp)}
				},
				InvalidRowsProvider: func() []*sqlmock.Rows {
					return []*sqlmock.Rows{sqlmock.NewRows(fixColumns)}
				},
			},
		},
		ConverterMockProvider: func() testdb.Mock {
			return &automock.EntityConverter{}
		},
		RepoConstructorFunc:       runtime.NewRepository,
		ExpectedModelEntity:       rtModel,
		ExpectedDBEntity:          rtEntity,
		MethodName:                "GetOldestForFilters",
		MethodArgs:                []interface{}{tenantID, []*labelfilter.LabelFilter{labelfilter.NewForKeyWithQuery(model.ScenariosKey, `$[*] ? ( @ == "scenario" )`)}},
		DisableConverterErrorTest: true,
	}

	suite.Run(t)
}

func TestPgRepository_ListByFiltersGlobal(t *testing.T) {
	// GIVEN
	runtime1ID := uuid.New().String()
	runtime2ID := uuid.New().String()
	runtimeEntity1 := fixDetailedEntityRuntime(t, runtime1ID, "Runtime 1", "Runtime desc 1")
	runtimeEntity2 := fixDetailedEntityRuntime(t, runtime2ID, "Runtime 2", "Runtime desc 2")

	runtimeModel1 := fixModelRuntime(t, runtime1ID, tenantID, "Runtime 1", "Runtime desc 1")
	runtimeModel2 := fixModelRuntime(t, runtime2ID, tenantID, "Runtime 2", "Runtime desc 2")

	mockConverter := &automock.EntityConverter{}
	mockConverter.On("FromEntity", runtimeEntity1).Return(runtimeModel1)
	mockConverter.On("FromEntity", runtimeEntity2).Return(runtimeModel2)

	sqlxDB, sqlMock := testdb.MockDatabase(t)
	defer sqlMock.AssertExpectations(t)

	rows := sqlmock.NewRows([]string{"id", "name", "description", "status_condition", "status_timestamp", "creation_timestamp"}).
		AddRow(runtime1ID, runtimeModel1.Name, runtimeModel1.Description, runtimeModel1.Status.Condition, runtimeModel1.CreationTimestamp, runtimeModel1.CreationTimestamp).
		AddRow(runtime2ID, runtimeModel2.Name, runtimeModel2.Description, runtimeModel2.Status.Condition, runtimeModel2.CreationTimestamp, runtimeModel2.CreationTimestamp)

	sqlMock.ExpectQuery(`^SELECT (.+) FROM public.runtimes WHERE id IN \(SELECT "runtime_id" FROM public\.labels WHERE "runtime_id" IS NOT NULL AND "key" = \$1 AND "value" \@\> \$2\ INTERSECT SELECT "runtime_id" FROM public\.labels WHERE "runtime_id" IS NOT NULL AND "key" = \$3 AND "value" \@\> \$4\)$`).
		WithArgs("someKey", "someValue", "someKey2", "someValue2").
		WillReturnRows(rows)

	ctx := persistence.SaveToContext(context.TODO(), sqlxDB)

	pgRepository := runtime.NewRepository(mockConverter)

	filters := []*labelfilter.LabelFilter{
		{
			Key:   "someKey",
			Query: str.Ptr(`someValue`),
		},
		{
			Key:   "someKey2",
			Query: str.Ptr(`someValue2`),
		},
	}
	// WHEN
	modelRuntimes, err := pgRepository.ListByFiltersGlobal(ctx, filters)

	// THEN
	require.NoError(t, err)
	require.NotNil(t, modelRuntimes)
	require.NoError(t, sqlMock.ExpectationsWereMet())

	require.Len(t, modelRuntimes, 2)
	require.Equal(t, runtimeModel1, modelRuntimes[0])
	require.Equal(t, runtimeModel2, modelRuntimes[1])
}

func TestPgRepository_List(t *testing.T) {
	runtime1ID := "aec0e9c5-06da-4625-9f8a-bda17ab8c3b9"
	runtime2ID := "ccdbef8f-b97a-490c-86e2-2bab2862a6e4"
	runtimeEntity1 := fixDetailedEntityRuntime(t, runtime1ID, "Runtime 1", "Runtime desc 1")
	runtimeEntity2 := fixDetailedEntityRuntime(t, runtime2ID, "Runtime 2", "Runtime desc 2")

	runtimeModel1 := fixModelRuntime(t, runtime1ID, tenantID, "Runtime 1", "Runtime desc 1")
	runtimeModel2 := fixModelRuntime(t, runtime2ID, tenantID, "Runtime 2", "Runtime desc 2")

	suite := testdb.RepoListPageableTestSuite{
		Name: "List Runtimes",
		SqlQueryDetails: []testdb.SqlQueryDetails{
			{
				Query: regexp.QuoteMeta(`SELECT id, name, description, status_condition, status_timestamp, creation_timestamp FROM public.runtimes
												WHERE id IN (SELECT "runtime_id" FROM public.labels WHERE "runtime_id" IS NOT NULL AND (id IN (SELECT id FROM runtime_labels_tenants WHERE tenant_id = $1)) AND "key" = $2 AND "value" ?| array[$3])
												AND (id IN (SELECT id FROM tenant_runtimes WHERE tenant_id = $4)) ORDER BY name LIMIT 2 OFFSET 0`),
				Args:     []driver.Value{tenantID, model.ScenariosKey, "scenario", tenantID},
				IsSelect: true,
				ValidRowsProvider: func() []*sqlmock.Rows {
					return []*sqlmock.Rows{sqlmock.NewRows(fixColumns).
						AddRow(runtimeEntity1.ID, runtimeEntity1.Name, runtimeEntity1.Description, runtimeEntity1.StatusCondition, runtimeEntity1.StatusTimestamp, runtimeEntity1.CreationTimestamp).
						AddRow(runtimeEntity2.ID, runtimeEntity2.Name, runtimeEntity2.Description, runtimeEntity2.StatusCondition, runtimeEntity2.StatusTimestamp, runtimeEntity2.CreationTimestamp),
					}
				},
			},
			{
				Query: regexp.QuoteMeta(`SELECT COUNT(*) FROM public.runtimes
												WHERE id IN (SELECT "runtime_id" FROM public.labels WHERE "runtime_id" IS NOT NULL AND (id IN (SELECT id FROM runtime_labels_tenants WHERE tenant_id = $1)) AND "key" = $2 AND "value" ?| array[$3])
												AND (id IN (SELECT id FROM tenant_runtimes WHERE tenant_id = $4))`),
				Args:     []driver.Value{tenantID, model.ScenariosKey, "scenario", tenantID},
				IsSelect: true,
				ValidRowsProvider: func() []*sqlmock.Rows {
					return []*sqlmock.Rows{sqlmock.NewRows([]string{"count"}).AddRow(2)}
				},
			},
		},
		Pages: []testdb.PageDetails{
			{
				ExpectedModelEntities: []interface{}{runtimeModel1, runtimeModel2},
				ExpectedDBEntities:    []interface{}{runtimeEntity1, runtimeEntity2},
				ExpectedPage: &model.RuntimePage{
					Data: []*model.Runtime{runtimeModel1, runtimeModel2},
					PageInfo: &pagination.Page{
						StartCursor: "",
						EndCursor:   "",
						HasNextPage: false,
					},
					TotalCount: 2,
				},
			},
		},
		ConverterMockProvider: func() testdb.Mock {
			return &automock.EntityConverter{}
		},
		RepoConstructorFunc:       runtime.NewRepository,
		MethodArgs:                []interface{}{tenantID, []*labelfilter.LabelFilter{labelfilter.NewForKeyWithQuery(model.ScenariosKey, `$[*] ? ( @ == "scenario" )`)}, 2, ""},
		MethodName:                "List",
		DisableConverterErrorTest: true,
	}

	suite.Run(t)
}

func TestPgRepository_ListAll(t *testing.T) {
	runtime1ID := "aec0e9c5-06da-4625-9f8a-bda17ab8c3b9"
	runtime2ID := "ccdbef8f-b97a-490c-86e2-2bab2862a6e4"
	runtimeEntity1 := fixDetailedEntityRuntime(t, runtime1ID, "Runtime 1", "Runtime desc 1")
	runtimeEntity2 := fixDetailedEntityRuntime(t, runtime2ID, "Runtime 2", "Runtime desc 2")

	runtimeModel1 := fixModelRuntime(t, runtime1ID, tenantID, "Runtime 1", "Runtime desc 1")
	runtimeModel2 := fixModelRuntime(t, runtime2ID, tenantID, "Runtime 2", "Runtime desc 2")

	suite := testdb.RepoListTestSuite{
		Name: "List Runtimes Without Paging",
		SqlQueryDetails: []testdb.SqlQueryDetails{
			{
				Query: regexp.QuoteMeta(`SELECT id, name, description, status_condition, status_timestamp, creation_timestamp FROM public.runtimes 
												WHERE id IN (SELECT "runtime_id" FROM public.labels WHERE "runtime_id" IS NOT NULL AND (id IN (SELECT id FROM runtime_labels_tenants WHERE tenant_id = $1)) AND "key" = $2 AND "value" ?| array[$3])
												AND (id IN (SELECT id FROM tenant_runtimes WHERE tenant_id = $4))`),
				Args:     []driver.Value{tenantID, model.ScenariosKey, "scenario", tenantID},
				IsSelect: true,
				ValidRowsProvider: func() []*sqlmock.Rows {
					return []*sqlmock.Rows{sqlmock.NewRows(fixColumns).
						AddRow(runtimeEntity1.ID, runtimeEntity1.Name, runtimeEntity1.Description, runtimeEntity1.StatusCondition, runtimeEntity1.StatusTimestamp, runtimeEntity1.CreationTimestamp).
						AddRow(runtimeEntity2.ID, runtimeEntity2.Name, runtimeEntity2.Description, runtimeEntity2.StatusCondition, runtimeEntity2.StatusTimestamp, runtimeEntity2.CreationTimestamp),
					}
				},
				InvalidRowsProvider: func() []*sqlmock.Rows {
					return []*sqlmock.Rows{sqlmock.NewRows(fixColumns)}
				},
			},
		},
		ConverterMockProvider: func() testdb.Mock {
			return &automock.EntityConverter{}
		},
		RepoConstructorFunc:       runtime.NewRepository,
		ExpectedModelEntities:     []interface{}{runtimeModel1, runtimeModel2},
		ExpectedDBEntities:        []interface{}{runtimeEntity1, runtimeEntity2},
		MethodArgs:                []interface{}{tenantID, []*labelfilter.LabelFilter{labelfilter.NewForKeyWithQuery(model.ScenariosKey, `$[*] ? ( @ == "scenario" )`)}},
		MethodName:                "ListAll",
		DisableConverterErrorTest: true,
	}

	suite.Run(t)
}

func TestPgRepository_Create(t *testing.T) {
	var nilRtModel *model.Runtime
	rtModel := fixDetailedModelRuntime(t, "foo", "Foo", "Lorem ipsum")
	rtEntity := fixDetailedEntityRuntime(t, "foo", "Foo", "Lorem ipsum")

	suite := testdb.RepoCreateTestSuite{
		Name: "Generic Create Runtime",
		SqlQueryDetails: []testdb.SqlQueryDetails{
			{
				Query:       regexp.QuoteMeta(`INSERT INTO public.runtimes ( id, name, description, status_condition, status_timestamp, creation_timestamp ) VALUES ( ?, ?, ?, ?, ?, ? )`),
				Args:        []driver.Value{rtModel.ID, rtModel.Name, rtModel.Description, rtModel.Status.Condition, rtModel.Status.Timestamp, rtModel.CreationTimestamp},
				ValidResult: sqlmock.NewResult(-1, 1),
			},
			{
				Query:       regexp.QuoteMeta(`INSERT INTO tenant_runtimes ( tenant_id, id, owner ) VALUES ( ?, ?, ? )`),
				Args:        []driver.Value{tenantID, rtModel.ID, true},
				ValidResult: sqlmock.NewResult(-1, 1),
			},
		},
		ConverterMockProvider: func() testdb.Mock {
			return &automock.EntityConverter{}
		},
		RepoConstructorFunc: runtime.NewRepository,
		ModelEntity:         rtModel,
		DBEntity:            rtEntity,
		NilModelEntity:      nilRtModel,
		TenantID:            tenantID,
		IsTopLevelEntity:    true,
	}

	suite.Run(t)
}

func TestPgRepository_Update(t *testing.T) {
	var nilRtModel *model.Runtime
	rtModel := fixDetailedModelRuntime(t, "foo", "Foo", "Lorem ipsum")
	rtEntity := fixDetailedEntityRuntime(t, "foo", "Foo", "Lorem ipsum")

	suite := testdb.RepoUpdateTestSuite{
		Name: "Update Runtime",
		SqlQueryDetails: []testdb.SqlQueryDetails{
			{
				Query:         regexp.QuoteMeta(fmt.Sprintf(`UPDATE public.runtimes SET name = ?, description = ?, status_condition = ?, status_timestamp = ? WHERE id = ? AND (id IN (SELECT id FROM tenant_runtimes WHERE tenant_id = '%s' AND owner = true))`, tenantID)),
				Args:          []driver.Value{rtModel.Name, rtModel.Description, rtModel.Status.Condition, rtModel.Status.Timestamp, rtModel.ID},
				ValidResult:   sqlmock.NewResult(-1, 1),
				InvalidResult: sqlmock.NewResult(-1, 0),
			},
		},
		ConverterMockProvider: func() testdb.Mock {
			return &automock.EntityConverter{}
		},
		RepoConstructorFunc: runtime.NewRepository,
		ModelEntity:         rtModel,
		DBEntity:            rtEntity,
		NilModelEntity:      nilRtModel,
		TenantID:            tenantID,
	}

	suite.Run(t)
}

func TestPgRepository_Delete(t *testing.T) {
	suite := testdb.RepoDeleteTestSuite{
		Name: "Runtime Delete",
		SqlQueryDetails: []testdb.SqlQueryDetails{
			{
				Query:    regexp.QuoteMeta(`SELECT id FROM public.runtimes WHERE id = $1 AND (id IN (SELECT id FROM tenant_runtimes WHERE tenant_id = $2 AND owner = true))`),
				Args:     []driver.Value{runtimeID, tenantID},
				IsSelect: true,
				ValidRowsProvider: func() []*sqlmock.Rows {
					return []*sqlmock.Rows{sqlmock.NewRows([]string{"id"}).AddRow(runtimeID)}
				},
				InvalidRowsProvider: func() []*sqlmock.Rows {
					return []*sqlmock.Rows{sqlmock.NewRows([]string{"id"}).AddRow(runtimeID).AddRow("secondID")}
				},
			},
			{
				Query:       regexp.QuoteMeta(`DELETE FROM tenant_runtimes WHERE id IN ($1)`),
				Args:        []driver.Value{runtimeID},
				ValidResult: sqlmock.NewResult(-1, 1),
			},
		},
		ConverterMockProvider: func() testdb.Mock {
			return &automock.EntityConverter{}
		},
		RepoConstructorFunc: runtime.NewRepository,
		MethodArgs:          []interface{}{tenantID, runtimeID},
		IsTopLeveEntity:     true,
	}

	suite.Run(t)
}

func TestPgRepository_Exist(t *testing.T) {
	suite := testdb.RepoExistTestSuite{
		Name: "Runtime Exists",
		SqlQueryDetails: []testdb.SqlQueryDetails{
			{
				Query:    regexp.QuoteMeta(`SELECT 1 FROM public.runtimes WHERE id = $1 AND (id IN (SELECT id FROM tenant_runtimes WHERE tenant_id = $2))`),
				Args:     []driver.Value{runtimeID, tenantID},
				IsSelect: true,
				ValidRowsProvider: func() []*sqlmock.Rows {
					return []*sqlmock.Rows{testdb.RowWhenObjectExist()}
				},
				InvalidRowsProvider: func() []*sqlmock.Rows {
					return []*sqlmock.Rows{testdb.RowWhenObjectDoesNotExist()}
				},
			},
		},
		ConverterMockProvider: func() testdb.Mock {
			return &automock.EntityConverter{}
		},
		RepoConstructorFunc: runtime.NewRepository,
		TargetID:            runtimeID,
		TenantID:            tenantID,
	}

	suite.Run(t)
}
