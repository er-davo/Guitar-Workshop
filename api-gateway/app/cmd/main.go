package main

import (
	"api-gateway/internal/handlers"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func main() {
	e := echo.New()

	e.Static("/static", "sctatic")

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/", func(c echo.Context) error {
		return c.File("static/index.html")
	})

	e.POST("/generate-tab", handlers.TabGenerate)

	e.Logger.Fatal(e.Start("0.0.0.0:8080"))
}
