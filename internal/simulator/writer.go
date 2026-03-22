package simulator

import (
	"context"
	"encoding/json"
	"os"

	"github.com/joho/godotenv"
	"github.com/segmentio/kafka-go"
)

type RocketTelemetryWriter struct {
	Writer *kafka.Writer
}

func NewRocketTelemetryWriter() *RocketTelemetryWriter {
	godotenv.Load("../../test_kafka.env")

	writer := &kafka.Writer{
		Addr:                   kafka.TCP(os.Getenv("KAFKA_URL")),
		Topic:                  os.Getenv("KAFKA_TOPIC_ROCKET_TELEMETRY"),
		Balancer:               &kafka.ReferenceHash{},
		Async:                  false,
		AllowAutoTopicCreation: true,
		RequiredAcks:           kafka.RequireAll,
	}

	return &RocketTelemetryWriter{
		Writer: writer,
	}
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
