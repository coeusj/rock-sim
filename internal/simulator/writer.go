package simulator

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/segmentio/kafka-go"
)

type RocketTelemetryWriter struct {
	Writer *kafka.Writer
}

func NewRocketTelemetryWriter(brokers []string, topicName string) (*RocketTelemetryWriter, error) {
	godotenv.Load("../../test_kafka.env")

	if len(brokers) <= 0 {
		return nil, errors.New("Brokers not provided")
	}

	if len(strings.TrimSpace(topicName)) <= 0 {
		return nil, errors.New("Topic not provided")
	}

	writer := &kafka.Writer{
		Addr:                   kafka.TCP(os.Getenv("KAFKA_URL")),
		Topic:                  topicName,
		Balancer:               &kafka.ReferenceHash{},
		Async:                  false,
		AllowAutoTopicCreation: true,
		RequiredAcks:           kafka.RequireAll,
	}

	return &RocketTelemetryWriter{
		Writer: writer,
	}, nil
}

func (w *RocketTelemetryWriter) WriteMessage(ctx context.Context, telemetry *RocketTelemetry) error {
	jsonData, err := json.Marshal(telemetry)
	if err != nil {
		return err
	}

	writeErr := w.Writer.WriteMessages(ctx,
		kafka.Message{
			Key:   []byte(telemetry.RocketID),
			Value: jsonData,
		},
	)

	if writeErr != nil {
		return writeErr
	}

	return nil
}
