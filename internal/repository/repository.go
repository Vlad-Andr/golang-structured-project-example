package repository

import (
	"best-structure-example/internal/model"
	"log"
	"sync"

	"github.com/google/uuid"
)

// Repository provides access to the storage
type Repository interface {
	CreateItem(item *model.Item) error
	GetItem(id string) (*model.Item, error)
	ListItems() ([]model.Item, error)
}

type inMemoryRepository struct {
	items  map[string]model.Item
	mu     sync.RWMutex
	logger *log.Logger
}

// NewRepository creates a new in-memory repository
func NewRepository(logger *log.Logger) Repository {
	return &inMemoryRepository{
		items:  make(map[string]model.Item),
		logger: logger,
	}
}

func (r *inMemoryRepository) CreateItem(item *model.Item) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Generate ID if not provided
	if item.ID == "" {
		item.ID = uuid.New().String()
	}

	r.items[item.ID] = *item
	r.logger.Printf("Item created: %s", item.ID)
	return nil
}

func (r *inMemoryRepository) GetItem(id string) (*model.Item, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	item, ok := r.items[id]
	if !ok {
		return nil, nil // Not found
	}
	return &item, nil
}

func (r *inMemoryRepository) ListItems() ([]model.Item, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	items := make([]model.Item, 0, len(r.items))
	for _, item := range r.items {
		items = append(items, item)
	}
	return items, nil
}
