package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"bank-api/internal/models"
	"bank-api/internal/repository"
	"bank-api/internal/service"

	"github.com/gorilla/mux"
)

type CreditHandler struct {
	creditService *service.CreditService
	accountRepo   *repository.AccountRepository
}

func NewCreditHandler(
	creditService *service.CreditService,
	accountRepo *repository.AccountRepository,
) *CreditHandler {
	return &CreditHandler{
		creditService: creditService,
		accountRepo:   accountRepo,
	}
}

func (h *CreditHandler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/credits", h.CreateCredit).Methods("POST")
	router.HandleFunc("/credits/{id}/schedule", h.GetPaymentSchedule).Methods("GET")
}

func (h *CreditHandler) CreateCredit(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(int)

	// В реальном приложении нужно получить email из БД
	userEmail := "user@example.com"

	var req models.CreateCreditRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	credit, err := h.creditService.CreateCredit(userID, &req, userEmail)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response := models.CreditResponse{
		ID:           credit.ID,
		Amount:       credit.Amount,
		InterestRate: credit.InterestRate,
		TermMonths:   credit.TermMonths,
		StartDate:    credit.StartDate,
		Status:       credit.Status,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func (h *CreditHandler) GetPaymentSchedule(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(int)
	vars := mux.Vars(r)
	creditID, _ := strconv.Atoi(vars["id"])

	// Проверяем права доступа
	credit, err := h.creditService.GetCreditByID(creditID)
	if err != nil {
		http.Error(w, "Credit not found", http.StatusNotFound)
		return
	}

	account, err := h.accountRepo.GetAccountByID(credit.AccountID)
	if err != nil || account.UserID != userID {
		http.Error(w, "Access denied", http.StatusForbidden)
		return
	}

	schedule, err := h.creditService.GetPaymentSchedule(creditID)
	if err != nil {
		http.Error(w, "Failed to get schedule", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(schedule)
}
