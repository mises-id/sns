package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/mises-id/sns/config/env"
	"github.com/mises-id/sns/config/route"
	"github.com/mises-id/sns/lib/db"
	_ "github.com/mises-id/sns/lib/mises"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*60)
	defer cancel()
	db.SetupMongo(ctx)
	e := echo.New()
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Skipper: func(c echo.Context) bool { return c.Path() == "/" },
		Format: `{"timestamp":"${time_rfc3339}","serviceContext":{"service":"mises-sns"},"message":"${remote_ip} ${status} ${method} ${uri}",` +
			`"severity":"INFO","context":{"request_id":"${id}","remote_ip":"${remote_ip}","host":"${host}","method":"${method}","uri":"${uri}",` +
			`"user_agent":"${user_agent}","status":"${status}","error":"${error}","latency_human":"${latency_human}","device_id":"${header:x-device-id}"}}` + "\n",
	}))
	e.Use(middleware.Recover())
	e.Use(middleware.RequestID())
	route.SetRoutes(e)
	go func() {
		if err := e.Start(fmt.Sprintf(":%d", env.Envs.Port)); err != nil {
			log.Fatal(err)
		}
	}()
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	_ = e.Shutdown(ctx)
}
