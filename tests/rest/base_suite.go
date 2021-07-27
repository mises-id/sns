package rest

import (
	"net/http"

	"github.com/gavv/httpexpect"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/mises-id/sns/config/route"
	"github.com/mises-id/sns/tests"
)

type RestBaseTestSuite struct {
	tests.BaseTestSuite
	Handler http.Handler
	Expect  *httpexpect.Expect
}

func (s *RestBaseTestSuite) SetupSuite() {
	s.BaseTestSuite.SetupSuite()
	s.SetupEchoHandler()
	s.InitExpect()
}

func (s *RestBaseTestSuite) SetupEchoHandler() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.DefaultCORSConfig))

	route.SetRoutes(e)
	s.Handler = e
}

func (s *RestBaseTestSuite) InitExpect() {
	s.Expect = httpexpect.WithConfig(httpexpect.Config{
		Client: &http.Client{
			Transport: httpexpect.NewBinder(s.Handler),
		},
		Reporter: httpexpect.NewRequireReporter(s.T()),
	})
}

func (suite *RestBaseTestSuite) LoginUser(misesid string) string {
	resp := suite.Expect.POST("/api/v1/signin").WithJSON(map[string]interface{}{
		"provider": "mises",
		"user_authz": map[string]interface{}{
			"misesid":   misesid,
			"auth_code": "123",
		},
	}).Expect().Status(http.StatusOK).JSON().Object()
	return resp.Value("data").Object().Value("token").String().Raw()
}
