package consumer

import (
	"context"
	"encoding/json"
	"errors"

	"orchestrator/internal/models"

	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

type AudioSepService interface {
	HandleAudioSepTaskCompletedEvent(ctx context.Context, event *models.AudioSepTaskCompletedEvent) (string, error)
}

type AudioSepTaskCompletedConsumer struct {
	reader          *kafka.Reader
	service         AudioSepService
	tgStartProducer TabGenStartProducer
	log             *zap.Logger
}

func NewAudioSepTaskCompletedConsumer(
	reader *kafka.Reader,
	service AudioSepService,
	tgStartProducer TabGenStartProducer,
	log *zap.Logger,
) *AudioSepTaskCompletedConsumer {
	return &AudioSepTaskCompletedConsumer{
		reader:          reader,
		service:         service,
		tgStartProducer: tgStartProducer,
		log:             log,
	}
}

func (c *AudioSepTaskCompletedConsumer) Run(ctx context.Context) error {
	c.log.Info("AudioSepTaskCompletedConsumer started")

	for {
		m, err := c.reader.FetchMessage(ctx)
		if err != nil {
			if errors.Is(err, context.Canceled) {
				c.log.Info("AudioSepTaskCompletedConsumer stopped by context")
				return nil
			}

			c.log.Error("failed to fetch kafka message",
				zap.Error(err),
			)
			continue
		}

		c.log.Debug("message fetched",
			zap.String("topic", m.Topic),
			zap.Int("partition", m.Partition),
			zap.Int64("offset", m.Offset),
			zap.ByteString("key", m.Key),
		)

		if err := c.consume(ctx, m); err != nil {
			if errors.Is(err, context.Canceled) {
				c.log.Info("AudioSepTaskCompletedConsumer stopped by context")
				return nil
			}

			c.log.Error("failed to process message",
				zap.Error(err),
				zap.String("topic", m.Topic),
				zap.Int("partition", m.Partition),
				zap.Int64("offset", m.Offset),
			)

			return err
		}

		if err := c.reader.CommitMessages(ctx, m); err != nil {
			c.log.Error("failed to commit message",
				zap.Error(err),
				zap.String("topic", m.Topic),
				zap.Int("partition", m.Partition),
				zap.Int64("offset", m.Offset),
			)
			continue
		}

		c.log.Debug("message committed",
			zap.String("topic", m.Topic),
			zap.Int("partition", m.Partition),
			zap.Int64("offset", m.Offset),
		)
	}
}
func (c *AudioSepTaskCompletedConsumer) consume(ctx context.Context, m kafka.Message) error {
	var ascEvent models.AudioSepTaskCompletedEvent

	if err := json.Unmarshal(m.Value, &ascEvent); err != nil {
		c.log.Error("failed to unmarshal AudioSepTaskCompletedEvent",
			zap.Error(err),
			zap.ByteString("key", m.Key),
			zap.Int64("offset", m.Offset),
		)
		return err
	}

	log := c.log.With(
		zap.String("audio_sep_task_id", ascEvent.ID),
		zap.Int64("offset", m.Offset),
		zap.Int("partition", m.Partition),
	)

	log.Info("audio sep task completed event received")

	tabGenTaskID, err := c.service.HandleAudioSepTaskCompletedEvent(ctx, &ascEvent)
	if err != nil {
		log.Error("service failed to handle AudioSepTaskCompletedEvent",
			zap.Error(err),
		)
		return err
	}

	if tabGenTaskID == "" {
		log.Warn("audio sep task not associated with tab gen task")
		return nil
	}

	log = log.With(zap.String("tab_gen_task_id", tabGenTaskID))

	log.Info("producing TabGenTaskStartEvent")

	err = c.tgStartProducer.Produce(ctx, &models.TabGenTaskStartEvent{
		ID: tabGenTaskID,
	})
	if err != nil {
		log.Error("failed to produce TabGenTaskStartEvent",
			zap.Error(err),
		)
		return err
	}

	log.Info("TabGenTaskStartEvent successfully produced")

	return nil
}

func (c *AudioSepTaskCompletedConsumer) Close() error {
	return c.reader.Close()
}
