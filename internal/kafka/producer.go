package kafka

import (
	"best-structure-example/internal/config"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"log"
)

// Producer defines the Kafka producer interface
type Producer interface {
	Produce(key string, value []byte) error
	Close() error
}

type kafkaProducer struct {
	producer *kafka.Producer
	topic    string
}

// NewProducer creates a new Kafka producer
func NewProducer(config config.KafkaConfig) (Producer, error) {
	p, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": config.Brokers[0],
		"client.id":         "example-producer",
		"acks":              "all",
	})

	if err != nil {
		return nil, err
	}

	// Start goroutine to handle delivery reports
	go func() {
		for e := range p.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					log.Println("Failed to deliver message:", ev.TopicPartition)
				}
			}
		}
	}()

	return &kafkaProducer{
		producer: p,
		topic:    config.Topic,
	}, nil
}

func (p *kafkaProducer) Produce(key string, value []byte) error {
	return p.producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &p.topic, Partition: kafka.PartitionAny},
		Key:            []byte(key),
		Value:          value,
	}, nil)
}

func (p *kafkaProducer) Close() error {
	p.producer.Flush(15 * 1000) // Wait for any outstanding messages to be delivered
	p.producer.Close()
	return nil
}
