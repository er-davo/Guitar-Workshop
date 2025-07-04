package main

import (
	"fmt"
	"net"

	"tabgen/internal/clients"
	"tabgen/internal/config"
	"tabgen/internal/logger"
	"tabgen/internal/proto/tab"
	"tabgen/internal/service"

	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func main() {
	clients.InitClients()
	defer clients.CloseClients()
	defer logger.Log.Sync()

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", config.Load().PORT))
	if err != nil {
		logger.Log.Fatal("Failed to listen", zap.Error(err))
	}

	s := grpc.NewServer(
		grpc.MaxRecvMsgSize(100*1024*1024), // 100 MB
		grpc.MaxSendMsgSize(100*1024*1024), // 100 MB
	)
	tab.RegisterTabGenerateServer(s, &service.TabService{})

	logger.Log.Info("gRPC server running on " + fmt.Sprintf(":%s", config.Load().PORT))
	if err := s.Serve(lis); err != nil {
		logger.Log.Fatal("Failed to serve", zap.Error(err))
	}
}
