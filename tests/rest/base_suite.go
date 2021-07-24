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
