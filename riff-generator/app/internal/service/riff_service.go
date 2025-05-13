package service

import (
	"context"
	riffpb "riffgen/internal/riffproto"
)

type RiffService struct {
	riffpb.UnimplementedRiffGeneratorServer
}

func (s *RiffService) GenerateRiff(ctx context.Context, req *riffpb.RiffRequest) (*riffpb.RiffResponse, error) {
	return nil, nil
}
