package status

import (
	"context"
	"net/http"
	"testing"

	"github.com/mises-id/sns/app/models"
	"github.com/mises-id/sns/app/models/enum"
	"github.com/mises-id/sns/lib/codes"
	"github.com/mises-id/sns/lib/db"
	"github.com/mises-id/sns/tests/factories"
	"github.com/mises-id/sns/tests/rest"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson"
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
	suite.Clean(suite.collections...)
}

func TestStatusServer(t *testing.T) {
	suite.Run(t, &StatusServerSuite{})
}

func (suite *StatusServerSuite) TestListStatus() {
	token := suite.MockLoginUser("1001:123")
	suite.T().Run("list me status", func(t *testing.T) {
		resp := suite.Expect.GET("/api/v1/status/recommend").
			WithHeader("Authorization", "Bearer "+token).Expect().Status(http.StatusOK).JSON().Object()
		resp.Value("data").Array()
	})
}

func (suite *StatusServerSuite) TestCreateStatus() {
	token := suite.MockLoginUser("1001:123")
	linkMeta := &map[string]interface{}{
		"title":         "Test link title",
		"host":          "www.test.com",
		"attachment_id": uint64(1),
		"link":          "http://www.test.com/articles/test/1",
	}
	suite.T().Run("create a text status", func(t *testing.T) {
		resp := suite.Expect.POST("/api/v1/status").WithJSON(map[string]interface{}{
			"status_type": "text",
			"content":     "post a text status",
		}).WithHeader("Authorization", "Bearer "+token).Expect().Status(http.StatusOK).JSON().Object()
		resp.Value("code").Equal(codes.SuccessCode)
		status := &models.Status{}
		err := db.ODM(context.Background()).Last(status).Error
		suite.Nil(err)
		suite.Equal("post a text status", status.Content)
		suite.Equal(enum.TextStatus, status.StatusType)
		suite.Equal(uint64(1001), status.UID)
	})
	suite.T().Run("create a link status", func(t *testing.T) {
		resp := suite.Expect.POST("/api/v1/status").WithJSON(map[string]interface{}{
			"status_type": "link",
			"content":     "post a link status",
			"meta":        linkMeta,
		}).WithHeader("Authorization", "Bearer "+token).Expect().Status(http.StatusOK).JSON().Object()
		resp.Value("code").Equal(codes.SuccessCode)
		status := &models.Status{}
		err := db.ODM(context.Background()).Last(status).Error
		suite.Nil(err)
		suite.Equal("post a link status", status.Content)
		suite.Equal(enum.LinkStatus, status.StatusType)
		suite.Equal(uint64(1001), status.UID)
	})
	suite.T().Run("forward a text status", func(t *testing.T) {
		resp := suite.Expect.POST("/api/v1/status").WithJSON(map[string]interface{}{
			"status_type":      "text",
			"parent_status_id": suite.statuses[0].ID.Hex(),
			"origin_status_id": suite.statuses[0].ID.Hex(),
			"content":          "forward a text status",
		}).WithHeader("Authorization", "Bearer "+token).Expect().Status(http.StatusOK).JSON().Object()
		resp.Value("code").Equal(codes.SuccessCode)
		status := &models.Status{}
		err := db.ODM(context.Background()).Last(status).Error
		suite.Nil(err)
		suite.Equal("forward a text status", status.Content)
		suite.Equal(enum.TextStatus, status.StatusType)
		suite.Equal(suite.statuses[0].ID.Hex(), status.ParentID.Hex())
		suite.Equal(suite.statuses[0].ID.Hex(), status.OriginID.Hex())
		suite.Equal(uint64(1001), status.UID)

		parentStatus := &models.Status{}
		err = db.ODM(context.Background()).First(parentStatus, bson.M{"_id": suite.statuses[0].ID}).Error
		suite.Nil(err)
		suite.Equal(uint64(1), parentStatus.ForwardsCount)
	})
	suite.T().Run("forward a link status", func(t *testing.T) {
		resp := suite.Expect.POST("/api/v1/status").WithJSON(map[string]interface{}{
			"status_type":      "text",
			"parent_status_id": suite.statuses[1].ID.Hex(),
			"origin_status_id": suite.statuses[1].ID.Hex(),
			"content":          "forward a link status",
		}).WithHeader("Authorization", "Bearer "+token).Expect().Status(http.StatusOK).JSON().Object()
		resp.Value("code").Equal(codes.SuccessCode)
		status := &models.Status{}
		err := db.ODM(context.Background()).Last(status).Error
		suite.Nil(err)
		suite.Equal("forward a link status", status.Content)
		suite.Equal(enum.TextStatus, status.StatusType)
		suite.Equal(suite.statuses[1].ID.Hex(), status.ParentID.Hex())
		suite.Equal(suite.statuses[1].ID.Hex(), status.OriginID.Hex())
		suite.Equal(uint64(1001), status.UID)

		parentStatus := &models.Status{}
		err = db.ODM(context.Background()).First(parentStatus, bson.M{"_id": suite.statuses[1].ID}).Error
		suite.Nil(err)
		suite.Equal(uint64(1), parentStatus.ForwardsCount)
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
