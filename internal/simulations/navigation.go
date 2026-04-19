package simulations

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/IBM/sarama"
	navigation "github.com/coeusj/rock-sim/pkg/api/rocket/v1"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type NavigationSimulation struct {
	producer sarama.AsyncProducer
	key      string
	topic    string
	value    *navigation.Navigation
}

func NewNavigationSimulation(producer sarama.AsyncProducer, initialValue *navigation.Navigation) *NavigationSimulation {
	return &NavigationSimulation{
		producer: producer,
		key:      initialValue.Key,
		topic:    "navigation",
		value:    initialValue,
	}
}

func (ns *NavigationSimulation) Start(ctx context.Context, wg *sync.WaitGroup) {
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
				err := ns.Update()
				if err != nil {
					log.Printf("error while trying to update simulation at iteration %d: %v\n", i, err)
					continue
				}

				msgProto, err := proto.Marshal(ns.value)
				if err != nil {
					log.Printf("error while trying to marshal the message at iteration %d: %v\n", i, err)
					continue
				}

				msg := &sarama.ProducerMessage{
					Topic: ns.topic,
					Key:   sarama.ByteEncoder(ns.key),
					Value: sarama.ByteEncoder(msgProto),
				}

				ns.producer.Input() <- msg

				select {
				case success := <-ns.producer.Successes():
					log.Printf("[Navigation] - Partition: %d, Offset: %d\n", success.Partition, success.Offset)
				case err := <-ns.producer.Errors():
					log.Printf("error while trying to send message at iteration %d: %v\n", i, err)
				case <-ctx.Done():
					log.Printf("simulation cancelled while trying to send message: %v, iterations completed: %d\n", ctx.Err(), i)
					return
				}
			}
		}
	}()
}

func (ns *NavigationSimulation) Stop() error {
	return nil
}

func (ns *NavigationSimulation) Update() error {
	ns.value.Altitude += 10.3
	ns.value.Velocity += 2.5
	ns.value.Timestamp = timestamppb.New(time.Now())
	return nil
}
