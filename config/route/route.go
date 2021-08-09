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
	groupV1.GET("/user/:uid/friendship", v1.ListFriendship)

	userGroup := e.Group("/api/v1", mw.ErrorResponseMiddleware, appmw.SetCurrentUserMiddleware, appmw.RequireCurrentUserMiddleware)
	userGroup.GET("/user/me", v1.MyProfile)
	userGroup.PATCH("/user/me", v1.UpdateUser)
	userGroup.POST("/user/follow", v1.Follow)
	userGroup.DELETE("/user/follow", v1.Unfollow)
	userGroup.POST("/status", v1.CreateStatus)
	userGroup.DELETE("/status/:id", v1.DeleteStatus)
	userGroup.POST("/status/:id/like", v1.LikeStatus)
	userGroup.DELETE("/status/:id/like", v1.UnlikeStatus)

	groupV1.GET("/comment", v1.ListComment)
	userGroup.POST("/comment", v1.CreateComment)
}
