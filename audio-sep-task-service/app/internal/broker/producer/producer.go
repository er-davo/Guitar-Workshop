package producer

import (
	"context"
	"encoding/json"

	"audio-sep-task-service/internal/models"

	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

type AudioSepTaskProducer struct {
	writer *kafka.Writer
	log    *zap.Logger
}

func NewAudioSepTaskProducer(writer *kafka.Writer, log *zap.Logger) *AudioSepTaskProducer {
	return &AudioSepTaskProducer{
		writer: writer,
		log:    log,
	}
}

func (p *AudioSepTaskProducer) Produce(ctx context.Context, task *models.StartAudioSepTask) error {
	val, err := json.Marshal(task)
	if err != nil {
		p.log.Error("failed to marshal audio separation task", zap.String("task_id", task.ID), zap.Error(err))
		return err
	}

	err = p.writer.WriteMessages(ctx, kafka.Message{
		Key:   []byte(task.ID),
		Value: val,
	})
	if err != nil {
		p.log.Error("failed to produce audio separation message", zap.String("task_id", task.ID), zap.Error(err))
		return err
	}

	p.log.Info("audio separation message produced", zap.String("task_id", task.ID))
	return nil
}

func (p *AudioSepTaskProducer) Close() error {
	if err := p.writer.Close(); err != nil {
		p.log.Error("failed to close kafka writer", zap.Error(err))
		return err
	}
	p.log.Info("kafka writer closed")
	return nil
}
