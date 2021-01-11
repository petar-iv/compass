package healthcheck_test

import (
	"credentials-broker/tests/common"
	"net/http"
	"testing"

	"github.com/stretchr/testify/suite"
)

func TestHealthcheckTestSuite(t *testing.T) {
	suite.Run(t, new(HealthcheckTestSuite))
}

type HealthcheckTestSuite struct {
	suite.Suite

	testContext *common.TestContext
}

func (suite *HealthcheckTestSuite) SetupSuite() {
	suite.testContext = common.NewTestContextBuilder().Build(suite.T())
}

func (suite *HealthcheckTestSuite) TearDownSuite() {
	suite.testContext.CleanUp()
}

func (suite *HealthcheckTestSuite) TestHealthcheck() {
	suite.testContext.SystemBroker.GET("/healthz").Expect().
		Status(http.StatusOK).
		Body().Equal(`{"status": "success"}`)
}
