package infrastructure

import (
	"encoding/json"

	"github.com/ZOLUXERO/gotest/core"
	"github.com/confluentinc/confluent-kafka-go/kafka"
)

var eventTopic = "points-events"

type EventEmitter interface {
	Emit(event *core.PointsEvent) error
}
type KafkaEventEmitter struct {
	Producer *kafka.Producer
}

func (k *KafkaEventEmitter) Emit(event *core.PointsEvent) error {
	payload, err := json.Marshal(event)
	if err != nil {
		return err
	}

	message := &kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &eventTopic, Partition: kafka.PartitionAny},
		Value:          payload,
	}
	err = k.Producer.Produce(message, nil)
	return err
}

// Instancia kafka
func NewKafkaEventEmitter() (*KafkaEventEmitter, error) {
	config := &kafka.ConfigMap{
		"bootstrap.servers":        "localhost:9092",
		"go.events.channel.enable": true,
	}
	producer, err := kafka.NewProducer(config)
	if err != nil {
		return nil, err
	}
	return &KafkaEventEmitter{Producer: producer}, nil
}
