package kafka

import (
	"context"

	"github.com/segmentio/kafka-go"
)

func NewWriter(brokers []string, topic string) *kafka.Writer {
	return kafka.NewWriter(kafka.WriterConfig{
		Brokers: brokers,
		Topic:   topic,
	})
}

func Write(ctx context.Context, w *kafka.Writer, key string, value []byte) error {
	msg := kafka.Message{Value: value}
	if key != "" {
		msg.Key = []byte(key)
	}
	return w.WriteMessages(ctx, msg)
}
