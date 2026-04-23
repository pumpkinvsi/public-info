package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"strings"

	confluent "github.com/confluentinc/confluent-kafka-go/v2/kafka"

	"src/email-sender/internal/config"
	"src/email-sender/internal/email"
	"src/email-sender/internal/metrics"
)

type Consumer struct {
	consumer *confluent.Consumer
	topic    string
	metrics  *metrics.Metrics
}

func New(cfg config.KafkaConfig, m *metrics.Metrics) (*Consumer, error) {
	c, err := confluent.NewConsumer(&confluent.ConfigMap{
		"bootstrap.servers":  strings.Join(cfg.Brokers, ","),
		"group.id":           cfg.GroupID,
		"auto.offset.reset":  "earliest",
		"enable.auto.commit": true,
	})
	if err != nil {
		return nil, fmt.Errorf("create kafka consumer: %w", err)
	}

	if err := c.Subscribe(cfg.Topic, nil); err != nil {
		c.Close()
		return nil, fmt.Errorf("subscribe to topic %q: %w", cfg.Topic, err)
	}

	slog.Info("kafka consumer subscribed", "topic", cfg.Topic, "group_id", cfg.GroupID)

	return &Consumer{
		consumer: c,
		topic:    cfg.Topic,
		metrics:  m,
	}, nil
}

func (c *Consumer) Listen(ctx context.Context, jobs chan<- *email.Email) {
	defer func() {
		if err := c.consumer.Close(); err != nil {
			slog.Error("kafka consumer close error", "error", err)
		} else {
			slog.Info("kafka consumer closed")
		}
	}()

	for {
		select {
		case <-ctx.Done():
			slog.Info("kafka consumer stopping", "reason", ctx.Err())
			return
		default:
		}

		ev := c.consumer.Poll(100)
		if ev == nil {
			continue
		}

		switch e := ev.(type) {
		case *confluent.Message:
			c.handleMessage(ctx, e, jobs)

		case confluent.Error:
			c.metrics.KafkaErrors.Inc()
			slog.Error("kafka consumer error",
				"code", e.Code(),
				"error", e.Error(),
				"fatal", e.IsFatal(),
			)
			if e.IsFatal() {
				slog.Error("fatal kafka error, stopping consumer")
				return
			}
		}
	}
}
func (c *Consumer) handleMessage(ctx context.Context, km *confluent.Message, jobs chan<- *email.Email) {
	c.metrics.MessagesReceived.Inc()

	var msg email.Email
	if err := json.Unmarshal(km.Value, &msg); err != nil {
		c.metrics.MessagesFailed.WithLabelValues(metrics.ReasonParseError).Inc()
		slog.Error("failed to unmarshal kafka message",
			"error", err,
			"partition", km.TopicPartition.Partition,
			"offset", km.TopicPartition.Offset,
		)
		return
	}

	slog.Debug("kafka message received",
		"from", msg.Sender,
		"contact", msg.Contact,
		"partition", km.TopicPartition.Partition,
		"offset", km.TopicPartition.Offset,
	)

	select {
	case jobs <- &msg:
	case <-ctx.Done():
		slog.Warn("context cancelled while submitting message to worker pool, message dropped",
			"from", msg.Sender,
			"contact", msg.Contact,
		)
	}
}

func (c *Consumer) Close() error {
	return c.consumer.Close()
}
