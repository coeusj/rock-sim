package simulator

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"strings"

	"github.com/IBM/sarama"
)

type Sender[K any, V any] interface {
	Send(key K, value V) error
	Close() error
}

type Writer[K any, V any] struct {
	sender Sender[K, V]
}

// NewWriter creates a new telemetry writer
func NewWriter[K any, V any](s Sender[K, V]) *Writer[K, V] {
	return &Writer[K, V]{sender: s}
}

// Write sends telemetry data with context support for cancellation
func (w *Writer[K, V]) Write(ctx context.Context, key K, value V) error {
	// Check if context is cancelled before sending
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	return w.sender.Send(key, value)
}

type KafkaSender[K any, V any] struct {
	producer sarama.SyncProducer
	topic    string
	logger   *slog.Logger
}

func NewKafkaSender[K any, V any](brokers []string, topic string, logger *slog.Logger) (*KafkaSender[K, V], error) {
	if len(brokers) <= 0 {
		return nil, errors.New("brokers not provided")
	}

	if len(strings.TrimSpace(topic)) <= 0 {
		return nil, errors.New("topic not provided")
	}

	if logger == nil {
		logger = slog.Default()
	}

	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5
	config.Producer.Return.Successes = true
	config.Producer.Partitioner = sarama.NewHashPartitioner

	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create Kafka producer: %w", err)
	}

	return &KafkaSender[K, V]{
		producer: producer,
		topic:    topic,
		logger:   logger,
	}, nil
}

func (s *KafkaSender[K, V]) Send(key K, value V) error {
	msgValue, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	msg := &sarama.ProducerMessage{
		Topic: s.topic,
		Key:   sarama.StringEncoder(fmt.Sprint(key)),
		Value: sarama.ByteEncoder(msgValue),
	}

	partition, offset, err := s.producer.SendMessage(msg)
	if err != nil {
		return fmt.Errorf("failed to send message to Kafka: %w", err)
	}

	s.logger.Debug("message sent to Kafka",
		"topic", s.topic,
		"partition", partition,
		"offset", offset,
		"key", key)

	return nil
}

func (s *KafkaSender[K, V]) Close() error {
	return s.producer.Close()
}
