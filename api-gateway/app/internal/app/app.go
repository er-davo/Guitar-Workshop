package app

import (
	"api-gateway/internal/clients"
	"api-gateway/internal/config"
	"api-gateway/internal/database"
	"api-gateway/internal/handlers"
	"api-gateway/internal/repository"
	"api-gateway/internal/service"
	"context"
	"fmt"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/supabase-community/supabase-go"
	"go.uber.org/zap"
)

type App struct {
	log *zap.Logger
	cfg *config.Config

	db *pgxpool.Pool

	tabGen   clients.TabGenerator
	audioSep clients.AudioSeparator

	server *echo.Echo
}

func New(ctx context.Context, cfg *config.Config, log *zap.Logger) (*App, error) {
	if cfg == nil {
		return nil, fmt.Errorf("config is nil")
	}
	if log == nil {
		return nil, fmt.Errorf("logger is nil")
	}

	supabaseClient, err := supabase.NewClient(cfg.Supabase.URL, cfg.Supabase.Key, &supabase.ClientOptions{})
	if err != nil {
		return nil, fmt.Errorf("error on initializing supabase client: %w", err)
	}

	tabGen, err := clients.NewTabGenerator(cfg.Tabgen.Host+":"+cfg.Tabgen.Port, log)
	if err != nil {
		return nil, fmt.Errorf("error on initializing tab generator client: %w", err)
	}

	audioSep, err := clients.NewAudioSeparator(cfg.AudioSep.Host+":"+cfg.AudioSep.Port, log)
	if err != nil {
		tabGen.Close()
		return nil, fmt.Errorf("error on initializing audio separator client: %w", err)
	}

	dbConn, err := database.Connect(ctx, cfg.DatabaseURL)
	if err != nil {
		tabGen.Close()
		audioSep.Close()
		return nil, fmt.Errorf("error on initializing database connection: %w", err)
	}

	repo := repository.NewTabRepository(dbConn)

	tabService := service.NewTabService(repo, supabaseClient, tabGen, audioSep)
	audioService := service.NewAudioService(audioSep)

	tabHandler := handlers.NewTabHandler(tabService, log)
	audioHandler := handlers.NewAudioHandler(audioService, log)

	e := echo.New()

	e.Static("/static", "static")

	e.GET("/", func(c echo.Context) error {
		return c.File("static/index.html")
	})

	e.Use(ZapLogger(log))
	e.Use(middleware.Recover())

	tabHandler.RegisterRoutes(e)
	audioHandler.RegisterRoutes(e)

	return &App{
		log:      log,
		cfg:      cfg,
		tabGen:   tabGen,
		audioSep: audioSep,
		db:       dbConn,
		server:   e,
	}, nil
}

func (a *App) Run(ctx context.Context) error {
	go func() {
		if err := a.server.Start(":" + a.cfg.App.Port); err != nil && err != http.ErrServerClosed {
			a.log.Error("server start failed", zap.Error(err))
		}
	}()

	<-ctx.Done()
	return a.Shutdown()
}

func (a *App) Shutdown() error {
	ctxTimeout, cancelTimeout := context.WithTimeout(context.Background(), a.cfg.App.ShutdownTimeout)
	defer cancelTimeout()
	if err := a.server.Shutdown(ctxTimeout); err != nil {
		return fmt.Errorf("failed to shutdown echo server: %w", err)
	}

	a.db.Close()

	a.tabGen.Close()
	a.audioSep.Close()

	return nil
}
