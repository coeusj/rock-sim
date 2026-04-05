package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"strings"
	"sync"
	"time"

	"github.com/coeusj/rock-sim/internal/simulator"
	"github.com/coeusj/rock-sim/pkg/utils"
	"github.com/joho/godotenv"
)

func main() {
	// Command line flags for configuration
	envFile, logLevel := utils.GetCliFlags()

	// Setup structured logging
	logger := utils.CreateLogger(logLevel)
	slog.SetDefault(logger)

	// Load environment variables
	if err := godotenv.Load(envFile); err != nil {
		logger.Warn("failed to load env file, using command line or defaults", "error", err)
	}

	// Get Kafka brokers
	brokersStr := os.Getenv("KAFKA_BROKERS")
	if brokersStr == "" {
		logger.Error("Kafka URL not provided")
		os.Exit(1)
	}
	brokers := strings.Split(brokersStr, ",")

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Setup signal handling for graceful shutdown
	ctx, stop := signal.NotifyContext(ctx, os.Interrupt, os.Kill)
	defer stop()

	wg := &sync.WaitGroup{}

	// Initialize navigation simulation with initial values
	navSimulator := simulator.NewNavigationSimulation(simulator.Navigation{
		Key:      "Electron-Beta",
		Velocity: 7500.2,
		Altitude: 150000.5,
		Pitch:    90.0,
		Yaw:      0.0,
		Roll:     0.0,
	})

	// Create Kafka sender with proper error handling
	navSender, err := simulator.NewKafkaSender[string, simulator.Navigation](brokers, "navigation", logger)
	if err != nil {
		logger.Error("failed to create Kafka sender", "error", err)
	}

	defer func() {
		if err := navSender.Close(); err != nil {
			logger.Error("failed to close Kafka sender", "error", err)
		}
	}()

	// Create writer and engine
	navWriter := simulator.NewWriter(navSender)
	navEngine := simulator.NewEngine(navSimulator, navWriter, logger)

	// Start simulation
	logger.Info("starting navigation simulation", "rocket", navSimulator.Value.Key)
	navEngine.StartSimulation(ctx, wg, "Electron-Beta", navSimulator.Value)

	// Wait for completion
	wg.Wait()
	logger.Info("simulation completed successfully")
}
