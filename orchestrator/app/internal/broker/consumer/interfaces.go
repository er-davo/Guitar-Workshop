package consumer

import (
	"context"
	"orchestrator/internal/models"
)

type TabGenStartProducer interface {
	Produce(ctx context.Context, tgStart *models.TabGenTaskStartEvent) error
}
