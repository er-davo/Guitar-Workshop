package produser

import (
	"context"
	"encoding/json"

	"api-gateway/internal/models"

	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

type TabGenTaskProduser struct {
	writer *kafka.Writer
	log    *zap.Logger
}

func NewTabGenTaskProduser(writer *kafka.Writer, log *zap.Logger) *TabGenTaskProduser {
	return &TabGenTaskProduser{
		writer: writer,
		log:    log,
	}
}

func (p *TabGenTaskProduser) Produce(ctx context.Context, task *models.StartTabGenTask) error {
	val, err := json.Marshal(task)
	if err != nil {
		p.log.Error("failed to marshal tab generation task",
			zap.String("task_id", task.ID),
			zap.Error(err),
		)
		return err
	}

	err = p.writer.WriteMessages(ctx, kafka.Message{
		Key:   []byte(task.ID),
		Value: val,
	})
	if err != nil {
		p.log.Error("failed to produce tab generation message",
			zap.String("task_id", task.ID),
			zap.Error(err),
		)
		return err
	}

	p.log.Info("tab generation message produced", zap.String("task_id", task.ID))
	return nil
}

func (p *TabGenTaskProduser) Close() error {
	if err := p.writer.Close(); err != nil {
		p.log.Error("failed to close kafka writer", zap.Error(err))
		return err
	}
	p.log.Info("kafka writer closed")
	return nil
}
