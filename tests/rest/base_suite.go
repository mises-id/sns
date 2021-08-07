package rest

import (
	"net/http"
	"strings"

	"github.com/gavv/httpexpect"
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/mises-id/sns/app/services/session"
	"github.com/mises-id/sns/config/route"
	"github.com/mises-id/sns/tests"
	misesMock "github.com/mises-id/sns/tests/mocks/lib/mises"
)

type RestBaseTestSuite struct {
	tests.BaseTestSuite
	Handler   http.Handler
	Expect    *httpexpect.Expect
	mockMises *misesMock.MockClient
	ctrl      *gomock.Controller
}

func (suite *RestBaseTestSuite) SetupSuite() {
	suite.BaseTestSuite.SetupSuite()
	suite.SetupEchoHandler()
	suite.InitExpect()
	suite.ctrl = gomock.NewController(suite.T())
	suite.mockMises = misesMock.NewMockClient(suite.ctrl)
}

func (suite *RestBaseTestSuite) TearDownSuite() {
	suite.ctrl.Finish()
	suite.BaseTestSuite.TearDownSuite()
}

func (suite *RestBaseTestSuite) SetupEchoHandler() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.DefaultCORSConfig))

	route.SetRoutes(e)
	suite.Handler = e
}

func (suite *RestBaseTestSuite) InitExpect() {
	suite.Expect = httpexpect.WithConfig(httpexpect.Config{
		Client: &http.Client{
			Transport: httpexpect.NewBinder(suite.Handler),
		},
		Reporter: httpexpect.NewRequireReporter(suite.T()),
	})
}

func (suite *RestBaseTestSuite) LoginUser(auth string) string {
	resp := suite.Expect.POST("/api/v1/signin").WithJSON(map[string]interface{}{
		"provider": "mises",
		"user_authz": map[string]interface{}{
			"auth": auth,
		},
	}).Expect().Status(http.StatusOK).JSON().Object()
	return resp.Value("data").Object().Value("token").String().Raw()
}

func (suite *RestBaseTestSuite) MockLoginUser(auth string) string {
	session.MockMisesClient(suite.mockMises)
	suite.mockMises.EXPECT().Auth(auth).Return(strings.Split(auth, ":")[0], nil)
	resp := suite.Expect.POST("/api/v1/signin").WithJSON(map[string]interface{}{
		"provider": "mises",
		"user_authz": map[string]interface{}{
			"auth": auth,
		},
	}).Expect().Status(http.StatusOK).JSON().Object()
	return resp.Value("data").Object().Value("token").String().Raw()
}

func (suite *RestBaseTestSuite) MockMisesAuth(auth string) {
	session.MockMisesClient(suite.mockMises)
	suite.mockMises.EXPECT().Auth(auth).Return(strings.Split(auth, ":")[0], nil)
}
