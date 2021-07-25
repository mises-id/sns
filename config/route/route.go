package route

import (
	"github.com/labstack/echo"
	"github.com/mises-id/sns/app/apis/rest"
	v1 "github.com/mises-id/sns/app/apis/rest/v1"
	mw "github.com/mises-id/sns/lib/middleware"
)

// SetRoutes sets the routes of echo http server
func SetRoutes(e *echo.Echo) {
	e.GET("/", rest.Probe)
	e.GET("/healthz", rest.Probe)

	groupV1 := e.Group("/api/v1", mw.ErrorResponseMiddleware)
	groupV1.POST("/attachment", v1.Upload)
	groupV1.GET("/user/:uid", v1.FindUser)
}
