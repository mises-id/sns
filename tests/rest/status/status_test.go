package status

import (
	"net/http"
	"testing"

	"github.com/mises-id/sns/app/models"
	"github.com/mises-id/sns/lib/codes"
	"github.com/mises-id/sns/tests/factories"
	"github.com/mises-id/sns/tests/rest"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type StatusServerSuite struct {
	rest.RestBaseTestSuite
	collections []string
	statuses    []*models.Status
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
	factories.InitUsers(&models.User{
		UID:     uint64(1001),
		Misesid: "1001",
	}, &models.User{
		UID:     uint64(1002),
		Misesid: "1002",
	})
	suite.statuses = factories.InitDefaultStatuses()
}

func (suite *StatusServerSuite) TearDownTest() {
	// suite.Clean(suite.collections...)
}

func TestStatusServer(t *testing.T) {
	suite.Run(t, &StatusServerSuite{})
}

func (suite *StatusServerSuite) TestListStatus() {
}

func (suite *StatusServerSuite) TestCreateStatus() {
	token := suite.MockLoginUser("1001:123")
	suite.T().Run("create a text status", func(t *testing.T) {
		resp := suite.Expect.POST("/api/v1/status").WithJSON(map[string]interface{}{
			"status_type": "text",
			"content":     "post a text status",
		}).WithHeader("Authorization", "Bearer "+token).Expect().Status(http.StatusOK).JSON().Object()
		resp.Value("code").Equal(codes.SuccessCode)
	})
	suite.T().Run("create a link status", func(t *testing.T) {
	})
	suite.T().Run("forward a text status", func(t *testing.T) {
	})
	suite.T().Run("forward a link status", func(t *testing.T) {
	})
}

func (suite *StatusServerSuite) TestDeleteStatus() {
	token := suite.MockLoginUser("1001:123")
	suite.T().Run("delete status not found", func(t *testing.T) {
		resp := suite.Expect.DELETE("/api/v1/status/xxxxxxx").
			WithHeader("Authorization", "Bearer "+token).Expect().Status(http.StatusBadRequest).JSON().Object()
		resp.Value("code").Equal(codes.InvalidArgumentCode)

		resp = suite.Expect.DELETE("/api/v1/status/"+primitive.NewObjectID().Hex()).
			WithHeader("Authorization", "Bearer "+token).Expect().Status(http.StatusNotFound).JSON().Object()
		resp.Value("code").Equal(codes.NotFoundCode)
	})

	suite.T().Run("delete status forbidden", func(t *testing.T) {
		resp := suite.Expect.DELETE("/api/v1/status/"+suite.statuses[1].ID.Hex()).
			WithHeader("Authorization", "Bearer "+token).Expect().Status(http.StatusForbidden).JSON().Object()
		resp.Value("code").Equal(codes.ForbiddenCode)
	})

	suite.T().Run("delete status success", func(t *testing.T) {
		resp := suite.Expect.DELETE("/api/v1/status/"+suite.statuses[0].ID.Hex()).
			WithHeader("Authorization", "Bearer "+token).Expect().Status(http.StatusOK).JSON().Object()
		resp.Value("code").Equal(codes.SuccessCode)
	})
}
