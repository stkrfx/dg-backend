package handler

import (
	"encoding/json"
	"net/http"

	"github.com/stkrfx/dg-backend/internal/service"
)

type ClientHandler struct {
	svc *service.ClientService
}

func NewClientHandler(s *service.ClientService) *ClientHandler {
	return &ClientHandler{svc: s}
}

// CreateClient: POST /clients
func (h *ClientHandler) CreateClient(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}

	// 1. Decode JSON (Like req.body in Express)
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	// 2. Call Service with Request Context
	// r.Context() carries the deadline/cancel signal of the HTTP request
	if err := h.svc.RegisterClient(r.Context(), body.Email, body.Name); err != nil {
		http.Error(w, "Failed to register client: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// 3. Professional Response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Client created successfully",
	})
}

// GetClients: GET /clients
func (h *ClientHandler) GetClients(w http.ResponseWriter, r *http.Request) {
	clients, err := h.svc.FetchAllClients(r.Context())
	if err != nil {
		http.Error(w, "Failed to fetch clients", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(clients)
}
