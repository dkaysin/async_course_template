package eventbus

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/segmentio/kafka-go"
)

func StartReaders(brokers []string, groupID string) {

	// topic A
	topicAReader := newTopicReader(brokers, groupID, "topic-A")
	go handleMessage(context.Background(), topicAReader, echo)

	// topic B
	topicBReader := newTopicReader(brokers, groupID, "topic-B")
	go handleMessage(context.Background(), topicBReader, echo)

}

func newTopicReader(brokers []string, groupID string, topic string) *kafka.Reader {
	return kafka.NewReader(kafka.ReaderConfig{
		Brokers:   brokers,
		GroupID:   groupID,
		Topic:     topic,
		Partition: 0,
		MaxBytes:  10e6, // 10MB
	})
}

func closeReader(r *kafka.Reader) {
	if err := r.Close(); err != nil {
		slog.Error("failed to close reader", "error", err)
		os.Exit(1)
	}
}

type messageHandler func(m kafka.Message) error

func handleMessage(ctx context.Context, r *kafka.Reader, fn messageHandler) {
	defer closeReader(r)
	for {
		m, err := r.ReadMessage(ctx)
		if err != nil {
			slog.Error("error while reading message", "error", err)
			break
		}
		slog.Info("received message from broker", "topic", r.Config().Topic, "key", string(m.Key), "value", string(m.Value))
		if err := fn(m); err != nil {
			slog.Error("error while handling message", "error", err)
		}
	}
}

func echo(m kafka.Message) error {
	fmt.Printf("message at offset %d: %s = %s\n", m.Offset, string(m.Key), string(m.Value))
	return nil
}
