package rest

import (
	"net/url"

	"github.com/dexconv/hasin/retrieve/common"
	"github.com/dexconv/hasin/retrieve/config"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var (
	log       = common.Log
	EP        = echo.New()
	JwtConfig = middleware.JWTConfig{
		SigningMethod: "HS512",
		SigningKey:    []byte(secret),
	}
)

func Run() {

	u := EP.Group("/user")
	u.POST("/register", userRegister)
	u.POST("/login", userLogin)

	e := EP.Group("/api", middleware.JWTWithConfig(JwtConfig))
	url, err := url.Parse("http://store:8080")
	if err != nil {
		EP.Logger.Fatal(err)
	}
	targets := []*middleware.ProxyTarget{
		{
			URL: url,
		},
	}
	e.Use(middleware.Proxy(middleware.NewRoundRobinBalancer(targets)))

	log.Info("starting rest server")
	EP.Logger.Fatal(EP.Start(config.GLB.Port))
}
