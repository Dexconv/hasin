package rest

import (
	"net/http"

	"github.com/dexconv/hasin/store/common"
	"github.com/dexconv/hasin/store/config"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var (
	log = common.Log
	EP  = echo.New()
)

func Run() {
	e := EP.Group("/api")

	e.GET("/ping", func(c echo.Context) error {
		return c.String(http.StatusOK, "pong")
	})

	e.POST("/upload", upload, middleware.BodyLimit("1M"))
	e.GET("/download/:mode/:tags", getWithTags)

	log.Info("starting rest server")
	EP.Logger.Fatal(EP.Start(config.GLB.Port))
}
