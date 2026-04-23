package kafka

import (
	"context"
	"fmt"
	"strings"

	confluent "github.com/confluentinc/confluent-kafka-go/v2/kafka"

	"src/backend/internal/shared/config"
)

type Producer struct {
	producer *confluent.Producer
	topic    string
}

func New(cfg config.KafkaConfig) (*Producer, error) {
	p, err := confluent.NewProducer(&confluent.ConfigMap{
		"bootstrap.servers": strings.Join(cfg.Brokers, ","),
	})
	if err != nil {
		return nil, fmt.Errorf("create kafka producer: %w", err)
	}

	return &Producer{
		producer: p,
		topic:    cfg.Topic,
	}, nil
}

func (p *Producer) Publish(ctx context.Context, value []byte) error {
	deliveryCh := make(chan confluent.Event, 1)

	err := p.producer.Produce(&confluent.Message{
		TopicPartition: confluent.TopicPartition{
			Topic:     &p.topic,
			Partition: confluent.PartitionAny,
		},
		Value: value,
	}, deliveryCh)
	if err != nil {
		return fmt.Errorf("enqueue message: %w", err)
	}

	select {
	case event := <-deliveryCh:
		msg, ok := event.(*confluent.Message)
		if !ok {
			return fmt.Errorf("unexpected event type: %T", event)
		}
		if msg.TopicPartition.Error != nil {
			return fmt.Errorf("delivery failed: %w", msg.TopicPartition.Error)
		}
		return nil

	case <-ctx.Done():
		return fmt.Errorf("publish cancelled: %w", ctx.Err())
	}
}

func (p *Producer) Close() error {
	if unflashed := p.producer.Flush(5000); unflashed > 0 {
		return fmt.Errorf("failed to flush %d messages", unflashed)
	}

	defer p.producer.Close()
	return nil
}
