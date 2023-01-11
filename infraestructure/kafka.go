package infrastructure

import (
	"encoding/json"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type KafkaEventEmitter struct {
	producer *kafka.Producer
}

func (k *kafkaEventEmitter) Emit(event *core.PointsEvent) error {
	payload, err := json.Marshal(event)
	if err != nil {
		return err
	}

	message := &kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &eventTopic, Partition: kafka.PartitionAny},
		Value:          payload,
	}
	_, err := k.producer.Produce(message, nil)
	return err
}

// NewKafkaEventEmitter returns a new instance of KafkaEventEmitter
func NewKafkaEventEmitter() (*KafkaEventEmitter, error) {
	config := &kafka.ConfigMap{
		"bootstrap.servers":        "localhost:9092",
		"go.events.channel.enable": true,
	}
	producer, err := kafka.NewProducer(config)
	if err != nil {
		return nil, err
	}
	return &KafkaEventEmitter{producer: producer}, nil
}
