package events

import (
    "context"

    "github.com/twmb/franz-go/pkg/kgo"
)

type Producer struct {
    client *kgo.Client
}

func NewProducer(client *kgo.Client) *Producer {
    return &Producer{client: client}
}

// Produce satisfies the service.EventProducer interface
func (p *Producer) Produce(ctx context.Context, topic string, key string, value []byte) error {
    if p.client == nil {
        return nil
    }
    record := &kgo.Record{
        Topic: topic,
        Key:   []byte(key),
        Value: value,
    }
    // ProduceSync is safest for low-volume critical events
    return p.client.ProduceSync(ctx, record).FirstErr()
}
