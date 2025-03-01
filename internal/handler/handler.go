package handler

import (
	"best-structure-example/internal/model"
	"best-structure-example/internal/service"
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// Handler handles HTTP requests
type Handler struct {
	service service.Service
	logger  *log.Logger
}

// NewHandler creates a new Handler
func NewHandler(service service.Service, logger *log.Logger) *Handler {
	return &Handler{
		service: service,
		logger:  logger,
	}
}

// Routes returns the router with all routes
func (h *Handler) Routes() http.Handler {
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Routes
	r.Get("/health", h.healthCheck)

	r.Route("/api/v1/items", func(r chi.Router) {
		r.Get("/", h.listItems)
		r.Post("/", h.createItem)
		r.Get("/{id}", h.getItem)
	})

	return r
}

func (h *Handler) healthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func (h *Handler) createItem(w http.ResponseWriter, r *http.Request) {
	var item model.Item
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		h.respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	createdItem, err := h.service.CreateItem(item)
	if err != nil {
		h.respondWithError(w, http.StatusInternalServerError, "Failed to create item")
		return
	}

	h.respondWithJSON(w, http.StatusCreated, createdItem)
}

func (h *Handler) getItem(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	item, err := h.service.GetItem(id)
	if err != nil {
		h.respondWithError(w, http.StatusInternalServerError, "Failed to get item")
		return
	}

	if item == nil {
		h.respondWithError(w, http.StatusNotFound, "Item not found")
		return
	}

	h.respondWithJSON(w, http.StatusOK, item)
}

func (h *Handler) listItems(w http.ResponseWriter, r *http.Request) {
	items, err := h.service.ListItems()
	if err != nil {
		h.respondWithError(w, http.StatusInternalServerError, "Failed to list items")
		return
	}

	h.respondWithJSON(w, http.StatusOK, items)
}

func (h *Handler) respondWithError(w http.ResponseWriter, code int, message string) {
	h.respondWithJSON(w, code, map[string]string{"error": message})
}

func (h *Handler) respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(payload)
}
