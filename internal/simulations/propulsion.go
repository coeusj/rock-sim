package simulations

import (
	"context"
	"encoding/json"
	"log"
	"sync"
	"time"

	"github.com/IBM/sarama"
)

type PropulsionSimulation struct {
	producer sarama.SyncProducer
	key      string
	topic    string
	value    Propulsion
}

func NewPropulsionSimulation(producer sarama.SyncProducer, initialValue Propulsion) *PropulsionSimulation {
	return &PropulsionSimulation{
		producer: producer,
		key:      "electron-beta-propulsion",
		topic:    "propulsion",
		value:    initialValue,
	}
}

func (s *PropulsionSimulation) Start(ctx context.Context, wg *sync.WaitGroup) {
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
				err := s.Update()
				if err != nil {
					log.Printf("error while trying to update simulation at iteration %d: %v\n", i, err)
					continue
				}

				jsonMsg, err := json.Marshal(s.value)
				if err != nil {
					log.Printf("error while trying to marshal JSON at iteration %d: %v\n", i, err)
					continue
				}

				msg := &sarama.ProducerMessage{
					Topic: s.topic,
					Key:   sarama.StringEncoder(s.key),
					Value: sarama.ByteEncoder(jsonMsg),
				}

				partition, offset, err := s.producer.SendMessage(msg)
				if err != nil {
					log.Printf("error while trying to send message to Kafka at iteration %d: %v\n", i, err)
					continue
				}

				log.Printf("[Propulsion] - Partition: %d, Offset: %d\n", partition, offset)
			}
		}
	}()
}

func (s *PropulsionSimulation) Stop() error {
	return nil
}

func (s *PropulsionSimulation) Update() error {
	s.value.FuelPerc -= 0.003
	return nil
}
