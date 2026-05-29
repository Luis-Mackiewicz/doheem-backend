package messaging

import (
	"context"

	"doheem-backend/internal/domain"
)

type NoopEventBus struct{}

func NewNoopEventBus() *NoopEventBus {
	return &NoopEventBus{}
}

func (b *NoopEventBus) Publish(_ context.Context, _ domain.DomainEvent) error {
	return nil
}
