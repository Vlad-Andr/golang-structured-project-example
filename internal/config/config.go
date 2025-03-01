package config

import (
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"log"
)

// Config represents the application configuration
type Config struct {
	Server ServerConfig
	Kafka  KafkaConfig
}

// ServerConfig holds HTTP server configuration
type ServerConfig struct {
	Addr         string `envconfig:"SERVER_ADDR" default:":8081"`
	ReadTimeout  int    `envconfig:"SERVER_READ_TIMEOUT" default:"5"`
	WriteTimeout int    `envconfig:"SERVER_WRITE_TIMEOUT" default:"10"`
	IdleTimeout  int    `envconfig:"SERVER_IDLE_TIMEOUT" default:"12000"`
}

// KafkaConfig holds Kafka configuration
type KafkaConfig struct {
	Brokers       []string `envconfig:"KAFKA_BROKERS" default:"localhost:9093"`
	ConsumerGroup string   `envconfig:"KAFKA_CONSUMER_GROUP" default:"example-service-group"`
	Topic         string   `envconfig:"KAFKA_TOPIC" default:"example-topic"`
}

func Load() (*Config, error) {
	// Load .env file and set environment variables
	if err := godotenv.Load("internal/config/.env"); err != nil {
		log.Println("No .env file found")
	}

	var cfg Config
	// Load environment variables
	if err := envconfig.Process("", &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
