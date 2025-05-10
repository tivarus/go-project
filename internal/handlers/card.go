package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"bank-api/internal/models"
	"bank-api/internal/service"

	"github.com/gorilla/mux"
)

type CardHandler struct {
	cardService *service.CardService
}

func NewCardHandler(cardService *service.CardService) *CardHandler {
	return &CardHandler{cardService: cardService}
}

func (h *CardHandler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/cards", h.CreateCard).Methods("POST")
	router.HandleFunc("/cards/{id}", h.GetCard).Methods("GET")
	router.HandleFunc("/accounts/{id}/cards", h.GetAccountCards).Methods("GET")
}

func (h *CardHandler) CreateCard(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(int)

	var req models.CreateCardRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	card, err := h.cardService.CreateCard(userID, &req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	response := models.CardResponse{
		ID:         card.ID,
		AccountID:  card.AccountID,
		LastFour:   card.Number[len(card.Number)-4:],
		ExpiryDate: card.ExpiryDate,
		CreatedAt:  card.CreatedAt,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func (h *CardHandler) GetCard(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(int)
	vars := mux.Vars(r)
	cardID, _ := strconv.Atoi(vars["id"])

	card, err := h.cardService.GetCard(userID, cardID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(card)
}

func (h *CardHandler) GetAccountCards(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(int)
	vars := mux.Vars(r)
	cardID, _ := strconv.Atoi(vars["id"])

	card, err := h.cardService.GetCard(userID, cardID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(card)
}
