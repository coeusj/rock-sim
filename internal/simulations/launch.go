package simulations

import (
	"context"
	"encoding/csv"
	"io"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/IBM/sarama"
	"github.com/coeusj/rock-sim/pkg/api/rocket/v1"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type LaunchSimulation struct {
	simFilePath string
	key         string
	topic       string
	producer    sarama.AsyncProducer
}

func NewLaunchSimulation(simFilePath string, key string, topic string, producer sarama.AsyncProducer) *LaunchSimulation {
	return &LaunchSimulation{
		simFilePath: simFilePath,
		key:         key,
		topic:       topic,
		producer:    producer,
	}
}

func (s *LaunchSimulation) Start(ctx context.Context, loopInterval time.Duration) {
	simFile, err := os.Open(s.simFilePath)
	if err != nil {
		log.Fatalf("error while trying to open simulation file: %v", err)
	}
	defer simFile.Close()

	reader := csv.NewReader(simFile)

	// Skip the header row
	_, headerReadErr := reader.Read()
	if headerReadErr != nil {
		log.Fatalf("error while trying to read header row: %v", err)
	}

	ticker := time.NewTicker(loopInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Printf("simulation cancelled")
			return
		case <-ticker.C:
			record, err := reader.Read()
			if err == io.EOF {
				log.Println("Launch simulation completed")
				return
			}

			if err != nil {
				log.Fatalf("error while trying to read row: %v", err)
			}

			altitude, _ := strconv.ParseFloat(record[3], 64)
			velocity, _ := strconv.ParseFloat(record[4], 64)
			latitude, _ := strconv.ParseFloat(record[5], 64)
			longitude, _ := strconv.ParseFloat(record[6], 64)
			chamberPsi, _ := strconv.ParseFloat(record[7], 64)
			batteryPerc, _ := strconv.ParseFloat(record[8], 64)
			timestamp, err := time.Parse(time.RFC3339, record[0])
			if err != nil {
				log.Fatalf("Failed to parse date: %v", err)
			}

			sendErr := s.Send(ctx, &rocket.Telemetry{
				Status:             record[2],
				Altitude:           altitude,
				Velocity:           velocity,
				Latitude:           latitude,
				Longitude:          longitude,
				ChamberPressurePsi: chamberPsi,
				BatteryPerc:        batteryPerc,
				Timestamp:          timestamppb.New(timestamp),
			})

			if sendErr != nil {
				log.Printf("error while trying to send data: %v", err)
				continue
			}
		}
	}
}

func (s *LaunchSimulation) Send(ctx context.Context, data *rocket.Telemetry) error {
	protoData, err := proto.Marshal(data)
	if err != nil {
		return err
	}

	msg := &sarama.ProducerMessage{
		Topic: s.topic,
		Key:   sarama.ByteEncoder(s.key),
		Value: sarama.ByteEncoder(protoData),
	}

	s.producer.Input() <- msg

	select {
	case _ = <-s.producer.Successes():
		log.Println("message sent successfully")
		return nil
	case err := <-s.producer.Errors():
		return err
	case <-ctx.Done():
		log.Println("context done while sending message")
		return nil
	}
}
