package main

import (
	"fmt"

	"api-gateway/internal/clients"
	"api-gateway/internal/config"
	"api-gateway/internal/handlers"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func main() {
	clients.InitClients()
	defer clients.CloseClients()

	e := echo.New()

	e.Static("/static", "static")

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/", func(c echo.Context) error {
		return c.File("static/index.html")
	})

	e.POST("/tab/generate", handlers.TabGenerate)

	e.POST("/tab/save", handlers.SaveTab)

	e.POST("/separate-audio", handlers.SeparateAudio)

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", config.Load().PORT)))
}
