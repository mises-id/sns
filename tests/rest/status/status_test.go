package status

import (
	"testing"

	"github.com/mises-id/sns/tests/rest"
	"github.com/stretchr/testify/suite"
)

type StatusServerSuite struct {
	rest.RestBaseTestSuite
	collections []string
}

func (suite *StatusServerSuite) SetupSuite() {
	suite.RestBaseTestSuite.SetupSuite()
	suite.collections = []string{"counters", "users", "follows", "statuses"}
}

func (suite *StatusServerSuite) TearDownSuite() {
	suite.RestBaseTestSuite.TearDownSuite()
}

func (suite *StatusServerSuite) SetupTest() {
	suite.Clean(suite.collections...)
	suite.Acquire(suite.collections...)
}

func (suite *StatusServerSuite) TearDownTest() {
	suite.Clean(suite.collections...)
}

func TestStatusServer(t *testing.T) {
	suite.Run(t, &StatusServerSuite{})
}

func (suite *StatusServerSuite) TestListStatus() {
	// user1 := factories.UserFactory.MustCreate().(*models.User)

}

func (suite *StatusServerSuite) TestCreateStatus() {
	// user1 := factories.UserFactory.MustCreate().(*models.User)

}

func (suite *StatusServerSuite) TestDeleteStatus() {
	// user1 := factories.UserFactory.MustCreate().(*models.User)

}
