package main

import (
	"os"
	"strings"
	"sync"

	"github.com/coeusj/rock-sim/internal/simulator"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load("../../test_kafka.env")

	brokers := strings.Split(os.Getenv("KAFKA_URL"), ",")
	wg := &sync.WaitGroup{}

	navigationSim, err := simulator.NewSimulator[string, simulator.Navigation](brokers, "navigation")
	if err == nil {
		navigationSim.Update = func() {
			navigationSim.Value.Altitude += 10.3
			navigationSim.Value.Velocity += 2.5
		}

		navigationSim.StartSimulation(wg, "Electron-Beta", simulator.Navigation{
			Key:      "Electron-Beta",
			Velocity: 7500.2,
			Altitude: 150000.5,
			Pitch:    90.0,
			Yaw:      0.0,
			Roll:     0.0,
		})
	}

	propulsionSim, err := simulator.NewSimulator[string, simulator.Propulsion](brokers, "propulsion")
	if err == nil {
		propulsionSim.Update = func() {
			propulsionSim.Value.FuelPerc -= 0.003
		}

		propulsionSim.StartSimulation(wg, "Electron-Beta", simulator.Propulsion{
			Key:      "Electron-Beta",
			FuelPerc: 100.00,
		})
	}

	wg.Wait()
	navigationSim.Dispose()
	propulsionSim.Dispose()
}
