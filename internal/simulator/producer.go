package simulator

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/IBM/sarama"
)

type TelemetryProducer[K any, V any] struct {
	producer sarama.SyncProducer
	topic    string
}

func NewTelemetryProducer[K any, V any](brokers []string, topic string) (*TelemetryProducer[K, V], error) {
	if len(brokers) <= 0 {
		return nil, errors.New("Brokers not provided")
	}

	if len(strings.TrimSpace(topic)) <= 0 {
		return nil, errors.New("Topic not provided")
	}

	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5
	config.Producer.Return.Successes = true
	config.Producer.Partitioner = sarama.NewHashPartitioner

	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		return nil, err
	}

	return &TelemetryProducer[K, V]{
		producer: producer,
		topic:    topic,
	}, nil
}

func (p *TelemetryProducer[K, V]) Send(key K, value V) error {
	msgValue, err := json.Marshal(value)
	if err != nil {
		return err
	}

	msg := &sarama.ProducerMessage{
		Topic: p.topic,
		Key:   sarama.StringEncoder(fmt.Sprint(key)),
		Value: sarama.ByteEncoder(msgValue),
	}

	partition, offset, err := p.producer.SendMessage(msg)
	if err != nil {
		return err
	}

	log.Printf("Message stored in partition %d at offset %d\n", partition, offset)
	return nil
}

func (p *TelemetryProducer[K, V]) Close() {
	if p.producer != nil {
		p.producer.Close()
	}
}
