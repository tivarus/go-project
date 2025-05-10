package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"bank-api/internal/models"
	"bank-api/internal/service"

	"github.com/gorilla/mux"
)

type AccountHandler struct {
	accountService *service.AccountService
}

func NewAccountHandler(accountService *service.AccountService) *AccountHandler {
	return &AccountHandler{accountService: accountService}
}

func (h *AccountHandler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/accounts", h.CreateAccount).Methods("POST")
	router.HandleFunc("/accounts/{id}/balance", h.UpdateBalance).Methods("PATCH")
	router.HandleFunc("/accounts/{id}/transfer", h.Transfer).Methods("POST")
	router.HandleFunc("/accounts/{id}", h.GetAccount).Methods("GET")
}

func (h *AccountHandler) CreateAccount(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(int)

	var req models.CreateAccountRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	account, err := h.accountService.CreateAccount(userID, &req)
	if err != nil {
		http.Error(w, "Failed to create account", http.StatusInternalServerError)
		return
	}

	response := models.AccountResponse{
		ID:        account.ID,
		Balance:   account.Balance,
		Currency:  account.Currency,
		CreatedAt: account.CreatedAt,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func (h *AccountHandler) UpdateBalance(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(int)
	vars := mux.Vars(r)
	accountID, _ := strconv.Atoi(vars["id"])

	var req models.UpdateBalanceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Проверяем что счет принадлежит пользователю
	account, err := h.accountService.GetAccountByID(accountID)
	if err != nil {
		http.Error(w, "Account error", http.StatusInternalServerError)
		return
	}
	if account == nil || account.UserID != userID {
		http.Error(w, "Account not found", http.StatusNotFound)
		return
	}

	if err := h.accountService.UpdateBalance(accountID, req.Amount); err != nil {
		switch err {
		case service.ErrInsufficientFunds:
			http.Error(w, "Insufficient funds", http.StatusBadRequest)
		default:
			http.Error(w, "Update failed", http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Balance updated"})
}

func (h *AccountHandler) Transfer(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(int)
	vars := mux.Vars(r)
	accountID, _ := strconv.Atoi(vars["id"])

	var req models.TransferRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Проверяем что счет отправителя принадлежит пользователю
	if req.FromAccountID != accountID {
		http.Error(w, "Invalid source account", http.StatusForbidden)
		return
	}

	fromAccount, err := h.accountService.GetAccountByID(req.FromAccountID)
	if err != nil {
		http.Error(w, "Account error", http.StatusInternalServerError)
		return
	}
	if fromAccount == nil || fromAccount.UserID != userID {
		http.Error(w, "Source account not found", http.StatusNotFound)
		return
	}

	if err := h.accountService.Transfer(&req); err != nil {
		switch err {
		case service.ErrAccountNotFound:
			http.Error(w, "Account not found", http.StatusNotFound)
		case service.ErrInsufficientFunds:
			http.Error(w, "Insufficient funds", http.StatusBadRequest)
		default:
			http.Error(w, "Transfer failed", http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Transfer successful"})
}

func (h *AccountHandler) GetAccount(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(int)
	vars := mux.Vars(r)
	accountID, _ := strconv.Atoi(vars["id"])

	account, err := h.accountService.GetAccountByID(accountID)
	if err != nil {
		http.Error(w, "Account error", http.StatusInternalServerError)
		return
	}
	if account == nil || account.UserID != userID {
		http.Error(w, "Account not found", http.StatusNotFound)
		return
	}

	response := models.AccountResponse{
		ID:        account.ID,
		Balance:   account.Balance,
		Currency:  account.Currency,
		CreatedAt: account.CreatedAt,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
