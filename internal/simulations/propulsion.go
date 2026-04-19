package simulations

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/IBM/sarama"
	rocket "github.com/coeusj/rock-sim/pkg/api/rocket/v1"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type PropulsionSimulation struct {
	producer sarama.AsyncProducer
	key      string
	topic    string
	value    *rocket.Propulsion
}

func NewPropulsionSimulation(producer sarama.AsyncProducer, initialValue *rocket.Propulsion) *PropulsionSimulation {
	return &PropulsionSimulation{
		producer: producer,
		key:      initialValue.Key,
		topic:    "propulsion",
		value:    initialValue,
	}
}

func (s *PropulsionSimulation) Start(ctx context.Context, wg *sync.WaitGroup) {
	wg.Add(1)

	go func() {
		defer wg.Done()

		ticker := time.NewTicker(time.Millisecond * 50)
		defer ticker.Stop()

		for i := 0; i < 500; i++ {
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

				protoMsg, err := proto.Marshal(s.value)
				if err != nil {
					log.Printf("error while trying to marshal protobuf at iteration %d: %v\n", i, err)
					continue
				}

				msg := &sarama.ProducerMessage{
					Topic: s.topic,
					Key:   sarama.ByteEncoder(s.key),
					Value: sarama.ByteEncoder(protoMsg),
				}

				s.producer.Input() <- msg

				select {
				case success := <-s.producer.Successes():
					log.Printf("[Propulsion] - Partition: %d, Offset: %d\n", success.Partition, success.Offset)
				case err := <-s.producer.Errors():
					log.Printf("error while trying to send message at iteration %d: %v\n", i, err)
				case <-ctx.Done():
					log.Printf("simulation cancelled while trying to send message: %v, iterations completed: %d\n", ctx.Err(), i)
					return
				}
			}
		}
	}()
}

func (s *PropulsionSimulation) Stop() error {
	return nil
}

func (s *PropulsionSimulation) Update() error {
	s.value.FuelPerc -= 0.003
	s.value.Timestamp = timestamppb.New(time.Now())
	return nil
}
