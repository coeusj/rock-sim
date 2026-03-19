package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/coeusj/rock-sim/internal/simulator"
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
	writeErr := telemetryWriter.WriteMessage(context.Background(), &rocketTelemetry)
	if writeErr != nil {
		log.Fatal("Could not write message:", writeErr)
		os.Exit(1)
	}

	fmt.Printf("Telemetry sent at: %v\n", time.Now())

	for i := 0; i < 5000; i++ {
		rocketTelemetry.Altitude += 15.2
		rocketTelemetry.Velocity += 5.5

		writerWg.Add(1)

		go func() {
			defer writerWg.Done()

			writeErr = telemetryWriter.WriteMessage(context.Background(), &rocketTelemetry)
			if writeErr != nil {
				log.Fatal("Could not write message:", writeErr)
			}

			fmt.Printf("Telemetry sent at: %v\n", time.Now())
		}()

		// time.Sleep(time.Second / 10000)
	}

	writerWg.Wait()
	fmt.Println("Done")
}
