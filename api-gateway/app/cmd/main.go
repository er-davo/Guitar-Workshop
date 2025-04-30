package main

import (
	"api-gateway/internal/logger"
	tabpb "api-gateway/internal/tabproto"
	"context"
	"net/http"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"

	"go.uber.org/zap"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	e := echo.New()
	conn, err := grpc.NewClient(
		"tab-generate:8081",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		logger.Log.Fatal("gRPC connection failed", zap.Error(err))
	}
	defer conn.Close()

	client := tabpb.NewTabGenerateClient(conn)

	e.Static("/static", "sctatic")

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/", func(c echo.Context) error {
		return c.File("index.html")
	})

	// TODO
	e.POST("tab-generate", func(c echo.Context) error {
		var tabReq tabpb.TabRequest
		if err := c.Bind(&tabReq); err != nil {
			return c.JSON(http.StatusBadRequest, "Invalid request")
		}

		tabResp, err := client.GenerateTab(context.Background(), &tabReq)
		if err != nil {
			logger.Log.Error("error on process audio", zap.Error(err))
			return c.JSON(http.StatusInternalServerError, "Tab generation failed")
		}

		return c.JSON(http.StatusOK, tabResp)
	})

	e.Logger.Fatal(e.Start("0.0.0.0:8080"))

}
