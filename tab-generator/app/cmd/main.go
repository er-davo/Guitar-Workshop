package main

import (
	"net"
	"tabgen/internal/logger"
	"tabgen/internal/service"
	tabpb "tabgen/internal/tabproto"

	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func main() {
	defer logger.Log.Sync()

	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		logger.Log.Fatal("Failed to listen", zap.Error(err))
	}

	s := grpc.NewServer()
	tabpb.RegisterTabGenerateServer(s, &service.TabService{})

	logger.Log.Info("gRPC server running on :8080")
	if err := s.Serve(lis); err != nil {
		logger.Log.Fatal("Failed to serve", zap.Error(err))
	}

}
