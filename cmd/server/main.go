package main

import (
	"context"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/IBM/sarama"
	"github.com/coeusj/rock-sim/internal/simulations"
	"github.com/coeusj/rock-sim/pkg/utils"
	"github.com/joho/godotenv"
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

	sim := simulations.NewLaunchSimulation("../../rocket_telemetry_sim.csv", "rocket-sim-1", "rocket-sim", asyncProducer)
	sim.Start(ctx, (time.Second / 6))
}
