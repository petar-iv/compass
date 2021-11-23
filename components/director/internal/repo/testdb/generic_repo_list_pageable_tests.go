package testdb

import (
	"context"
	"fmt"
	"testing"

	"github.com/kyma-incubator/compass/components/director/pkg/apperrors"
	"github.com/kyma-incubator/compass/components/director/pkg/persistence"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
)

type PageDetails struct {
	ExpectedModelEntities []interface{}
	ExpectedDBEntities    []interface{}
	ExpectedPage          interface{}
}

type RepoListPageableTestSuite struct {
	Name                      string
	SqlQueryDetails           []SqlQueryDetails
	ConverterMockProvider     func() Mock
	RepoConstructorFunc       interface{}
	AdditionalConverterArgs   []interface{}
	Pages                     []PageDetails
	MethodArgs                []interface{}
	DisableConverterErrorTest bool
	MethodName                string
}

func (suite *RepoListPageableTestSuite) Run(t *testing.T) bool {
	if len(suite.MethodName) == 0 {
		panic("missing method name")
	}

	for _, queryDetails := range suite.SqlQueryDetails {
		if !queryDetails.IsSelect {
			panic("list pageable suite should expect only select SQL statements")
		}
	}

	for _, page := range suite.Pages {
		if len(page.ExpectedDBEntities) != len(page.ExpectedModelEntities) {
			panic("for each DB entity a corresponding model Entity is expected")
		}
	}

	return t.Run(suite.Name, func(t *testing.T) {
		testErr := errors.New("test error")

		t.Run("success", func(t *testing.T) {
			sqlxDB, sqlMock := MockDatabase(t)
			ctx := persistence.SaveToContext(context.TODO(), sqlxDB)

			for _, sqlDetails := range suite.SqlQueryDetails {
				sqlMock.ExpectQuery(sqlDetails.Query).WithArgs(sqlDetails.Args...).WillReturnRows(sqlDetails.ValidRowsProvider()...)
			}

			convMock := suite.ConverterMockProvider()
			for _, page := range suite.Pages {
				for i := range page.ExpectedDBEntities {
					convMock.On("FromEntity", append([]interface{}{page.ExpectedDBEntities[i]}, suite.AdditionalConverterArgs...)...).Return(page.ExpectedModelEntities[i], nil).Once()
				}
			}

			pgRepository := createRepo(suite.RepoConstructorFunc, convMock)
			//WHEN
			res, err := callList(pgRepository, ctx, suite.MethodName, suite.MethodArgs)
			//THEN
			require.NoError(t, err)
			if len(suite.Pages) > 1 { // entity uses dataloaders and load more than one page on a single call.
				require.Len(t, res, len(suite.Pages))
				for _, page := range suite.Pages {
					require.Contains(t, res, page.ExpectedPage)
				}
			} else {
				require.Equal(t, suite.Pages[0].ExpectedPage, res)
			}

			sqlMock.AssertExpectations(t)
			convMock.AssertExpectations(t)
		})

		for i := range suite.SqlQueryDetails {
			t.Run(fmt.Sprintf("error if SQL query %d fail", i), func(t *testing.T) {
				sqlxDB, sqlMock := MockDatabase(t)
				ctx := persistence.SaveToContext(context.TODO(), sqlxDB)

				for _, sqlDetails := range suite.SqlQueryDetails {
					if sqlDetails.Query == suite.SqlQueryDetails[i].Query {
						sqlMock.ExpectQuery(sqlDetails.Query).WithArgs(sqlDetails.Args...).WillReturnError(testErr)
						break
					} else {
						sqlMock.ExpectQuery(sqlDetails.Query).WithArgs(sqlDetails.Args...).WillReturnRows(sqlDetails.ValidRowsProvider()...)
					}
				}

				convMock := suite.ConverterMockProvider()
				pgRepository := createRepo(suite.RepoConstructorFunc, convMock)

				//WHEN
				res, err := callList(pgRepository, ctx, suite.MethodName, suite.MethodArgs)

				//THEN
				require.Nil(t, res)

				require.Error(t, err)
				require.Equal(t, apperrors.InternalError, apperrors.ErrorCode(err))
				require.Contains(t, err.Error(), "Internal Server Error: Unexpected error while executing SQL query")

				sqlMock.AssertExpectations(t)
				convMock.AssertExpectations(t)
			})
		}

		if !suite.DisableConverterErrorTest {
			t.Run("error when conversion fail", func(t *testing.T) {
				sqlxDB, sqlMock := MockDatabase(t)
				ctx := persistence.SaveToContext(context.TODO(), sqlxDB)

				for _, sqlDetails := range suite.SqlQueryDetails {
					sqlMock.ExpectQuery(sqlDetails.Query).WithArgs(sqlDetails.Args...).WillReturnRows(sqlDetails.ValidRowsProvider()...)
				}

				convMock := suite.ConverterMockProvider()
				for _, page := range suite.Pages {
					if len(page.ExpectedDBEntities) > 0 {
						convMock.On("FromEntity", append([]interface{}{page.ExpectedDBEntities[0]}, suite.AdditionalConverterArgs...)...).Return(nil, testErr).Once()
						break
					}
				}
				pgRepository := createRepo(suite.RepoConstructorFunc, convMock)
				//WHEN
				res, err := callList(pgRepository, ctx, suite.MethodName, suite.MethodArgs)
				//THEN
				require.Nil(t, res)

				require.Error(t, err)
				require.Contains(t, err.Error(), testErr.Error())

				sqlMock.AssertExpectations(t)
				convMock.AssertExpectations(t)
			})
		}
	})
}