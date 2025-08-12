package main

import (
	"context"
	"fmt"

	"api-gateway/internal/clients"
	"api-gateway/internal/config"
	"api-gateway/internal/database"
	"api-gateway/internal/handlers"
	"api-gateway/internal/logger"
	"api-gateway/internal/repository"
	"api-gateway/internal/router"
	"api-gateway/internal/service"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/supabase-community/supabase-go"
)

func main() {
	cfg := config.Load()

	dbConn := database.Connect(cfg.DatabaseURL)
	defer dbConn.Close(context.Background())

	tabRepo := repository.NewTabRepository(dbConn)

	supabaseClient, err := supabase.NewClient(cfg.SupabaseURL, cfg.SupabaseKey, &supabase.ClientOptions{})
	if err != nil {
		logger.Log.Fatal("error on connecting to supabase: " + err.Error())
	}

	tabGen, err := clients.NewTabGenerator(cfg.TabgenHost + ":" + cfg.TabgenPort)
	if err != nil {
		logger.Log.Fatal("error on connecting to grpc tab generator: " + err.Error())
	}
	defer tabGen.Close()

	audioSeparator, err := clients.NewAudioSeparator(cfg.AudioSeparatorHost + ":" + cfg.AudioSeparatorPort)
	if err != nil {
		logger.Log.Fatal("error on connecting to grpc audio separator: " + err.Error())
	}
	defer audioSeparator.Close()

	tabService := service.NewTabService(tabRepo, supabaseClient, tabGen, audioSeparator)
	audioService := service.NewAudioService(audioSeparator)

	tabHandler := handlers.NewTabHandler(tabService)
	audioHandler := handlers.NewAudioHandler(audioService)

	router := router.NewRouter(tabHandler, audioHandler)

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	router.RegisterRoutes(e)

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", cfg.PORT)))
}
