package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"bank-api/internal/service"

	"github.com/gorilla/mux"
)

type TransactionHandler struct {
	transactionSvc *service.TransactionService
}

func NewTransactionHandler(transactionSvc *service.TransactionService) *TransactionHandler {
	return &TransactionHandler{transactionSvc: transactionSvc}
}

func (h *TransactionHandler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/accounts/{id}/deposit", h.Deposit).Methods("POST")
	router.HandleFunc("/accounts/{id}/withdraw", h.Withdraw).Methods("POST")
	router.HandleFunc("/accounts/{id}/transactions", h.GetTransactions).Methods("GET")
}

func (h *TransactionHandler) Deposit(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(int)
	_ = userID
	vars := mux.Vars(r)
	accountID, _ := strconv.Atoi(vars["id"])

	var req struct {
		Amount      float64 `json:"amount" validate:"required,gt=0"`
		Description string  `json:"description"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// В реальном приложении email должен браться из БД по userID
	userEmail := "user@example.com"

	if err := h.transactionSvc.ProcessDeposit(
		accountID,
		req.Amount,
		req.Description,
		userEmail,
	); err != nil {
		switch err {
		case service.ErrAccountNotFound:
			http.Error(w, "Account not found", http.StatusNotFound)
		default:
			http.Error(w, "Deposit failed", http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Deposit successful"})
}

func (h *TransactionHandler) Withdraw(w http.ResponseWriter, r *http.Request) {
	// Аналогично Deposit, но вызывает ProcessWithdrawal
}

func (h *TransactionHandler) GetTransactions(w http.ResponseWriter, r *http.Request) {
	// Реализация получения истории транзакций
}
