package messaging

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"doheem-backend/internal/domain"

	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"
)

type KafkaEventBus struct {
	writer *kafka.Writer
}

func NewKafkaEventBus(brokers string) *KafkaEventBus {
	w := &kafka.Writer{
		Addr:     kafka.TCP(brokers),
		Topic:    "doheem.events",
		Balancer: &kafka.Hash{},
	}
	return &KafkaEventBus{writer: w}
}

func (b *KafkaEventBus) Close() error {
	return b.writer.Close()
}

func (b *KafkaEventBus) Publish(ctx context.Context, event domain.DomainEvent) error {
	if event.ID == "" {
		event.ID = uuid.New().String()
	}
	if event.Timestamp.IsZero() {
		event.Timestamp = time.Now()
	}
	payload, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("marshal event: %w", err)
	}
	msg := kafka.Message{
		Key:   []byte(event.ID),
		Value: payload,
	}
	return b.writer.WriteMessages(ctx, msg)
}
