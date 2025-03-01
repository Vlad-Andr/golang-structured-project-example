package kafka

import (
	"best-structure-example/internal/config"
	"context"
	"log"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

// MessageHandler is a function that processes Kafka messages
type MessageHandler func(ctx context.Context, message []byte) error

// Consumer defines the Kafka consumer interface
type Consumer interface {
	Start(ctx context.Context, handler MessageHandler) error
	Stop() error
}

type kafkaConsumer struct {
	consumer *kafka.Consumer
	topic    string
	logger   *log.Logger
	stopCh   chan struct{}
}

// NewConsumer creates a new Kafka consumer
func NewConsumer(config config.KafkaConfig, logger *log.Logger) (Consumer, error) {
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":     config.Brokers[0],
		"broker.address.family": "v4",
		"group.id":              config.ConsumerGroup,
		"auto.offset.reset":     "earliest",
	})

	if err != nil {
		return nil, err
	}

	if err := c.SubscribeTopics([]string{config.Topic}, nil); err != nil {
		c.Close()
		return nil, err
	}

	return &kafkaConsumer{
		consumer: c,
		topic:    config.Topic,
		logger:   logger,
		stopCh:   make(chan struct{}),
	}, nil
}

func (c *kafkaConsumer) Start(ctx context.Context, handler MessageHandler) error {
	c.logger.Printf("Starting Kafka consumer for topic: %s", c.topic)

	go func() {
		run := true
		for run {
			select {
			case <-c.stopCh:
				run = false
			default:
				msg, err := c.consumer.ReadMessage(1 * time.Second)
				if err != nil {
					if err.(kafka.Error).Code() != kafka.ErrTimedOut {
						c.logger.Printf("Consumer error: %v", err)
					}
					continue
				}

				if err := handler(ctx, msg.Value); err != nil {
					c.logger.Printf("Failed to process message: %v", err)
				}
			}
		}
		c.logger.Println("Kafka consumer stopped")
	}()

	return nil
}

func (c *kafkaConsumer) Stop() error {
	close(c.stopCh)
	return c.consumer.Close()
}
