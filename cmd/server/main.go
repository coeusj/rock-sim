package main

import (
	"context"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"strings"
	"sync"
	"time"

	"github.com/IBM/sarama"
	"github.com/coeusj/rock-sim/internal/simulations"
	rocket "github.com/coeusj/rock-sim/pkg/api/rocket/v1"
	"github.com/coeusj/rock-sim/pkg/utils"
	"github.com/joho/godotenv"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func main() {
	envFile, logLevel := utils.GetCliFlags()

	logger := utils.CreateLogger(logLevel)
	slog.SetDefault(logger)

	if err := godotenv.Load(envFile); err != nil {
		logger.Warn("failed to load env file, using command line or defaults", "error", err)
	}

	brokersStr := os.Getenv("KAFKA_BROKERS")
	if brokersStr == "" {
		logger.Error("Kafka URL not provided")
		os.Exit(1)
	}
	brokers := strings.Split(brokersStr, ",")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ctx, stop := signal.NotifyContext(ctx, os.Interrupt, os.Kill)
	defer stop()

	wg := &sync.WaitGroup{}

	kafkaConfig := sarama.NewConfig()
	kafkaConfig.Producer.RequiredAcks = sarama.WaitForAll
	kafkaConfig.Producer.Retry.Max = 5
	kafkaConfig.Producer.Return.Successes = true
	kafkaConfig.Producer.Partitioner = sarama.NewHashPartitioner
	asyncProducer, err := sarama.NewAsyncProducer(brokers, kafkaConfig)
	if err != nil {
		log.Fatalf("could not create Kafka Producer: %v", err)
	}
	defer asyncProducer.Close()

	navSim := simulations.NewNavigationSimulation(asyncProducer, &rocket.Navigation{
		Key:       "electron-beta-1",
		Velocity:  7500.2,
		Altitude:  150000.5,
		Pitch:     90.0,
		Yaw:       0.0,
		Roll:      0.0,
		Timestamp: timestamppb.New(time.Now()),
	})
	navSim.Start(ctx, wg)

	propulsionSim := simulations.NewPropulsionSimulation(asyncProducer, &rocket.Propulsion{
		Key:       "electron-beta-1",
		FuelPerc:  100.0,
		Timestamp: timestamppb.New(time.Now()),
	})
	propulsionSim.Start(ctx, wg)

	wg.Wait()
	log.Println("simulation completed")
}
