package main

import (
	"fmt"
	"net"

	"tabgen/internal/config"
	"tabgen/internal/logger"
	"tabgen/internal/service"
	tabpb "tabgen/internal/tabproto"

	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func main() {
	defer logger.Log.Sync()

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", config.Load().PORT))
	if err != nil {
		logger.Log.Fatal("Failed to listen", zap.Error(err))
	}

	s := grpc.NewServer()
	tabpb.RegisterTabGenerateServer(s, &service.TabService{})

	logger.Log.Info("gRPC server running on " + fmt.Sprintf(":%s", config.Load().PORT))
	if err := s.Serve(lis); err != nil {
		logger.Log.Fatal("Failed to serve", zap.Error(err))
	}
}
