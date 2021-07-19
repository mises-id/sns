package route

import (
	"github.com/labstack/echo"
	"github.com/mises-id/sns/app/apis/rest"
)

// SetRoutes sets the routes of echo http server
func SetRoutes(e *echo.Echo) {
	e.GET("/", rest.Probe)
	e.GET("/healthz", rest.Probe)

}
