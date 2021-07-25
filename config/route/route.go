package route

import (
	"github.com/labstack/echo"
	"github.com/mises-id/sns/app/apis/rest"
	v1 "github.com/mises-id/sns/app/apis/rest/v1"
	appmw "github.com/mises-id/sns/app/middleware"
	mw "github.com/mises-id/sns/lib/middleware"
)

// SetRoutes sets the routes of echo http server
func SetRoutes(e *echo.Echo) {
	e.GET("/", rest.Probe)
	e.GET("/healthz", rest.Probe)

	groupV1 := e.Group("/api/v1", mw.ErrorResponseMiddleware, appmw.SetCurrentUserMiddleware)
	groupV1.POST("/attachment", v1.Upload)
	groupV1.GET("/user/:uid", v1.FindUser)
	groupV1.POST("/signin", v1.SignIn)

	userGroup := e.Group("/api/v1", mw.ErrorResponseMiddleware, appmw.SetCurrentUserMiddleware, appmw.RequireCurrentUserMiddleware)
	userGroup.GET("/user/me", v1.MyProfile)
	userGroup.PATCH("/user/me", v1.UpdateUser)
}
