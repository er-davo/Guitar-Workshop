package consumer

import (
	"context"
	"encoding/json"
	"errors"

	"orchestrator/internal/models"

	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

type TabGenEventHandler interface {
	HandleTabGenEvent(ctx context.Context, e *models.TabGenTaskRequestedEvent) (bool, error)
}

type TabGenTaskRequestedConsumer struct {
	reader          *kafka.Reader
	service         TabGenEventHandler
	tgStartProducer TabGenStartProducer
	log             *zap.Logger
}

func NewTabGenTaskRequstedConsumer(
	reader *kafka.Reader,
	service TabGenEventHandler,
	tgStartProducer TabGenStartProducer,
	log *zap.Logger,
) *TabGenTaskRequestedConsumer {
	return &TabGenTaskRequestedConsumer{
		reader:          reader,
		service:         service,
		tgStartProducer: tgStartProducer,
		log:             log,
	}
}
func (c *TabGenTaskRequestedConsumer) Run(ctx context.Context) error {
	c.log.Info("TabGenTaskRequestedConsumer started")

	for {
		m, err := c.reader.FetchMessage(ctx)
		if err != nil {
			if errors.Is(err, context.Canceled) {
				c.log.Info("TabGenTaskRequestedConsumer stopped by context")
				return nil
			}
			c.log.Error("failed to fetch kafka message",
				zap.Error(err),
			)
			return err
		}

		c.log.Debug("message fetched",
			zap.String("topic", m.Topic),
			zap.Int("partition", m.Partition),
			zap.Int64("offset", m.Offset),
			zap.ByteString("key", m.Key),
		)

		if err := c.consume(ctx, m); err != nil {
			if errors.Is(err, context.Canceled) {
				c.log.Info("TabGenTaskRequestedConsumer stopped by context")
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
			return err
		}

		c.log.Debug("message committed",
			zap.String("topic", m.Topic),
			zap.Int("partition", m.Partition),
			zap.Int64("offset", m.Offset),
		)
	}
}

func (c *TabGenTaskRequestedConsumer) consume(ctx context.Context, m kafka.Message) error {
	var tgtask models.TabGenTaskRequestedEvent

	if err := json.Unmarshal(m.Value, &tgtask); err != nil {
		c.log.Error("failed to unmarshal TabGenTaskRequestedEvent",
			zap.Error(err),
			zap.ByteString("key", m.Key),
			zap.Int64("offset", m.Offset),
		)
		return err
	}

	log := c.log.With(
		zap.String("task_id", tgtask.ID),
		zap.Int64("offset", m.Offset),
		zap.Int("partition", m.Partition),
	)

	log.Info("tab gen task received")

	ok, err := c.service.HandleTabGenEvent(ctx, &tgtask)
	if err != nil {
		log.Error("service failed to handle tab gen event",
			zap.Error(err),
		)
		return err
	}

	if !ok {
		log.Warn("tab gen event skipped by service logic")
		return nil
	}

	log.Info("producing TabGenTaskStartEvent")

	err = c.tgStartProducer.Produce(ctx, &models.TabGenTaskStartEvent{
		ID: tgtask.ID,
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

func (c *TabGenTaskRequestedConsumer) Close() error {
	return c.reader.Close()
}
