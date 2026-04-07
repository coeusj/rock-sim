package simulations

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/IBM/sarama"
)

type NavigationSimulation struct {
	*Simulation[Navigation]
}

func NewNavigationSimulation(producer sarama.SyncProducer, initialValue Navigation) *NavigationSimulation {
	return &NavigationSimulation{
		Simulation: NewSimulation(producer, "electron-beta-navigation", "navigation", initialValue),
	}
}

func (ns *NavigationSimulation) Start(ctx context.Context, wg *sync.WaitGroup) {
	wg.Add(1)

	go func() {
		defer wg.Done()

		ticker := time.NewTicker(time.Millisecond * 100)
		defer ticker.Stop()

		for i := 0; i < 100; i++ {
			select {
			case <-ctx.Done():
				log.Printf("simulation cancelled: %v, iterations completed: %d\n", ctx.Err(), i)
				return
			case <-ticker.C:
				err := ns.Update()
				if err != nil {
					log.Printf("error while trying to update simulation at iteration %d: %v\n", i, err)
					continue
				}

				jsonMsg, err := json.Marshal(ns.value)
				if err != nil {
					log.Printf("error while trying to marshal the message at iteration %d: %v\n", i, err)
					continue
				}

				msg := &sarama.ProducerMessage{
					Topic: ns.topic,
					Key:   sarama.StringEncoder(fmt.Sprint(ns.key)),
					Value: sarama.ByteEncoder(jsonMsg),
				}

				partition, offset, err := ns.producer.SendMessage(msg)
				if err != nil {
					log.Printf("error while trying to send message at iteration %d: %v\n", i, err)
					continue
				}

				log.Printf("Partition: %d, Offset: %d\n", partition, offset)
			}
		}
	}()
}

func (ns *NavigationSimulation) Stop() error {
	err := ns.producer.Close()
	return err
}

func (ns *NavigationSimulation) Update() error {
	ns.value.Altitude += 10.3
	ns.value.Velocity += 2.5
	return nil
}
