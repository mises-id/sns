package user

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/mises-id/sns/app/models"
	"github.com/mises-id/sns/app/models/enum"
	"github.com/mises-id/sns/tests/factories"
	"github.com/mises-id/sns/tests/rest"
	"github.com/stretchr/testify/suite"
)

type UserServerSuite struct {
	rest.RestBaseTestSuite
	collections []string
}

func (suite *UserServerSuite) SetupSuite() {
	suite.RestBaseTestSuite.SetupSuite()
	suite.collections = []string{"counters", "attachments", "users"}
}

func (suite *UserServerSuite) TearDownSuite() {
	suite.RestBaseTestSuite.TearDownSuite()
}

func (suite *UserServerSuite) SetupTest() {
	suite.Clean(suite.collections...)
	suite.Acquire(suite.collections...)
}

func (suite *UserServerSuite) TearDownTest() {
	suite.Clean(suite.collections...)
}

func TestUserServer(t *testing.T) {
	suite.Run(t, &UserServerSuite{})
}

func (suite *UserServerSuite) TestFindUser() {
	factories.InitAttachments(&models.Attachment{
		ID:        1,
		Filename:  "test.jpg",
		FileType:  enum.ImageFile,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})
	factories.InitUsers(&models.User{
		UID:      1,
		AvatarID: 1,
		Gender:   enum.GenderMale,
	}, &models.User{
		UID:      2,
		Gender:   enum.GenderMale,
		AvatarID: 0,
	})
	suite.T().Run("not found user", func(t *testing.T) {
		resp := suite.Expect.GET("/api/v1/user/999").
			Expect().Status(http.StatusNotFound).JSON().Object()
		resp.Value("code").Equal(404000)
	})

	suite.T().Run("find user with avatar", func(t *testing.T) {
		resp := suite.Expect.GET("/api/v1/user/1").
			Expect().Status(http.StatusOK).JSON().Object()
		resp.Value("code").Equal(0)
		resp.Value("data").Object().Value("uid").Equal(1)
		url := fmt.Sprintf("http://localhost/upload/attachment/%s/1/test.jpg", time.Now().Format("2006/01/02"))
		resp.Value("data").Object().Value("avatar").Object().Value("small").Equal(url)
	})

	suite.T().Run("find user without avatar", func(t *testing.T) {
		resp := suite.Expect.GET("/api/v1/user/2").
			Expect().Status(http.StatusOK).JSON().Object()
		resp.Value("code").Equal(0)
		resp.Value("data").Object().Value("uid").Equal(2)
		resp.Value("data").Object().Value("avatar").Null()
	})
}

func (suite *UserServerSuite) TestSignin() {
	factories.InitUsers(&models.User{
		UID:      1001,
		AvatarID: 0,
		Misesid:  "123",
		Gender:   enum.GenderMale,
	})
	factories.InitAttachments(&models.Attachment{
		ID:        1,
		Filename:  "test.jpg",
		FileType:  enum.ImageFile,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})
	suite.T().Run("user signin success", func(t *testing.T) {
		resp := suite.Expect.POST("/api/v1/signin").WithJSON(map[string]interface{}{
			"provider": "mises",
			"user_authz": map[string]interface{}{
				"misesid":   "123",
				"auth_code": "123",
			},
		}).Expect().Status(http.StatusOK).JSON().Object()
		resp.Value("code").Equal(0)
	})
	suite.T().Run("create new success", func(t *testing.T) {
		resp := suite.Expect.POST("/api/v1/signin").WithJSON(map[string]interface{}{
			"provider": "mises",
			"user_authz": map[string]interface{}{
				"misesid":   "234",
				"auth_code": "234",
			},
		}).Expect().Status(http.StatusOK).JSON().Object()
		resp.Value("code").Equal(0)
	})
}

func (suite *UserServerSuite) TestUpdateUser() {
	factories.InitUsers(&models.User{
		UID:      1001,
		Gender:   enum.GenderFemale,
		AvatarID: 0,
		Misesid:  "123",
	})
	token := suite.LoginUser("123")
	suite.T().Run("update username success", func(t *testing.T) {
		resp := suite.Expect.PATCH("/api/v1/user/me").WithJSON(map[string]interface{}{
			"by": "username",
			"username": map[string]interface{}{
				"username": "Hahaha",
			},
		}).WithHeader("Authorization", "Bearer "+token).Expect().Status(http.StatusOK).JSON().Object()
		resp.Value("code").Equal(0)
	})
	suite.T().Run("update user avatar success", func(t *testing.T) {
		resp := suite.Expect.PATCH("/api/v1/user/me").WithJSON(map[string]interface{}{
			"by": "avatar",
			"avatar": map[string]interface{}{
				"attachment_id": 1,
			},
		}).WithHeader("Authorization", "Bearer "+token).Expect().Status(http.StatusOK).JSON().Object()
		resp.Value("code").Equal(0)
	})
	suite.T().Run("update user profile success", func(t *testing.T) {
		resp := suite.Expect.PATCH("/api/v1/user/me").WithJSON(map[string]interface{}{
			"by": "profile",
			"profile": map[string]interface{}{
				"email":   "test@t.com",
				"gender":  "female",
				"mobile":  "123456",
				"address": "xxxx",
			},
		}).WithHeader("Authorization", "Bearer "+token).Expect().Status(http.StatusOK).JSON().Object()
		resp.Value("code").Equal(0)
	})
}
