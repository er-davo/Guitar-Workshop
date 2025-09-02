package handler

import (
	"context"
	"tabgen/internal/proto/tab"
	"tabgen/internal/service"

	"go.uber.org/zap"
)

type TabHandler struct {
	tab.UnimplementedTabGenerateServer
	service *service.TabService
	log     *zap.Logger
}

func NewTabHandler(service *service.TabService, log *zap.Logger) *TabHandler {
	return &TabHandler{
		service: service,
		log:     log,
	}
}

func (h *TabHandler) GenerateTab(ctx context.Context, req *tab.TabRequest) (*tab.TabResponse, error) {
	return h.service.GenerateTab(ctx, req)
}
