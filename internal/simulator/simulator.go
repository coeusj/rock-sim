package simulator

import (
	"errors"
	"log"
	"sync"
	"time"
)

type Simulator[K any, V any] struct {
	brokers  []string
	topic    string
	producer *TelemetryProducer[K, V]
	Value    V
	Update   func()
}

func NewSimulator[K any, V any](brokers []string, topic string) (*Simulator[K, V], error) {
	telemetryProducer, err := NewTelemetryProducer[K, V](brokers, topic)
	if err != nil {
		log.Fatalf("Could not create instance of Simulator: %v", err.Error())
		return nil, err
	}

	return &Simulator[K, V]{
		brokers:  brokers,
		topic:    topic,
		producer: telemetryProducer,
	}, nil
}

func (s *Simulator[K, V]) StartSimulation(wg *sync.WaitGroup, key K, intialValue V) error {
	if s.Update == nil {
		return errors.New("Update is undefined")
	}

	wg.Add(1)

	go func() {
		for i := 0; i < 30; i++ {
			s.Update()

			err := s.producer.Send(key, s.Value)
			if err != nil {
				log.Fatalf("Could not send message: %v", err.Error())
				continue
			}

			time.Sleep(time.Second / 100)
		}

		wg.Done()
	}()

	return nil
}

func (s *Simulator[K, V]) Dispose() {
	if s != nil {
		s.producer.Close()
	}
}
