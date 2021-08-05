package rest

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/mises-id/sns/config/env"
	"github.com/mises-id/sns/config/route"
)

func Start(ctx context.Context) error {
	e := echo.New()
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Skipper: func(c echo.Context) bool { return c.Path() == "/" },
		Format: `{"timestamp":"${time_rfc3339}","serviceContext":{"service":"mises-sns"},"message":"${remote_ip} ${status} ${method} ${uri}",` +
			`"severity":"INFO","context":{"request_id":"${id}","remote_ip":"${remote_ip}","host":"${host}","method":"${method}","uri":"${uri}",` +
			`"user_agent":"${user_agent}","status":"${status}","error":"${error}","latency_human":"${latency_human}","device_id":"${header:x-device-id}"}}` + "\n",
	}))
	e.Use(middleware.Recover())
	e.Use(middleware.RequestID())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: strings.Split(env.Envs.AllowOrigins, ","),
		AllowMethods: []string{http.MethodGet, http.MethodHead, http.MethodDelete, http.MethodOptions, http.MethodPut, http.MethodPost, http.MethodDelete, http.MethodPatch},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderXRequestedWith, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
	}))
	route.SetRoutes(e)
	go func() {
		if err := e.Start(fmt.Sprintf(":%d", env.Envs.Port)); err != nil {
			log.Fatal(err)
		}
	}()
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	return e.Shutdown(ctx)
}
