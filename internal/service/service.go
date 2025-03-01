package service

import (
	"best-structure-example/internal/kafka"
	"best-structure-example/internal/model"
	"best-structure-example/internal/repository"
	"context"
	"encoding/json"
	"log"
	"time"
)

// Service provides business logic operations
type Service interface {
	CreateItem(item model.Item) (*model.Item, error)
	GetItem(id string) (*model.Item, error)
	ListItems() ([]model.Item, error)
	ProcessMessage(ctx context.Context, message []byte) error
}

type service struct {
	repo   repository.Repository
	kafka  kafka.Producer
	logger *log.Logger
}

// NewService creates a new service
func NewService(repo repository.Repository, kafka kafka.Producer, logger *log.Logger) Service {
	return &service{
		repo:   repo,
		kafka:  kafka,
		logger: logger,
	}
}

func (s *service) CreateItem(item model.Item) (*model.Item, error) {
	// Set timestamps
	now := time.Now()
	item.CreatedAt = now
	item.UpdatedAt = now

	// Save to repository
	if err := s.repo.CreateItem(&item); err != nil {
		return nil, err
	}

	// Create item event
	event := model.ItemEvent{
		EventType: "created",
		Timestamp: now,
		Item:      item,
	}

	// Publish to Kafka
	if err := s.publishEvent(event); err != nil {
		s.logger.Printf("Failed to publish item created event: %v", err)
		// Continue anyway - we successfully created the item
	}

	return &item, nil
}

func (s *service) GetItem(id string) (*model.Item, error) {
	return s.repo.GetItem(id)
}

func (s *service) ListItems() ([]model.Item, error) {
	return s.repo.ListItems()
}

func (s *service) ProcessMessage(ctx context.Context, message []byte) error {
	var event model.ItemEvent
	if err := json.Unmarshal(message, &event); err != nil {
		return err
	}

	s.logger.Printf("Processing event: %s for item: %s", event.EventType, event.Item.ID)

	return nil
}

func (s *service) publishEvent(event model.ItemEvent) error {
	data, err := json.Marshal(event)
	if err != nil {
		return err
	}

	return s.kafka.Produce(event.Item.ID, data)
}
