package main

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/coeusj/rock-sim/internal/simulator"
	"github.com/coeusj/rock-sim/pkg/logger"
)

func main() {
	telemetryWriter := simulator.NewRocketTelemetryWriter()

	defer telemetryWriter.Writer.Close()

	rocketTelemetry := simulator.RocketTelemetry{
		RocketID: "Electron-Beta",
		Velocity: 7500.2,
		Altitude: 150000.5,
	}

	var writerWg sync.WaitGroup
	for i := 0; i < 10000; i++ {
		rocketTelemetry.Altitude += 15.2
		rocketTelemetry.Velocity += 5.5

		writerWg.Add(1)

		go func() {
			defer writerWg.Done()

			writeErr := telemetryWriter.WriteMessage(context.Background(), &rocketTelemetry)
			if writeErr != nil {
				log.Fatal("Could not write message:", writeErr)
			}

			log.Printf("Telemetry sent at: %v\n", logger.DateTimeWithNanoseconds(time.Now()))
		}()

		time.Sleep(time.Second / 100)
	}

	writerWg.Wait()
	log.Println("Done")
}
