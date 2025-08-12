package router

import (
	"api-gateway/internal/handlers"

	"github.com/labstack/echo"
)

type Router struct {
	tabHandler   *handlers.TabHandler
	audioHandler *handlers.AudioHandler
}

func NewRouter(tabHandler *handlers.TabHandler, audioHandler *handlers.AudioHandler) *Router {
	return &Router{
		tabHandler:   tabHandler,
		audioHandler: audioHandler,
	}
}

func (r *Router) RegisterRoutes(e *echo.Echo) {
	e.Static("/static", "static")

	e.GET("/", func(c echo.Context) error {
		return c.File("static/index.html")
	})

	tabGroup := e.Group("/tab")

	tabGroup.POST("/generate", r.tabHandler.TabGenerate)
	tabGroup.GET("/search", r.tabHandler.SearchTabs)
	tabGroup.POST("/", r.tabHandler.CreateTab)
	tabGroup.DELETE("/:id", r.tabHandler.DeleteTab)
	tabGroup.GET("/:id", r.tabHandler.GetTab)
	tabGroup.GET("/view/:id", r.tabHandler.ViewTabPage)

	audioGroup := e.Group("/audio")

	audioGroup.POST("/separate", r.audioHandler.SeparateAudio)
}
