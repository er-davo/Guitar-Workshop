package clients

import (
	"api-gateway/internal/config"
	"api-gateway/internal/logger"

	"api-gateway/internal/proto/separator"
	"api-gateway/internal/proto/tab"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	tabGenClient TabGenerator
	//	audioProcessorClient audioproc.AudioProcessorServiceClient
	audioSeparatorClient AudioSeparator

	tabGenConn         *grpc.ClientConn
	audioProcessorConn *grpc.ClientConn
	audioSeparatorConn *grpc.ClientConn
)

func init() {
	cfg := config.Load()
	var err error

	tabGenConn, err = grpc.NewClient(
		cfg.TabgenHost+":"+cfg.TabgenPort,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultCallOptions(
			grpc.MaxCallRecvMsgSize(100*1024*1024), // 100 MB
			grpc.MaxCallSendMsgSize(100*1024*1024),
		),
	)
	if err != nil {
		logger.Log.Fatal("tab-generator gRPC connection failed", zap.Error(err))
	}

	// audioProcessorConn, err = grpc.NewClient(
	// 	config.Load().AudioProcHost+":"+config.Load().AudioProcPort,
	// 	grpc.WithTransportCredentials(insecure.NewCredentials()),
	// 	grpc.WithDefaultCallOptions(
	// 		grpc.MaxCallRecvMsgSize(100*1024*1024), // 100 MB
	// 		grpc.MaxCallSendMsgSize(100*1024*1024),
	// 	),
	// )
	// if err != nil {
	// 	logger.Log.Fatal("audio-processor gRPC connection failed", zap.Error(err))
	// }

	audioSeparatorConn, err = grpc.NewClient(
		cfg.AudioSeparatorHost+":"+cfg.AudioSeparatorPort,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultCallOptions(
			grpc.MaxCallRecvMsgSize(250<<20), // 250 MB
			grpc.MaxCallSendMsgSize(250<<20),
		),
	)
	if err != nil {
		logger.Log.Fatal("audio-separator gRPC connection failed", zap.Error(err))
	}

	tabGenClient = tab.NewTabGenerateClient(tabGenConn)
	//	audioProcessorClient = audioproc.NewAudioProcessorServiceClient(audioProcessorConn)
	audioSeparatorClient = separator.NewAudioSeparatorClient(audioSeparatorConn)

}

// func InitClients() {
// 	var err error
//
// 	tabGenConn, err = grpc.NewClient(
// 		config.Load().TabgenHost+":"+config.Load().TabgenPort,
// 		grpc.WithTransportCredentials(insecure.NewCredentials()),
// 		grpc.WithDefaultCallOptions(
// 			grpc.MaxCallRecvMsgSize(100*1024*1024), // 100 MB
// 			grpc.MaxCallSendMsgSize(100*1024*1024),
// 		),
// 	)
// 	if err != nil {
// 		logger.Log.Fatal("tab-generator gRPC connection failed", zap.Error(err))
// 	}
//
// 	// audioProcessorConn, err = grpc.NewClient(
// 	// 	config.Load().AudioProcHost+":"+config.Load().AudioProcPort,
// 	// 	grpc.WithTransportCredentials(insecure.NewCredentials()),
// 	// 	grpc.WithDefaultCallOptions(
// 	// 		grpc.MaxCallRecvMsgSize(100*1024*1024), // 100 MB
// 	// 		grpc.MaxCallSendMsgSize(100*1024*1024),
// 	// 	),
// 	// )
// 	// if err != nil {
// 	// 	logger.Log.Fatal("audio-processor gRPC connection failed", zap.Error(err))
// 	// }
//
// 	audioSeparatorConn, err = grpc.NewClient(
// 		config.Load().AudioSeparatorHost+":"+config.Load().AudioSeparatorPort,
// 		grpc.WithTransportCredentials(insecure.NewCredentials()),
// 		grpc.WithDefaultCallOptions(
// 			grpc.MaxCallRecvMsgSize(100*1024*1024), // 100 MB
// 			grpc.MaxCallSendMsgSize(100*1024*1024),
// 		),
// 	)
// 	if err != nil {
// 		logger.Log.Fatal("audio-separator gRPC connection failed", zap.Error(err))
// 	}
//
// 	TabGenClient = tab.NewTabGenerateClient(tabGenConn)
// 	AudioProcessorClient = audioproc.NewAudioProcessorServiceClient(audioProcessorConn)
// 	AudioSeparatorClient = separator.NewAudioSeparatorClient(audioSeparatorConn)
// }

func CloseClients() {
	if tabGenConn != nil {
		tabGenConn.Close()
	}

	if audioProcessorConn != nil {
		audioProcessorConn.Close()
	}

	if audioSeparatorConn != nil {
		audioSeparatorConn.Close()
	}
}
