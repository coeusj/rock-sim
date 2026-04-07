package simulations

import (
	"context"
	"sync"

	"github.com/IBM/sarama"
)

type Simulator interface {
	Start(ctx context.Context, wg *sync.WaitGroup)
	Stop() error
	Update() error
}

type Simulation[V any] struct {
	producer sarama.SyncProducer
	key      string
	topic    string
	value    V
}

func NewSimulation[V any](producer sarama.SyncProducer, key string, topic string, initialValue V) *Simulation[V] {
	return &Simulation[V]{
		producer: producer,
		key:      key,
		topic:    topic,
		value:    initialValue,
	}
}

func (s *Simulation[V]) Start(ctx context.Context, wg *sync.WaitGroup) {
	panic("Start() method not implemented")
}

func (s *Simulation[V]) Stop() error {
	panic("Stop() method not implemented")
}

func (s *Simulation[V]) Update() error {
	panic("Update() method not implemented")
}
