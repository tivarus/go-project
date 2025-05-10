package service

import (
	"bank-api/internal/models"
	"bank-api/internal/repository"
	"bank-api/pkg/cbr"
	"database/sql"
	"fmt"
	"log"
	"math"
	"time"
)

type CreditService struct {
	creditRepo      *repository.CreditRepository
	accountRepo     *repository.AccountRepository
	accountService  *AccountService
	notificationSvc *NotificationService
	db              *sql.DB
}

func NewCreditService(
	creditRepo *repository.CreditRepository,
	accountRepo *repository.AccountRepository,
	accountService *AccountService,
	notificationSvc *NotificationService,
	db *sql.DB,
) *CreditService {
	return &CreditService{
		creditRepo:      creditRepo,
		accountRepo:     accountRepo,
		accountService:  accountService,
		notificationSvc: notificationSvc,
		db:              db,
	}
}

func (s *CreditService) CreateCredit(userID int, req *models.CreateCreditRequest, userEmail string) (*models.Credit, error) {
	account, err := s.accountRepo.GetAccountByID(req.AccountID)
	if err != nil {
		return nil, err
	}
	if account == nil || account.UserID != userID {
		return nil, ErrAccountNotFound
	}

	keyRate, err := cbr.GetKeyRate()
	if err != nil {
		return nil, fmt.Errorf("failed to get key rate: %w", err)
	}

	credit := &models.Credit{
		AccountID:    req.AccountID,
		Amount:       req.Amount,
		InterestRate: keyRate,
		TermMonths:   req.TermMonths,
		StartDate:    time.Now(),
		Status:       models.CreditStatusActive,
	}

	tx, err := s.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	if err := s.creditRepo.CreateCredit(credit); err != nil {
		return nil, err
	}

	payments := s.generatePaymentSchedule(credit)
	for _, payment := range payments {
		if err := s.creditRepo.CreatePaymentSchedule(tx, payment); err != nil {
			return nil, err
		}
	}

	if err := s.accountService.ProcessCreditDeposit(credit.AccountID, credit.Amount, userEmail); err != nil {
		return nil, err
	}

	if s.notificationSvc != nil {
		if err := s.notificationSvc.SendCreditNotification(
			userEmail,
			credit.Amount,
			credit.TermMonths,
		); err != nil {
			log.Printf("Failed to send credit notification: %v", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return credit, nil
}

func (s *CreditService) generatePaymentSchedule(credit *models.Credit) []*models.PaymentSchedule {
	var payments []*models.PaymentSchedule

	monthlyRate := credit.InterestRate / 100 / 12
	annuityPayment := (credit.Amount * monthlyRate) / (1 - math.Pow(1+monthlyRate, float64(-credit.TermMonths)))

	remainingPrincipal := credit.Amount
	currentDate := credit.StartDate

	for i := 0; i < credit.TermMonths; i++ {
		currentDate = currentDate.AddDate(0, 1, 0)

		interest := remainingPrincipal * monthlyRate
		principal := annuityPayment - interest
		remainingPrincipal -= principal

		payments = append(payments, &models.PaymentSchedule{
			CreditID:    credit.ID,
			PaymentDate: currentDate,
			Amount:      annuityPayment,
			Principal:   principal,
			Interest:    interest,
			Status:      "pending",
		})
	}

	return payments
}

func (s *CreditService) ProcessDuePayments() error {
	// Реализация обработки просроченных платежей
	// (будет вызываться шедулером)
	return nil
}
