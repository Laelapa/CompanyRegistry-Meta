package service

import (
	"context"
)

type EventProducer interface {
	Produce(ctx context.Context, topic string, key string, value []byte) error
}
