package clients

import (
	"api-gateway/internal/config"
	"api-gateway/internal/logger"

	"api-gateway/internal/proto/separator"
	"api-gateway/internal/proto/tab"

	"github.com/supabase-community/supabase-go"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	Supabase *supabase.Client

	TabGenClient TabGenerator
	//	audioProcessorClient audioproc.AudioProcessorServiceClient
	AudioSeparatorClient AudioSeparator

	tabGenConn         *grpc.ClientConn
	audioProcessorConn *grpc.ClientConn
	audioSeparatorConn *grpc.ClientConn
)

func InitClients() {
	cfg := config.Load()
	var err error

	Supabase, err = supabase.NewClient(cfg.SupabaseURL, cfg.SupabaseKey, &supabase.ClientOptions{})
	if err != nil {
		logger.Log.Fatal("supabase connection failed", zap.Error(err))
	}

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

	TabGenClient = tab.NewTabGenerateClient(tabGenConn)
	//	audioProcessorClient = audioproc.NewAudioProcessorServiceClient(audioProcessorConn)
	AudioSeparatorClient = separator.NewAudioSeparatorClient(audioSeparatorConn)

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
