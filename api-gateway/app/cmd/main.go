package main

import (
	"context"
	"fmt"

	"api-gateway/internal/clients"
	"api-gateway/internal/config"
	"api-gateway/internal/database"
	"api-gateway/internal/handlers"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func main() {
	clients.InitClients()
	defer clients.CloseClients()

	cfg := config.Load()

	dbConn := database.Connect(cfg.DatabaseURL)
	defer dbConn.Close(context.Background())
	handlers.Init(dbConn)

	e := echo.New()

	e.Static("/static", "static")

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/", func(c echo.Context) error {
		return c.File("static/index.html")
	})

	tabGroup := e.Group("/tab")

	tabGroup.POST("/generate", handlers.TabGenerate)
	tabGroup.POST("/save", handlers.SaveTab)
	tabGroup.GET("/search", handlers.SearchTabs)
	tabGroup.GET("/:id", handlers.GetTab)
	tabGroup.GET("/view/:id", handlers.ViewTabPage)

	e.POST("/audio/separate", handlers.SeparateAudio)

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", cfg.PORT)))
}
