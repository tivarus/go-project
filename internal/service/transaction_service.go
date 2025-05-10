package service

import (
	"bank-api/internal/models"
	"bank-api/internal/repository"
	"database/sql"
	"log"
)

type TransactionService struct {
	transactionRepo *repository.TransactionRepository
	accountRepo     *repository.AccountRepository
	notificationSvc *NotificationService
	db              *sql.DB
}

func NewTransactionService(
	transactionRepo *repository.TransactionRepository,
	accountRepo *repository.AccountRepository,
	notificationSvc *NotificationService,
	db *sql.DB,
) *TransactionService {
	return &TransactionService{
		transactionRepo: transactionRepo,
		accountRepo:     accountRepo,
		notificationSvc: notificationSvc,
		db:              db,
	}
}

func (s *TransactionService) ProcessDeposit(accountID int, amount float64, description string, userEmail string) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Проверяем существование счета
	account, err := s.accountRepo.GetAccountByID(accountID)
	if err != nil {
		return err
	}
	if account == nil {
		return ErrAccountNotFound
	}

	// Обновляем баланс
	if err := s.accountRepo.UpdateBalance(accountID, amount); err != nil {
		return err
	}

	// Создаем запись о транзакции
	transaction := &models.Transaction{
		AccountID:   accountID,
		Amount:      amount,
		Type:        models.TransactionDeposit,
		Description: description,
	}

	if err := s.transactionRepo.CreateTransaction(tx, transaction); err != nil {
		return err
	}

	// Отправляем уведомление
	if err := s.notificationSvc.SendPaymentNotification(
		userEmail,
		amount,
	); err != nil {
		log.Printf("Failed to send deposit notification: %v", err)
	}

	return tx.Commit()
}

func (s *TransactionService) ProcessWithdrawal(accountID int, amount float64, description string, userEmail string) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Проверяем существование счета и достаточность средств
	account, err := s.accountRepo.GetAccountByID(accountID)
	if err != nil {
		return err
	}
	if account == nil {
		return ErrAccountNotFound
	}
	if account.Balance < amount {
		return ErrInsufficientFunds
	}

	// Обновляем баланс
	if err := s.accountRepo.UpdateBalance(accountID, -amount); err != nil {
		return err
	}

	// Создаем запись о транзакции
	transaction := &models.Transaction{
		AccountID:   accountID,
		Amount:      amount,
		Type:        models.TransactionWithdrawal,
		Description: description,
	}

	if err := s.transactionRepo.CreateTransaction(tx, transaction); err != nil {
		return err
	}

	// Отправляем уведомление
	if err := s.notificationSvc.SendPaymentNotification(
		userEmail,
		amount,
	); err != nil {
		log.Printf("Failed to send withdrawal notification: %v", err)
	}

	return tx.Commit()
}

func (s *TransactionService) ProcessTransfer(fromAccountID, toAccountID int, amount float64, description string, userEmail string) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Проверяем счета
	fromAccount, err := s.accountRepo.GetAccountByID(fromAccountID)
	if err != nil {
		return err
	}
	if fromAccount == nil {
		return ErrAccountNotFound
	}

	toAccount, err := s.accountRepo.GetAccountByID(toAccountID)
	if err != nil {
		return err
	}
	if toAccount == nil {
		return ErrAccountNotFound
	}

	// Проверяем достаточность средств
	if fromAccount.Balance < amount {
		return ErrInsufficientFunds
	}

	// Выполняем перевод
	if err := s.accountRepo.UpdateBalance(fromAccountID, -amount); err != nil {
		return err
	}
	if err := s.accountRepo.UpdateBalance(toAccountID, amount); err != nil {
		return err
	}

	// Создаем записи о транзакциях
	fromTransaction := &models.Transaction{
		AccountID:   fromAccountID,
		Amount:      amount,
		Type:        models.TransactionTransfer,
		Description: description,
	}

	toTransaction := &models.Transaction{
		AccountID:   toAccountID,
		Amount:      amount,
		Type:        models.TransactionTransfer,
		Description: description,
	}

	if err := s.transactionRepo.CreateTransaction(tx, fromTransaction); err != nil {
		return err
	}
	if err := s.transactionRepo.CreateTransaction(tx, toTransaction); err != nil {
		return err
	}

	// Отправляем уведомление
	if err := s.notificationSvc.SendPaymentNotification(
		userEmail,
		amount,
	); err != nil {
		log.Printf("Failed to send transfer notification: %v", err)
	}

	return tx.Commit()
}
