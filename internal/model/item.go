package model

import (
	"time"
)

// Item represents a basic data model
type Item struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Price     float64   `json:"price"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ItemEvent represents a Kafka message payload
type ItemEvent struct {
	EventType string    `json:"event_type"` // created, updated, deleted
	Timestamp time.Time `json:"timestamp"`
	Item      Item      `json:"item"`
}
