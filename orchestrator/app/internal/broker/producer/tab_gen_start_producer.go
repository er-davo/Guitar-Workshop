package producer

import (
	"context"
	"encoding/json"
	"orchestrator/internal/models"

	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

type TabGenStartProducer struct {
	writer *kafka.Writer
	log    *zap.Logger
}

func NewTabGenStartProducer(writer *kafka.Writer, log *zap.Logger) *TabGenStartProducer {
	return &TabGenStartProducer{
		writer: writer,
		log:    log,
	}
}

func (p *TabGenStartProducer) Produce(ctx context.Context, tgStart *models.TabGenTaskStartEvent) error {
	val, err := json.Marshal(tgStart)
	if err != nil {
		p.log.Error("failed to marshal tab gen start event", zap.Error(err))
		return err
	}

	err = p.writer.WriteMessages(ctx, kafka.Message{
		Key:   []byte(tgStart.ID),
		Value: val,
	})
	if err != nil {
		p.log.Error("failed to write tab gen start event", zap.Error(err))
		return err
	}

	p.log.Info("tab gen start event produced", zap.String("id", tgStart.ID))

	return nil
}

func (p *TabGenStartProducer) Close() error {
	return p.writer.Close()
}
