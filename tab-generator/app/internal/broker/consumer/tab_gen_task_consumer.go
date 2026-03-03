package consumer

import (
	"context"
	"encoding/json"
	"time"

	"tabgen/internal/models"
	"tabgen/internal/worker"

	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

type TabGenService interface {
	GenerateTab(ctx context.Context, task *models.TabGenTaskStartEvent) error
}

type TabGenTaskStartConsumer struct {
	reader  *kafka.Reader
	service TabGenService
	pool    *worker.Pool
	log     *zap.Logger
}

func NewTabGenTaskStartConsumer(
	reader *kafka.Reader,
	service TabGenService,
	pool *worker.Pool,
	log *zap.Logger,
) *TabGenTaskStartConsumer {
	return &TabGenTaskStartConsumer{
		reader:  reader,
		service: service,
		pool:    pool,
		log:     log,
	}
}

func (c *TabGenTaskStartConsumer) Run(ctx context.Context) {
	c.log.Info("TabGenTaskStartConsumer started")

	for {
		m, err := c.reader.FetchMessage(ctx)
		if err != nil {
			if ctx.Err() != nil {
				c.log.Info("TabGenTaskStartConsumer stopped by context")
				return
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
			if ctx.Err() != nil {
				c.log.Info("TabGenTaskStartConsumer stopped by context")
				return
			}

			c.log.Error("failed to submit task to worker pool",
				zap.Error(err),
				zap.Int64("offset", m.Offset),
			)
		}
	}
}

func (c *TabGenTaskStartConsumer) consume(ctx context.Context, msg kafka.Message) error {
	return c.pool.Submit(ctx, func(ctx context.Context) error {
		start := time.Now()

		var task models.TabGenTaskStartEvent
		if err := json.Unmarshal(msg.Value, &task); err != nil {
			c.log.Error("failed to unmarshal TabGenTaskStartEvent",
				zap.Error(err),
				zap.Int64("offset", msg.Offset),
				zap.Int("partition", msg.Partition),
			)

			if err := c.reader.CommitMessages(ctx, msg); err != nil {
				c.log.Error("failed to commit malformed message",
					zap.Error(err),
				)
			}
			return nil
		}

		log := c.log.With(
			zap.String("task_id", task.ID),
			zap.Int64("offset", msg.Offset),
			zap.Int("partition", msg.Partition),
		)

		log.Info("tab generation started")

		if err := c.service.GenerateTab(ctx, &task); err != nil {
			log.Error("tab generation failed",
				zap.Error(err),
				zap.Duration("duration", time.Since(start)),
			)
			return err
		}

		log.Info("tab generation completed",
			zap.Duration("duration", time.Since(start)),
		)

		if err := c.reader.CommitMessages(ctx, msg); err != nil {
			log.Error("failed to commit message",
				zap.Error(err),
			)
			return err
		}

		log.Debug("message committed")

		return nil
	})
}
func (c *TabGenTaskStartConsumer) Close() error {
	err := c.reader.Close()
	c.pool.Wait()
	return err
}
