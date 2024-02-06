package eventbus

import (
	"context"
	"log/slog"

	"github.com/segmentio/kafka-go"
)

type EventWriter struct {
	TopicAWriter *TopicWriter
	TopicBWriter *TopicWriter
}

type TopicWriter struct {
	w *kafka.Writer
}

func NewEventWriter(brokers []string) *EventWriter {
	return &EventWriter{
		TopicAWriter: newTopicWriter(brokers, "topic-A"),
		TopicBWriter: newTopicWriter(brokers, "topic-B"),
	}
}

func newTopicWriter(brokers []string, topic string) *TopicWriter {
	return &TopicWriter{&kafka.Writer{
		Addr:     kafka.TCP(brokers...),
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	}}
}

func (tr *TopicWriter) WriteMessage(ctx context.Context, key, value string) error {
	err := tr.w.WriteMessages(ctx, kafka.Message{
		Key:   []byte(key),
		Value: []byte(value),
	})
	if err != nil {
		slog.Error("failed to write message", "topic", tr.w.Topic, "key", key, "value", value, "error", err)
		return err
	}
	slog.Info("written message", "topic", tr.w.Topic, "key", key, "value", value)
	return nil
}

func (er *EventWriter) Close() {
	if err := er.TopicAWriter.w.Close(); err != nil {
		slog.Error("failed to close writer", "error", err)
	}
}
