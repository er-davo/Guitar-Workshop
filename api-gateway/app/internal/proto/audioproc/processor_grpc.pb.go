// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v6.30.2
// source: processor.proto

package audioproc

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	AudioProcessorService_ProcessAudio_FullMethodName    = "/audioproc.AudioProcessorService/ProcessAudio"
	AudioProcessorService_SplitIntoChunks_FullMethodName = "/audioproc.AudioProcessorService/SplitIntoChunks"
)

// AudioProcessorServiceClient is the client API for AudioProcessorService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type AudioProcessorServiceClient interface {
	ProcessAudio(ctx context.Context, in *ProcessAudioRequest, opts ...grpc.CallOption) (*ProcessAudioResponse, error)
	SplitIntoChunks(ctx context.Context, in *SplitAudioRequest, opts ...grpc.CallOption) (*SplitAudioResponse, error)
}

type audioProcessorServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewAudioProcessorServiceClient(cc grpc.ClientConnInterface) AudioProcessorServiceClient {
	return &audioProcessorServiceClient{cc}
}

func (c *audioProcessorServiceClient) ProcessAudio(ctx context.Context, in *ProcessAudioRequest, opts ...grpc.CallOption) (*ProcessAudioResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ProcessAudioResponse)
	err := c.cc.Invoke(ctx, AudioProcessorService_ProcessAudio_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *audioProcessorServiceClient) SplitIntoChunks(ctx context.Context, in *SplitAudioRequest, opts ...grpc.CallOption) (*SplitAudioResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(SplitAudioResponse)
	err := c.cc.Invoke(ctx, AudioProcessorService_SplitIntoChunks_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// AudioProcessorServiceServer is the server API for AudioProcessorService service.
// All implementations must embed UnimplementedAudioProcessorServiceServer
// for forward compatibility.
type AudioProcessorServiceServer interface {
	ProcessAudio(context.Context, *ProcessAudioRequest) (*ProcessAudioResponse, error)
	SplitIntoChunks(context.Context, *SplitAudioRequest) (*SplitAudioResponse, error)
	mustEmbedUnimplementedAudioProcessorServiceServer()
}

// UnimplementedAudioProcessorServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedAudioProcessorServiceServer struct{}

func (UnimplementedAudioProcessorServiceServer) ProcessAudio(context.Context, *ProcessAudioRequest) (*ProcessAudioResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ProcessAudio not implemented")
}
func (UnimplementedAudioProcessorServiceServer) SplitIntoChunks(context.Context, *SplitAudioRequest) (*SplitAudioResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SplitIntoChunks not implemented")
}
func (UnimplementedAudioProcessorServiceServer) mustEmbedUnimplementedAudioProcessorServiceServer() {}
func (UnimplementedAudioProcessorServiceServer) testEmbeddedByValue()                               {}

// UnsafeAudioProcessorServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to AudioProcessorServiceServer will
// result in compilation errors.
type UnsafeAudioProcessorServiceServer interface {
	mustEmbedUnimplementedAudioProcessorServiceServer()
}

func RegisterAudioProcessorServiceServer(s grpc.ServiceRegistrar, srv AudioProcessorServiceServer) {
	// If the following call pancis, it indicates UnimplementedAudioProcessorServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&AudioProcessorService_ServiceDesc, srv)
}

func _AudioProcessorService_ProcessAudio_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ProcessAudioRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AudioProcessorServiceServer).ProcessAudio(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: AudioProcessorService_ProcessAudio_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AudioProcessorServiceServer).ProcessAudio(ctx, req.(*ProcessAudioRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AudioProcessorService_SplitIntoChunks_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SplitAudioRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AudioProcessorServiceServer).SplitIntoChunks(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: AudioProcessorService_SplitIntoChunks_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AudioProcessorServiceServer).SplitIntoChunks(ctx, req.(*SplitAudioRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// AudioProcessorService_ServiceDesc is the grpc.ServiceDesc for AudioProcessorService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var AudioProcessorService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "audioproc.AudioProcessorService",
	HandlerType: (*AudioProcessorServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ProcessAudio",
			Handler:    _AudioProcessorService_ProcessAudio_Handler,
		},
		{
			MethodName: "SplitIntoChunks",
			Handler:    _AudioProcessorService_SplitIntoChunks_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "processor.proto",
}
