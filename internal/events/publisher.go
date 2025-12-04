package events

import "context"

type Publisher interface {
	Publish(ctx context.Context, topic string, payload []byte) error
	Close() error
}

type noop struct{}

func (noop) Publish(ctx context.Context, topic string, payload []byte) error { return nil }
func (noop) Close() error { return nil }

func NewNoopPublisher() Publisher { return noop{} }
